package source

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileServer 创建一个文件服务器，提供文件列表和下载功能
func FileServer(dest string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/filelist", func(w http.ResponseWriter, r *http.Request) {
		fileList, err := GetFileList(dest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fileList)
	})
	mux.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/download/"):]
		err := serveFile(dest, path, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error serving file: %v", err), http.StatusInternalServerError)
			return
		}
	})
	return mux
}

func serveFile(basePath, filePath string, w http.ResponseWriter) error {
	// 构建本地文件路径
	localPath := filepath.Join(basePath, filepath.Clean(filePath))

	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Error getting file info: %v", err)
	}

	// 设置响应头
	contentType := getContentTypeFromExtension(filepath.Ext(filePath))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 将文件内容复制到响应中
	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("Error writing file to response: %v", err)
	}

	return nil
}

// getContentTypeFromExtension 根据文件扩展名返回相应的 Content-Type
func getContentTypeFromExtension(ext string) string {
	// 转换为小写
	ext = strings.ToLower(ext)

	// MIME 类型映射表
	mimeTypes := map[string]string{
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".mp3":  "audio/mpeg",
		"":      "application/octet-stream",
	}

	// 根据扩展名查找 MIME 类型
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}

	// 如果找不到，则使用默认值
	return mimeTypes[""]
}
