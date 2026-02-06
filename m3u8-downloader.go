// @author:
// @date:
// @
package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/levigross/grequests"
)

const (
	// HEAD_TIMEOUT 请求头超时时间
	HEAD_TIMEOUT = 5 * time.Second
	// PROGRESS_WIDTH 进度条长度
	PROGRESS_WIDTH = 20
	// TS_NAME_TEMPLATE ts视频片段命名规则
	TS_NAME_TEMPLATE = "%05d.ts"
	// MAX_RETRY_COUNT 最大重试次数
	MAX_RETRY_COUNT = 5
)

// DownloadStats 下载统计信息
type DownloadStats struct {
	TotalCount    int64
	DownloadCount int64
	StartTime     time.Time
	LastUpdateTime time.Time
	LastDownloadCount int64
}

var (
	// 命令行参数
	urlFlag = flag.String("u", "", "m3u8下载地址(http(s)://url/xx/xx/index.m3u8)")
	nFlag   = flag.Int("n", 24, "num:下载线程数(默认24)")
	htFlag  = flag.String("ht", "v1", "hostType:设置getHost的方式(v1: `http(s):// + url.Host + filepath.Dir(url.Path)`; v2: `http(s)://+ u.Host`")
	oFlag   = flag.String("o", "movie", "movieName:自定义文件名(默认为movie)不带后缀")
	cFlag   = flag.String("c", "", "cookie:自定义请求cookie")
	rFlag   = flag.Bool("r", true, "autoClear:是否自动清除ts文件")
	sFlag   = flag.Int("s", 0, "InsecureSkipVerify:是否允许不安全的请求(默认0)")
	spFlag  = flag.String("sp", "", "savePath:文件保存的绝对路径(默认为当前路径,建议默认值)")

	logger *log.Logger
	stats  = &DownloadStats{}
	statsMu sync.Mutex
	ro     = &grequests.RequestOptions{
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
		RequestTimeout: HEAD_TIMEOUT,
		Headers: map[string]string{
			"Connection":      "keep-alive",
			"Accept":          "*/*",
			"Accept-Encoding": "*",
			"Accept-Language": "zh-CN,zh;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
		},
	}
)

// TsInfo 用于保存 ts 文件的下载地址和文件名
type TsInfo struct {
	Name string
	Url  string
}

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	Run()
}

func Run() {
	msgTpl := `
========================================
     M3U8 视频下载工具 v2.0
=========================================
[功能]  多线程下载直播流 m3u8 视频
[特性]  
  - 支持加密 TS 文件自动解密
  - 使用 FFmpeg 合并生成 MP4 格式
  - 支持断点续传和失败重试
  - 显示下载速度和 ETA
[提醒]  
  - 下载失败，请使用 -ht=v2 
  - 下载失败，m3u8 地址可能存在嵌套
  - 进度条中途下载失败，可重复执行
  - 需要系统安装 FFmpeg
========================================
`
	fmt.Println(msgTpl)
	runtime.GOMAXPROCS(runtime.NumCPU())
	now := time.Now()

	// 1、解析命令行参数
	flag.Parse()
	m3u8Url := *urlFlag
	maxGoroutines := *nFlag
	hostType := *htFlag
	movieName := *oFlag
	autoClearFlag := *rFlag
	cookie := *cFlag
	insecure := *sFlag
	savePath := *spFlag

	ro.Headers["Referer"] = getHost(m3u8Url, "v2")
	if insecure != 0 {
		ro.InsecureSkipVerify = true
	}
	// http 自定义 cookie
	if cookie != "" {
		ro.Headers["Cookie"] = cookie
	}

	// 参数验证
	if !strings.HasPrefix(m3u8Url, "http") || m3u8Url == "" {
		fmt.Println("[Error] 请输入有效的 m3u8 下载地址")
		flag.Usage()
		return
	}

	if maxGoroutines <= 0 || maxGoroutines > 256 {
		fmt.Println("[Error] 线程数必须在 1-256 之间")
		return
	}

	if movieName == "" {
		fmt.Println("[Error] 输出文件名不能为空")
		return
	}

	var download_dir string
	pwd, err := os.Getwd()
	if err != nil {
		logger.Printf("[Error] 获取当前目录失败: %v\n", err)
		return
	}

	if savePath != "" {
		pwd = savePath
	}

	// 初始化下载ts的目录，后面所有的ts文件会保存在这里
	download_dir = filepath.Join(pwd, movieName)
	if isExist, _ := pathExists(download_dir); !isExist {
		err := os.MkdirAll(download_dir, os.ModePerm)
		if err != nil {
			logger.Printf("[Error] 创建下载目录失败: %v\n", err)
			return
		}
		fmt.Printf("[Info] 创建临时目录: %s\n", download_dir)
	}

	// 2、检查ffmpeg是否可用
	if !checkFFmpeg() {
		fmt.Printf("\n[Failed] 系统未安装 FFmpeg，请先安装 FFmpeg\n")
		fmt.Println("macOS: brew install ffmpeg")
		fmt.Println("Linux: apt-get install ffmpeg 或 yum install ffmpeg")
		fmt.Println("Windows: https://ffmpeg.org/download.html")
		return
	}

	// 3、解析m3u8
	m3u8Host := getHost(m3u8Url, hostType)
	m3u8Body := getM3u8Body(m3u8Url)
	//m3u8Body := getFromFile()
	ts_key := getM3u8Key(m3u8Host, m3u8Body)
	if ts_key != "" {
		fmt.Printf("待解密 ts 文件 key : %s \n", ts_key)
	}
	ts_list := getTsList(m3u8Host, m3u8Body)
	fmt.Printf("待下载 ts 文件数量: %d 个\n", len(ts_list))

	// 4、初始化下载统计
	stats.TotalCount = int64(len(ts_list))
	stats.DownloadCount = 0
	stats.StartTime = time.Now()
	stats.LastUpdateTime = time.Now()

	// 5、下载ts文件到download_dir
	downloader(ts_list, maxGoroutines, download_dir, ts_key)
	if !checkTsDownDir(download_dir, len(ts_list)) {
		fmt.Printf("\n[Failed] 请检查url地址有效性，部分或全部 ts 文件下载失败\n")
		return
	}

	// 6、合并ts切割文件成mp4文件
	fmt.Println("\n[Info] 开始使用 FFmpeg 合并视频...")
	mv, err := mergeWithFFmpeg(download_dir, movieName)
	if err != nil {
		fmt.Printf("\n[Failed] 合并视频失败: %v\n", err)
		return
	}

	if autoClearFlag {
		//自动清除ts文件目录
		os.RemoveAll(download_dir)
		fmt.Println("[Info] 已清除临时 ts 文件")
	}

	//7、输出下载视频信息
	fmt.Printf("\n[Success] 下载保存路径：%s\n", mv)
	fmt.Printf("共耗时: %6.2fs\n", time.Since(now).Seconds())
}

// 获取m3u8地址的host
// @modify: 2026-02-06 改进错误处理
func getHost(Url, ht string) (host string) {
	u, err := url.Parse(Url)
	if err != nil {
		logger.Printf("[warn] 解析 URL 失败: %v，使用版本 v2\n", err)
		return ""
	}
	switch ht {
	case "v1":
		host = u.Scheme + "://" + u.Host + filepath.Dir(u.EscapedPath())
	case "v2":
		host = u.Scheme + "://" + u.Host
	default:
		logger.Printf("[warn] 未知的 hostType: %s，使用版本 v2\n", ht)
		host = u.Scheme + "://" + u.Host
	}
	return
}

// 获取m3u8地址的内容体
// @modify: 2026-02-06 改进错误处理和重试机制
func getM3u8Body(Url string) string {
	maxRetries := 3
	for i := 1; i <= maxRetries; i++ {
		r, err := grequests.Get(Url, ro)
		if err != nil {
			if i < maxRetries {
				logger.Printf("[warn] 第 %d 次获取 m3u8 文件失败，2秒后重试: %v\n", i, err)
				time.Sleep(2 * time.Second)
				continue
			}
			checkErr(fmt.Errorf("获取 m3u8 文件失败（已重试 %d 次）: %v", maxRetries, err))
		}
		if !r.Ok {
			if i < maxRetries {
				logger.Printf("[warn] 第 %d 次获取 m3u8 文件返回状态 %d，2秒后重试\n", i, r.StatusCode)
				time.Sleep(2 * time.Second)
				continue
			}
			checkErr(fmt.Errorf("获取 m3u8 文件失败，HTTP 状态码: %d", r.StatusCode))
		}
		return r.String()
	}
	return ""
}

// 获取m3u8加密的密钥
// @modify: 2026-02-06 改进日志输出，使用 logger 替代 fmt.Println
func getM3u8Key(host, html string) (key string) {
	lines := strings.Split(html, "\n")
	key = ""
	for _, line := range lines {
		if strings.Contains(line, "#EXT-X-KEY") {
			if !strings.Contains(line, "URI") {
				continue
			}
			logger.Println("[debug] line_key:", line)
			uri_pos := strings.Index(line, "URI")
			quotation_mark_pos := strings.LastIndex(line, "\"")
			if uri_pos == -1 || quotation_mark_pos == -1 {
				continue
			}
			key_url := strings.Split(line[uri_pos:quotation_mark_pos], "\"")[1]
			if !strings.Contains(key_url, "http") {
				key_url = fmt.Sprintf("%s/%s", host, key_url)
			}
			res, err := grequests.Get(key_url, ro)
			if err != nil {
				logger.Println("[warn] 获取密钥失败:", err)
				continue
			}
			if res.StatusCode == 200 {
				key = res.String()
				break
			}
		}
	}
	logger.Println("[debug] m3u8Host:", host, "m3u8Key:", key)
	return
}

func getTsList(host, body string) (tsList []TsInfo) {
	lines := strings.Split(body, "\n")
	index := 0
	var ts TsInfo
	for _, line := range lines {
		line = strings.TrimSpace(line) // 移除前后空白
		if !strings.HasPrefix(line, "#") && line != "" {
			//有可能出现的二级嵌套格式的m3u8,请自行转换！
			index++
			if strings.HasPrefix(line, "http") {
				ts = TsInfo{
					Name: fmt.Sprintf(TS_NAME_TEMPLATE, index),
					Url:  line,
				}
				tsList = append(tsList, ts)
			} else {
				line = strings.TrimPrefix(line, "/")
				ts = TsInfo{
					Name: fmt.Sprintf(TS_NAME_TEMPLATE, index),
					Url:  fmt.Sprintf("%s/%s", host, line),
				}
				tsList = append(tsList, ts)
			}
		}
	}
	return
}

func getFromFile() string {
	data, _ := ioutil.ReadFile("./ts.txt")
	return string(data)
}

// 下载ts文件
// @modify: 2020-08-13 修复ts格式SyncByte合并不能播放问题
// @modify: 2026-02-06 改进进度统计，使用原子操作
func downloadTsFile(ts TsInfo, download_dir, key string, retries int) {
	if retries <= 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			downloadTsFile(ts, download_dir, key, retries-1)
		}
	}()
	curr_path_file := filepath.Join(download_dir, ts.Name)
	if exists, _ := pathExists(curr_path_file); exists {
		atomic.AddInt64(&stats.DownloadCount, 1)
		return
	}
	res, err := grequests.Get(ts.Url, ro)
	if err != nil || !res.Ok {
		if retries > 0 {
			time.Sleep(time.Duration(MAX_RETRY_COUNT-retries) * time.Second) // 指数退避
			downloadTsFile(ts, download_dir, key, retries-1)
			return
		}
		return
	}
	// 校验长度是否合法
	var origData []byte
	origData = res.Bytes()
	contentLen := 0
	contentLenStr := res.Header.Get("Content-Length")
	if contentLenStr != "" {
		contentLen, _ = strconv.Atoi(contentLenStr)
	}
	if len(origData) == 0 || (contentLen > 0 && len(origData) < contentLen) || res.Error != nil {
		if retries > 0 {
			time.Sleep(time.Duration(MAX_RETRY_COUNT-retries) * time.Second)
			downloadTsFile(ts, download_dir, key, retries-1)
		}
		return
	}
	// 解密出视频 ts 源文件
	if key != "" {
		//解密 ts 文件，算法：aes 128 cbc pack5
		origData, err = AesDecrypt(origData, []byte(key))
		if err != nil {
			if retries > 0 {
				downloadTsFile(ts, download_dir, key, retries-1)
			}
			return
		}
	}
	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the bytes before the SyncByte 0x47(71).
	syncByte := uint8(71) //0x47
	bLen := len(origData)
	for j := 0; j < bLen; j++ {
		if origData[j] == syncByte {
			origData = origData[j:]
			break
		}
	}
	ioutil.WriteFile(curr_path_file, origData, 0666)
	atomic.AddInt64(&stats.DownloadCount, 1)
}

// downloader m3u8 下载器
// @modify: 2026-02-06 改进进度显示，添加下载速度和 ETA
func downloader(tsList []TsInfo, maxGoroutines int, downloadDir string, key string) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxGoroutines) //chan struct 内存占用 0 bool 占用 1
	tsLen := len(tsList)

	for _, ts := range tsList {
		wg.Add(1)
		limiter <- struct{}{}
		go func(ts TsInfo, downloadDir, key string) {
			defer func() {
				wg.Done()
				<-limiter
			}()
			downloadTsFile(ts, downloadDir, key, MAX_RETRY_COUNT)
			updateProgressBar(tsLen)
		}(ts, downloadDir, key)
	}
	wg.Wait()
	fmt.Println() // 打印换行
}

// updateProgressBar 更新进度条并显示速度和ETA
func updateProgressBar(total int) {
	downloadCount := atomic.LoadInt64(&stats.DownloadCount)
	proportion := float32(downloadCount) / float32(total)
	
	now := time.Now()
	elapsed := now.Sub(stats.StartTime).Seconds()
	speed := float64(downloadCount) / elapsed
	remainCount := float64(total) - float64(downloadCount)
	eta := remainCount / speed
	
	pos := int(proportion * float32(PROGRESS_WIDTH))
	speedStr := fmt.Sprintf("%.1f files/s", speed)
	etaStr := fmt.Sprintf("ETA: %.0fs", eta)
	if eta > 60 {
		etaStr = fmt.Sprintf("ETA: %.1fm", eta/60)
	}
	
	s := fmt.Sprintf("[Downloading] %s%*s %6.2f%% (%d/%d) %s %s",
		strings.Repeat("■", pos), PROGRESS_WIDTH-pos, "", proportion*100,
		downloadCount, total, speedStr, etaStr)
	fmt.Print("\r" + s)
}

func checkTsDownDir(dir string, expectedCount int) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	tsCount := 0
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".ts" {
			tsCount++
		}
	}
	// 允许损失率 10%
	threshold := int(float64(expectedCount) * 0.9)
	return tsCount >= threshold
}

// 合并ts文件 - 已过时，使用 mergeWithFFmpeg 代替
// @deprecated
func mergeTs(downloadDir string) string {
	mvName := downloadDir + ".mp4"
	outMv, _ := os.Create(mvName)
	defer outMv.Close()
	writer := bufio.NewWriter(outMv)
	err := filepath.Walk(downloadDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() || filepath.Ext(path) != ".ts" {
			return nil
		}
		bytes, _ := ioutil.ReadFile(path)
		_, err = writer.Write(bytes)
		return err
	})
	checkErr(err)
	_ = writer.Flush()
	return mvName
}

// checkFFmpeg 检查系统是否安装了 FFmpeg
func checkFFmpeg() bool {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.Output()
	return err == nil
}

// mergeWithFFmpeg 使用 FFmpeg 合并 ts 文件成 mp4
// 优点：正确的视频编码、容器格式、元数据等
// @modify: 2026-02-06 新增函数，使用 FFmpeg 进行合并
func mergeWithFFmpeg(tsDir, movieName string) (string, error) {
	// 获取所有 ts 文件并按名称排序
	files, err := ioutil.ReadDir(tsDir)
	if err != nil {
		return "", fmt.Errorf("读取目录失败: %v", err)
	}

	var tsFiles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".ts" {
			tsFiles = append(tsFiles, f.Name())
		}
	}

	if len(tsFiles) == 0 {
		return "", fmt.Errorf("目录中未找到 ts 文件")
	}

	// 按文件名排序（保证正确顺序）
	sort.Strings(tsFiles)

	// 创建 concat demuxer 文件
	concatFile := filepath.Join(tsDir, "concat.txt")
	concatContent := ""
	for _, f := range tsFiles {
		filePath := filepath.Join(tsDir, f)
		concatContent += fmt.Sprintf("file '%s'\n", filePath)
	}

	err = ioutil.WriteFile(concatFile, []byte(concatContent), 0644)
	if err != nil {
		return "", fmt.Errorf("创建 concat 文件失败: %v", err)
	}
	defer os.Remove(concatFile)

	// 输出文件路径
	pwd, _ := os.Getwd()
	outputPath := filepath.Join(pwd, movieName+".mp4")

	// 使用 FFmpeg 进行合并
	// -c copy: 不重新编码，直接复制流（快速）
	// -bsf:a aac_adtstoasc: 转换 AAC 格式（某些情况需要）
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		"-y",
		outputPath,
	)

	// 捕获 stderr 以显示进度
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("FFmpeg 执行失败: %v\nDetails: %s", err, stderr.String())
	}

	// 验证输出文件
	if exists, _ := pathExists(outputPath); !exists {
		return "", fmt.Errorf("输出文件不存在: %s", outputPath)
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil || fileInfo.Size() == 0 {
		return "", fmt.Errorf("输出文件为空或不可读: %s", outputPath)
	}

	return outputPath, nil
}

// 进度条（已改进）
// @modify: 2026-02-06 改进显示格式
func DrawProgressBar(prefix string, proportion float32, width int, suffix ...string) {
	pos := int(proportion * float32(width))
	s := fmt.Sprintf("[%s] %s%*s %6.2f%% \t%s",
		prefix, strings.Repeat("■", pos), width-pos, "", proportion*100, strings.Join(suffix, ""))
	fmt.Print("\r" + s)
}

// drawMergeProgressBar 显示合并进度
func drawMergeProgressBar(percentage float64, message string) {
	pos := int(percentage * float64(PROGRESS_WIDTH))
	s := fmt.Sprintf("[Merging] %s%*s %6.2f%% %s",
		strings.Repeat("■", pos), PROGRESS_WIDTH-pos, "", percentage*100, message)
	fmt.Print("\r" + s)
}

// ============================== shell相关 ==============================
// 判断文件是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// getFileSize 获取文件大小（单位MB）
func getFileSize(path string) (float64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return float64(info.Size()) / (1024 * 1024), nil
}

// 执行 shell
func execUnixShell(s string) {
	cmd := exec.Command("bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out.String())
}

func execWinShell(s string) error {
	cmd := exec.Command("cmd", "/C", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out.String())
	return nil
}

// windows 合并文件（已过时，使用 FFmpeg 代替）
// @deprecated
func win_merge_file(path string) {
	pwd, _ := os.Getwd()
	os.Chdir(path)
	execWinShell("copy /b *.ts merge.tmp")
	execWinShell("del /Q *.ts")
	os.Rename("merge.tmp", "merge.mp4")
	os.Chdir(pwd)
}

// unix 合并文件（已过时，使用 FFmpeg 代替）
// @deprecated
func unix_merge_file(path string) {
	pwd, _ := os.Getwd()
	os.Chdir(path)
	//cmd := `ls  *.ts |sort -t "\." -k 1 -n |awk '{print $0}' |xargs -n 1 -I {} bash -c "cat {} >> new.tmp"`
	cmd := `cat *.ts >> merge.tmp`
	execUnixShell(cmd)
	execUnixShell("rm -rf *.ts")
	os.Rename("merge.tmp", "merge.mp4")
	os.Chdir(pwd)
}

// ============================== 加解密相关 ==============================

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func checkErr(e error) {
	if e != nil {
		logger.Panic(e)
	}
}
