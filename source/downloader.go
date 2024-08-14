package source

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v50/github"
	"golang.org/x/net/proxy"
	"golang.org/x/oauth2"
)

// DownloadRepoFiles 下载仓库文件
func DownloadRepoFiles(config *Config) error {
	fmt.Printf("Start check and download files\n")
	// 创建 GitHub 客户端
	ctx := context.Background()
	client := createGitHubClient(ctx, config)

	// 获取存储库的文件树
	tree, _, err := client.Git.GetTree(ctx, config.GitHub.Owner, config.GitHub.Repo, "HEAD", true)
	if err != nil {
		return fmt.Errorf("Error getting tree: %v", err)
	}

	// 创建进度条
	bar := pb.StartNew(len(tree.Entries))

	// 创建工作池
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.MaxConcurrent)

	for _, entry := range tree.Entries {
		if *entry.Type == "blob" {
			// 使用工作池
			wg.Add(1)
			sem <- struct{}{} // Acquire a slot

			go func(entry *github.TreeEntry) {
				defer wg.Done()
				defer func() { <-sem }() // Release the slot

				err := downloadFileIfNotExists(client, config.GitHub.Owner, config.GitHub.Repo, *entry.Path, *entry.SHA, config.Dest)
				if err != nil {
					fmt.Printf("Error downloading file %s: %v\n", *entry.Path, err)
				}
				bar.Increment()
			}(entry)
		}
	}

	// 等待所有下载完成
	wg.Wait()
	bar.Finish()
	return nil
}

// createGitHubClient 创建一个支持 SOCKS 代理的 GitHub 客户端
func createGitHubClient(ctx context.Context, config *Config) *github.Client {
	var httpClient *http.Client
	if config.ProxyURL != "" {
		// 创建 SOCKS5 拨号器
		dialer, err := proxy.SOCKS5("tcp", config.ProxyURL, nil, proxy.Direct)
		if err != nil {
			fmt.Printf("Error creating SOCKS5 dialer: %v\n", err)
			return nil
		}

		// 创建带有 SOCKS 代理的 HTTP 传输
		httpTransport := &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
		httpClient = &http.Client{Transport: httpTransport}
	} else {
		// 没有代理
		httpClient = http.DefaultClient
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHub.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	tc.Transport = httpClient.Transport

	return github.NewClient(tc)
}

func downloadFileIfNotExists(client *github.Client, owner, repo, path, sha, destDir string) error {
	localPath := filepath.Join(destDir, path)

	// 检查本地文件是否存在且内容是否相同
	if fileExistsAndEqual(localPath, sha) {
		return nil
	}

	return downloadFile(client, owner, repo, path, sha, destDir)
}

func fileExistsAndEqual(filePath, sha string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		// 文件不存在
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return false
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Printf("Error computing hash for file %s: %v\n", filePath, err)
		return false
	}

	return fmt.Sprintf("%x", hash.Sum(nil)) == sha
}

func downloadFile(client *github.Client, owner, repo, path, sha, destDir string) error {
	ctx := context.Background()

	// 下载文件内容
	content, _, err := client.Repositories.DownloadContents(ctx, owner, repo, path, nil)
	if err != nil {
		return fmt.Errorf("Error downloading file: %v", err)
	}
	defer content.Close()

	// 创建文件夹结构
	localPath := filepath.Join(destDir, path)
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}

	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	// 保存文件内容到本地文件
	if _, err := io.Copy(file, content); err != nil {
		return fmt.Errorf("Error saving file: %v", err)
	}

	return nil
}
