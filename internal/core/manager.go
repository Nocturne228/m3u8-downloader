package core

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	httpClient "m3u8-downloader/internal/http"
	"m3u8-downloader/internal/logger"
	"m3u8-downloader/internal/m3u8"
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
}

// NewDownloadManager 创建新的下载管理器
func NewDownloadManager(httpClient httpClient.Client, maxGoroutines, maxRetries int, logger logger.Logger) *DownloadManager {
	return &DownloadManager{
		httpClient:     httpClient,
		maxGoroutines:  maxGoroutines,
		maxRetries:     maxRetries,
		tsNameTemplate: "%05d.ts",
		logger:         logger,
		stats:          &DownloadStats{},
	}
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
	fmt.Println() // 打印换行

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

	// 计算进度条
	progressWidth := 20
	pos := int(progress * float32(progressWidth))

	// 格式化 ETA
	etaStr := fmt.Sprintf("%.0fs", eta)
	if eta > 60 {
		etaStr = fmt.Sprintf("%.1fm", eta/60)
	}

	// 显示进度条
	fmt.Printf("\r[Downloading] %s%s %6.2f%% (%d/%d) %.1f files/s ETA: %s",
		repeatStr("■", pos),
		repeatStr(" ", progressWidth-pos),
		progress*100,
		downloadCount,
		total,
		speed,
		etaStr,
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
