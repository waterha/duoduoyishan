package utils

import (
	"duoduoyishan/config"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UploadInfo struct {
	FileName string
	FileSize int64
	FileExt  string
	SavePath string
	FileURL  string
}

// 检查文件扩展名是否允许
func CheckFileExt(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowExt := range config.GlobalConfig.Upload.AllowExts {
		if ext == allowExt {
			return true
		}
	}
	return false
}

// 检查文件大小
func CheckFileSize(size int64) bool {
	return size <= config.GlobalConfig.Upload.MaxSize
}

// 保存上传文件
func SaveUploadFile(file *multipart.FileHeader, subDir string) (*UploadInfo, error) {
	// 检查扩展名
	if !CheckFileExt(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型")
	}

	// 检查文件大小
	if !CheckFileSize(file.Size) {
		return nil, fmt.Errorf("文件大小超过限制")
	}

	// 生成新文件名
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext

	// 创建日期子目录
	dateDir := time.Now().Format("20060102")
	saveDir := filepath.Join(config.GlobalConfig.Upload.SavePath, subDir, dateDir)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, err
	}

	// 完整保存路径
	savePath := filepath.Join(saveDir, newFileName)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// 文件URL
	fileURL := fmt.Sprintf("/uploads/%s/%s/%s", subDir, dateDir, newFileName)

	return &UploadInfo{
		FileName: file.Filename,
		FileSize: file.Size,
		FileExt:  ext,
		SavePath: savePath,
		FileURL:  fileURL,
	}, nil
}
