package core

import (
	"fmt"
	"math"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	httpClient "m3u8-downloader/internal/http"
	"m3u8-downloader/internal/logger"
	"m3u8-downloader/internal/m3u8"
	"m3u8-downloader/internal/theme"
	"m3u8-downloader/internal/util"
)

// DownloadStats 下载统计信息
type DownloadStats struct {
	TotalCount     int64
	DownloadCount  int64
	SkippedCount   int64
	FailedCount    int64
	StartTime      time.Time
	LastUpdateTime time.Time
}

// DownloadManager 下载管理器
type DownloadManager struct {
	httpClient     httpClient.Client
	maxGoroutines  int
	maxRetries     int
	tsNameTemplate string
	logger         logger.Logger
	stats          *DownloadStats
	// progressActive indicates whether the progress line should be redrawn.
	// 1 = active, 0 = stopped.
	progressActive int32
}

// NewDownloadManager 创建新的下载管理器
func NewDownloadManager(httpClient httpClient.Client, maxGoroutines, maxRetries int, lg logger.Logger) *DownloadManager {
	dm := &DownloadManager{
		httpClient:     httpClient,
		maxGoroutines:  maxGoroutines,
		maxRetries:     maxRetries,
		tsNameTemplate: "%05d.ts",
		logger:         lg,
		stats:          &DownloadStats{},
	}

	// Register progress redraw so log messages won't leave the progress broken
	// on the terminal. The logger package will call this after printing logs.
	logger.RegisterProgressRedraw(func() { dm.displayProgress() })

	// enable progress redraw
	atomic.StoreInt32(&dm.progressActive, 1)

	return dm
}

// Download 下载所有 TS 段
func (dm *DownloadManager) Download(manifest *m3u8.Manifest, downloadDir string) error {
	// 确保目录存在
	err := util.EnsureDir(downloadDir)
	if err != nil {
		return err
	}

	segments := manifest.Segments
	dm.stats.TotalCount = int64(len(segments))
	dm.stats.DownloadCount = 0
	dm.stats.SkippedCount = 0
	dm.stats.FailedCount = 0
	dm.stats.StartTime = time.Now()

	dm.logger.Info("开始下载 %d 个 TS 文件到 %s", len(segments), downloadDir)

	// 使用 goroutine 限制器
	limiter := make(chan struct{}, dm.maxGoroutines)
	var wg sync.WaitGroup

	for i, segment := range segments {
		wg.Add(1)
		limiter <- struct{}{} // 获取许可

		go func(index int, seg *m3u8.TsSegment) {
			defer func() {
				wg.Done()
				<-limiter // 释放许可
			}()

			dm.downloadSingleSegment(index, seg, downloadDir, manifest.Key)
			dm.displayProgress()
		}(i, segment)
	}

	wg.Wait()
	// stop progress redraw and clear the line so subsequent operations (merge)
	// won't have the progress bar print over log output.
	atomic.StoreInt32(&dm.progressActive, 0)
	// unregister redraw callback
	logger.RegisterProgressRedraw(nil)
	// clear the progress line
	fmt.Print("\r\033[K")
	fmt.Println()

	dm.logger.Info("下载完成: 成功 %d, 跳过 %d, 失败 %d",
		atomic.LoadInt64(&dm.stats.DownloadCount),
		atomic.LoadInt64(&dm.stats.SkippedCount),
		atomic.LoadInt64(&dm.stats.FailedCount),
	)

	return nil
}

func (dm *DownloadManager) downloadSingleSegment(index int, segment *m3u8.TsSegment, downloadDir string, key *m3u8.EncryptionKey) {
	filePath := filepath.Join(downloadDir, segment.Name)

	// 检查文件是否已存在
	exists, _ := util.PathExists(filePath)
	if exists {
		atomic.AddInt64(&dm.stats.SkippedCount, 1)
		return
	}

	// 重试下载
	for attempt := 1; attempt <= dm.maxRetries; attempt++ {
		data, err := dm.httpClient.Get(segment.URL)
		if err != nil {
			if attempt < dm.maxRetries {
				dm.logger.Warn("下载段 %d 失败，重试 (%d/%d): %v", index, attempt, dm.maxRetries, err)
				time.Sleep(time.Duration(attempt-1) * time.Second)
				continue
			}
			dm.logger.Error("下载段 %d 失败: %v", index, err)
			atomic.AddInt64(&dm.stats.FailedCount, 1)
			return
		}

		// 验证数据
		if len(data) == 0 {
			if attempt < dm.maxRetries {
				dm.logger.Warn("段 %d 数据为空，重试 (%d/%d)", index, attempt, dm.maxRetries)
				time.Sleep(time.Duration(attempt-1) * time.Second)
				continue
			}
			dm.logger.Error("段 %d 数据为空", index)
			atomic.AddInt64(&dm.stats.FailedCount, 1)
			return
		}

		// 解密（如果需要）
		if key != nil && len(key.Data) > 0 {
			decrypted, err := util.AesDecrypt(data, key.Data)
			if err != nil {
				if attempt < dm.maxRetries {
					dm.logger.Warn("解密段 %d 失败，重试 (%d/%d): %v", index, attempt, dm.maxRetries, err)
					time.Sleep(time.Duration(attempt-1) * time.Second)
					continue
				}
				dm.logger.Error("解密段 %d 失败: %v", index, err)
				atomic.AddInt64(&dm.stats.FailedCount, 1)
				return
			}
			data = decrypted
		}

		// 移除 TS padding
		data = util.RemoveTSPadding(data)

		// 写入文件
		err = util.WriteFile(filePath, data)
		if err != nil {
			dm.logger.Error("写入文件 %s 失败: %v", filePath, err)
			atomic.AddInt64(&dm.stats.FailedCount, 1)
			return
		}

		atomic.AddInt64(&dm.stats.DownloadCount, 1)
		return
	}
}

func (dm *DownloadManager) displayProgress() {
	// if progress redraw was turned off, don't print anything
	if atomic.LoadInt32(&dm.progressActive) == 0 {
		return
	}
	downloadCount := atomic.LoadInt64(&dm.stats.DownloadCount)
	total := dm.stats.TotalCount

	if total == 0 {
		return
	}

	progress := float32(downloadCount) / float32(total)
	elapsed := time.Since(dm.stats.StartTime).Seconds()

	speed := float64(downloadCount) / elapsed
	remainCount := float64(total) - float64(downloadCount)
	eta := remainCount / speed

	// progress bar rendering — use shared theme colors
	// colors are provided by internal/theme

	progressWidth := 36
	pos := int(progress * float32(progressWidth))

	// format ETA into hh:mm:ss or mm:ss
	etaStr := "--:--"
	if !math.IsInf(eta, 0) && !math.IsNaN(eta) {
		d := time.Duration(eta) * time.Second
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		s := int(d.Seconds()) % 60
		if h > 0 {
			etaStr = fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
		} else {
			etaStr = fmt.Sprintf("%02dm%02ds", m, s)
		}
	}

	// spinner
	sp := []string{"⣽", "⣾", "⣻", "⣷", "⣯", "⣟"}
	spinner := sp[int(time.Now().UnixNano()/1e8)%len(sp)]

	// render bar with colors using theme
	filled := ""
	if pos > 0 {
		filled = theme.Green + repeatStr("━", pos) + theme.Reset
	}
	empty := theme.Surface1 + repeatStr(" ", progressWidth-pos) + theme.Reset

	// assemble and print
	fmt.Printf("\r%s %s%s %s%d/%d %6.2f%% %s%.2f files/s ETA:%s %s",
		theme.Lavender+"Vid Kbps"+theme.Reset,
		filled,
		empty,
		theme.Text,
		downloadCount,
		total,
		progress*100,
		theme.Text,
		speed,
		etaStr,
		spinner,
	)
}

func repeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// GetStats 获取下载统计信息
func (dm *DownloadManager) GetStats() *DownloadStats {
	return dm.stats
}
