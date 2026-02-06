package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"m3u8-downloader/internal/config"
	"m3u8-downloader/internal/errors"
	httpClient "m3u8-downloader/internal/http"
	"m3u8-downloader/internal/logger"
	"m3u8-downloader/internal/m3u8"
	"m3u8-downloader/internal/util"
	"m3u8-downloader/internal/video"
)

// Application 应用程序
type Application struct {
	cfg              *config.Config
	logger           logger.Logger
	httpClient       httpClient.Client
	m3u8Fetcher      m3u8.Fetcher
	downloadManager  *DownloadManager
	videoMerger      video.Merger
}

// NewApplication 创建新的应用程序
func NewApplication(cfg *config.Config, logger logger.Logger) (*Application, error) {
	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 创建 HTTP 客户端
	hc := httpClient.NewClient(
		cfg.HTTP.Timeout,
		cfg.HTTP.MaxRetries,
		cfg.HTTP.UserAgent,
		logger,
	)

	if cfg.Download.InsecureSkipVerify {
		hc.(*httpClient.HTTPClient).SetInsecureVerify(true)
	}

	// 创建 M3U8 获取器
	m3u8Fetcher := m3u8.NewFetcher(hc, logger)

	// 创建下载管理器
	downloadManager := NewDownloadManager(
		hc,
		cfg.Download.MaxGoroutines,
		cfg.HTTP.MaxRetries,
		logger,
	)

	// 创建视频合并器
	videoMerger := video.NewFFmpegMerger(cfg.FFmpeg.Path, logger)

	// 检查 FFmpeg
	if err := videoMerger.CheckFFmpeg(); err != nil {
		return nil, err
	}

	return &Application{
		cfg:             cfg,
		logger:          logger,
		httpClient:      hc,
		m3u8Fetcher:     m3u8Fetcher,
		downloadManager: downloadManager,
		videoMerger:     videoMerger,
	}, nil
}

// Run 运行应用程序
func (app *Application) Run(m3u8URL, movieName string) error {
	startTime := time.Now()

	// 1. 获取 M3U8 清单
	app.logger.Info("[准备] 获取 M3U8 清单...")
	manifest, err := app.m3u8Fetcher.FetchManifest(m3u8URL, app.cfg.Download.Cookie)
	if err != nil {
		return err
	}

	app.logger.Info("[准备] 共解析 %d 个 TS 文件", len(manifest.Segments))

	// 2. 确定保存路径
	savePath := app.cfg.Download.SavePath
	if savePath == "" {
		savePath, _ = os.Getwd()
	}

	downloadDir := filepath.Join(savePath, movieName)

	// 3. 下载 TS 文件
	app.logger.Info("[准备] 开始下载到: %s", downloadDir)
	err = app.downloadManager.Download(manifest, downloadDir)
	if err != nil {
		return err
	}

	// 4. 验证下载
	app.logger.Info("[验证] 检查下载完整性...")
	stats := app.downloadManager.GetStats()
	expectedCount := int(stats.TotalCount)
	successCount := int(stats.DownloadCount)
	allowedLoss := float64(expectedCount) * app.cfg.Download.LossTolerance

	if float64(expectedCount-successCount) > allowedLoss {
		return errors.New(errors.DownloadFailed,
			fmt.Sprintf("下载不完整: 期望 %d, 成功 %d, 允许丢失 %.0f",
				expectedCount, successCount, allowedLoss), nil)
	}

	// 5. 合并视频
	app.logger.Info("[合并] 使用 FFmpeg 合并视频...")
	outputPath := filepath.Join(savePath, movieName+".mp4")
	finalPath, err := app.videoMerger.Merge(downloadDir, outputPath)
	if err != nil {
		return err
	}

	// 6. 清理临时文件
	if app.cfg.Download.AutoClear {
		app.logger.Info("[清理] 删除临时 TS 文件...")
		util.RemoveDir(downloadDir)
	}

	// 7. 显示完成信息
	elapsed := time.Since(startTime)
	fileSize, _ := util.GetFileSize(finalPath)

	fmt.Printf("\n")
	app.logger.Info("[成功] 视频已保存: %s", finalPath)
	app.logger.Info("[统计] 下载耗时: %.1fs, 文件大小: %.2f MB", elapsed.Seconds(), fileSize)

	return nil
}
