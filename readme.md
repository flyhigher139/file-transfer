# File Transfer

这是一个用 Go 语言编写的文件传输工具，基于 Gin 框架，提供了一个简单而强大的文件上传和下载服务，并支持断点续传功能。

## ✨ 功能特性

- **文件上传**：支持普通文件上传和分块上传。
- **文件下载**：支持普通文件下载和带 `Range` 请求的断点续传下载。
- **前端界面**：提供一个基于 Bootstrap 的美观易用的前端页面，用于文件操作。
- **RESTful API**：提供清晰的 API 接口。
- **单元测试**：包含使用 `testify` 编写的完整单元测试。

## 🚀 快速开始

### 环境要求

- Go 1.18 或更高版本

### 安装

1. 克隆仓库:
   ```sh
   git clone https://github.com/your-username/file-transfer.git
   cd file-transfer
   ```

2. 安装依赖:
   ```sh
   go mod tidy
   ```

### 运行

你可以通过以下命令启动服务：

```sh
go run main.go --storage my_files
```

- `--storage` 参数指定了文件存储的目录。如果目录不存在，程序会自动创建。

服务启动后，你可以在浏览器中打开 `http://localhost:8080/static/index.html` 来访问前端页面。

## API 文档

### 获取文件列表

- **URL**: `/api/files`
- **Method**: `GET`
- **Success Response**:
  ```json
  {
    "files": ["file1.txt", "file2.jpg"]
  }
  ```

### 普通上传

- **URL**: `/api/simple/upload`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Form-Data**:
  - `file`: 要上传的文件

### 普通下载

- **URL**: `/api/simple/download/:filename`
- **Method**: `GET`

### 分块上传

1.  **上传文件块**
    - **URL**: `/api/files`
    - **Method**: `POST`
    - **Content-Type**: `multipart/form-data`
    - **Form-Data**:
      - `file`: 文件块
      - `chunkNumber`: 当前块的编号 (从 1 开始)
      - `totalChunks`: 总块数
      - `fileIdentifier`: 文件的唯一标识符

2.  **合并文件块**
    - **URL**: `/api/files/merge`
    - **Method**: `POST`
    - **Content-Type**: `multipart/form-data`
    - **Form-Data**:
      - `filename`: 原始文件名
      - `totalChunks`: 总块数
      - `fileIdentifier`: 文件的唯一标识符

### 断点续传下载

- **URL**: `/api/files/:filename`
- **Method**: `GET`
- **Headers**:
  - `Range`: `bytes={start}-{end}`

## ✅ 测试

运行单元测试，请执行以下命令：

```sh
go test -v ./...
```