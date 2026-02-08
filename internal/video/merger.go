package video

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"m3u8-downloader/internal/errors"
	"m3u8-downloader/internal/logger"
	"m3u8-downloader/internal/theme"
	"m3u8-downloader/internal/util"
)

// Merger 视频合并接口
type Merger interface {
	Merge(segmentDir, outputPath string) (string, error)
	Validate(outputPath string) error
}

// FFmpegMerger FFmpeg 视频合并实现
type FFmpegMerger struct {
	ffmpegPath string
	logger     logger.Logger
}

// NewFFmpegMerger 创建新的 FFmpeg 合并器
func NewFFmpegMerger(ffmpegPath string, logger logger.Logger) *FFmpegMerger {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	return &FFmpegMerger{
		ffmpegPath: ffmpegPath,
		logger:     logger,
	}
}

// CheckFFmpeg 检查 FFmpeg 是否可用
func (m *FFmpegMerger) CheckFFmpeg() error {
	cmd := exec.Command(m.ffmpegPath, "-version")
	_, err := cmd.Output()
	if err != nil {
		return errors.New(
			errors.FFmpegNotFound,
			fmt.Sprintf("FFmpeg 不可用 (%s)", m.ffmpegPath),
			err,
		)
	}
	m.logger.Debug("FFmpeg 可用: %s", m.ffmpegPath)
	return nil
}

// Merge 合并 TS 文件为 MP4
func (m *FFmpegMerger) Merge(segmentDir, outputPath string) (string, error) {
	// 列出所有 TS 文件
	tsFiles, err := util.ListTSFiles(segmentDir)
	if err != nil {
		return "", err
	}

	if len(tsFiles) == 0 {
		return "", errors.New(errors.MergeFailed, "目录中未找到 TS 文件", nil)
	}

	// 排序文件名
	sort.Strings(tsFiles)

	// 创建 concat 文件
	concatFile := filepath.Join(segmentDir, "concat.txt")
	concatContent := ""
	for _, f := range tsFiles {
		filePath := filepath.Join(segmentDir, f)
		concatContent += fmt.Sprintf("file '%s'\n", filePath)
	}

	err = util.WriteFile(concatFile, []byte(concatContent))
	if err != nil {
		return "", errors.New(errors.MergeFailed, "创建 concat 文件失败", err)
	}
	defer os.Remove(concatFile)

	m.logger.Info("开始合并 %d 个 TS 文件到 %s", len(tsFiles), theme.Lavender+outputPath+theme.Reset)

	// 执行 FFmpeg 合并
	cmd := exec.Command(
		m.ffmpegPath,
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		"-y",
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", errors.New(
			errors.FFmpegFailed,
			"FFmpeg 执行失败",
			fmt.Errorf("%v\n%s", err, stderr.String()),
		)
	}

	// 验证输出文件
	err = m.Validate(outputPath)
	if err != nil {
		return "", err
	}

	m.logger.Info("成功合并视频: %s", theme.Lavender+outputPath+theme.Reset)
	return outputPath, nil
}

// Validate 验证输出文件
func (m *FFmpegMerger) Validate(outputPath string) error {
	exists, err := util.PathExists(outputPath)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(errors.MergeFailed, "输出文件不存在", nil)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return errors.New(errors.MergeFailed, "无法读取输出文件", err)
	}

	if info.Size() == 0 {
		return errors.New(errors.MergeFailed, "输出文件为空", nil)
	}

	return nil
}
