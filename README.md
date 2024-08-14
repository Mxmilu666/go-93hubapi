# Go-93hubapi

`go-93hubapi` 是一个用于分发来自 ~~群友们的奇怪言论~~ [bangbang93hub](https://github.com/Mxmilu666/Bangbang93hub) 仓库的图片文件的简单程序 基于 Golang

但是理论支持所有 Github 仓库

## ⚙️ 使用方法
### Releases 下载

TODO..

### 源码构建

1. 克隆项目仓库：
    ```bash
    git clone https://github.com/Mxmilu666/go-93hubapi.git
    ```

2. 进入项目目录并构建：
    ```bash
    cd go-93hubapi
    go build
    ```

3. 启动服务：
    ```bash
    ./go-93hubapi
    ```

4. 编写配置文件
   
5. 启动！

## 🛠️ 配置

在启动前，请根据需要在 `config.yaml` 文件中配置以下选项：

```bash
github:
  token: 你的 Github Token
  owner: Mxmilu666
  repo: bangbang93hub
server:
  address: localhost
  port: 8080
//你的缓存目录
dest: your_dest
//支持socks代理
ProxyURL: ""
maxconcurrent: 4

```
## 📄 提供的 API
### 下载文件

你可以通过以下 API 下载指定文件：

```bash
GET http://127.0.0.1:8080/download/file/{filename}
```

示例

```bash
GET http://127.0.0.1:8080/download/if%20err%20!=%20nil%20return%20nil,%20err.png
```

### 获取 FileList

你可以通过以下 API 获取 FileList：

```bash
GET http://127.0.0.1:8080/api/filelist
```

### 随机获取文件

你可以通过以下 API 获取随机文件：

TODO

## 🤝 贡献

欢迎贡献代码和建议！请提交 Pull Request 或创建 Issue

## 📄 许可证

项目采用 ``Apache-2.0 license``  协议开源
