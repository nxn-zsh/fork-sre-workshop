# 03 — Dockerfile 基礎

## 3.1 什麼是 Dockerfile？

Ocean 用 Compose 把服務跑得有聲有色，但有一天 Andrew 推了新的 Go 程式碼，Ocean 發現 Docker Hub 上根本沒有現成的映像檔可以用。Andrew：「別人的映像檔當然沒有我們的程式碼。要把自己的程式打包成映像檔，就要寫 Dockerfile。」

前面用的 `nginx`、`postgres` 都是別人做好放在 Docker Hub 上的 Image。要把自己寫的程式變成 Image，就要寫一個 **Dockerfile**，一個放在專案根目錄的純文字檔，描述這個 Image 要怎麼建出來：基礎環境、要複製哪些檔案、要執行什麼安裝指令、程式怎麼啟動。寫好之後用 `docker build` 把它變成 Image，再用 `docker run` 跑起來。

---

## 3.2 第一個 Dockerfile

概念講完了，先用最簡單的例子走一次完整流程。建一個空目錄並切進去：

```bash
mkdir hello-docker && cd hello-docker
```

在裡面創造一個檔名叫 `Dockerfile` 的檔案：

```dockerfile
FROM alpine:3.21
CMD ["echo", "Hello from my first Dockerfile!"]
```

在這個目錄底下執行 build 指令，建立 Image：

```bash
docker build -t hello .
```

這會產出一個名為 `hello` 的 Image。用 `docker images` 接上名稱當作 filter，可以只列出這張 Image，確認它真的建出來了：

```bash
docker images hello
```

輸出會像這樣：

```
IMAGE          ID             DISK USAGE   CONTENT SIZE   EXTRA
hello:latest   c2de8808511c       8.17MB             0B
```

接著把它跑起來：

```bash
docker run --rm hello
```

你會看到：

```
Hello from my first Dockerfile!
```

---

## 3.3 Go 網頁服務的 Dockerfile

`exercises/my-app/` 準備好了一個最簡單的 Go 專案，結構長這樣：

```
exercises/my-app/
├── main.go         # Go 程式碼
├── go.mod          # Go 模組定義
└── Dockerfile      # 建置映像檔的腳本
```

`go.mod` 宣告模組名稱為 `my-app`、Go 版本為 `1.24`，只用了標準函式庫。

`main.go` 是一個簡單的 HTTP Server，監聽 container `8080` port，提供兩個 endpoint：

- `GET /`：回傳 `Hello from Go! Hostname: <hostname>`
- `GET /health`：回傳 `OK`，用於健康檢查

### Dockerfile

如果不用 Docker，要在一台新的機器上把這個 Go 服務跑起來，你大概會打這些指令：

```bash
brew install go                 # 1. 安裝 Go
mkdir /app && cd /app           # 2. 切到工作目錄
git clone <repo-url> .          # 3. 把程式碼放進去
go build -o server .            # 4. 編譯
./server                        # 5. 執行
```

Dockerfile 把這五件事寫成指令，讓 Docker 在一個乾淨的環境裡照著做：

```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o server .
CMD ["./server"]
```

五個指令剛好對應上面五個步驟。一行一行來看：

**`FROM golang:1.24-alpine`** — 每個 Dockerfile 都從 `FROM` 開始，指定你要在什麼基礎環境上建置。這裡用的是已經裝好 Go 1.24 的 Alpine Linux。

**`WORKDIR /app`** — 設定映像檔內部的工作目錄。後面的 `COPY`、`RUN` 都會在 `/app` 這個目錄下執行。

**`COPY . .`** — 把你本機目前目錄下的所有檔案複製到映像檔的 `/app` 裡面。第一個 `.` 是本機路徑，第二個 `.` 是容器內的路徑（也就是 `WORKDIR` 設定的 `/app`）。

**`RUN go build -o server .`** — 在建置過程中執行命令。這裡是編譯 Go 程式，產出一個叫 `server` 的執行檔。`RUN` 的結果會被寫進映像檔的 Layer 裡。

**`CMD ["./server"]`** — 容器啟動時的預設命令。當你 `docker run` 這個映像檔的時候，它就會執行 `./server`。

所以我們可以用和 3.2 相同的方法，把這個服務 build 成一張 Image，切到 `exercises/my-app/` 目錄下建置並執行：

```bash
cd exercises/my-app
docker build -t my-app .         # 建置映像檔，取名叫 my-app
docker run -p 8080:8080 my-app   # 跑起來
```

跑完這些，打開瀏覽器連 `http://localhost:8080`，就會看到 `Hello from Go! Hostname: ...` 的回應。

### 指令速查表

上面的範例只用了五個指令，但 Dockerfile 還有其他常用指令。完整列表如下：

| 指令 | 說明 | 範例 |
|------|------|------|
| `FROM` | 指定基礎映像檔（必須是第一行） | `FROM golang:1.24-alpine` |
| `WORKDIR` | 設定工作目錄 | `WORKDIR /app` |
| `COPY` | 複製檔案至映像檔 | `COPY . .` |
| `RUN` | 建置時執行命令 | `RUN go build -o server .` |
| `CMD` | 容器啟動時的預設命令 | `CMD ["./server"]` |
| `ENTRYPOINT` | 容器啟動時的固定入口 | `ENTRYPOINT ["./server"]` |
| `EXPOSE` | 宣告容器監聽的連接埠（僅文件用途） | `EXPOSE 8080` |
| `ENV` | 設定環境變數 | `ENV GIN_MODE=release` |
| `USER` | 指定執行命令的使用者 | `USER nonroot` |
| `HEALTHCHECK` | 定義健康檢查 | `HEALTHCHECK CMD curl -f http://localhost/` |

基本上 AI 很會寫這種東西，對於 SRE 言這部分也不是最重要的，這裡就不展開來講。

專案開發越多，專案啟動的設定過程越長，Dockerfile 也會長得越複雜。它最大的價值是強制要求把專案的設定流程定義成規範化語言，讓啟動流程可以維持統一，並減少新進開發者的設定負擔。做完這個章節，我們就正式完成了 Dockerfile → `docker build` → Image → `docker run` → Container 的完整循環。

→ 綜合練習：[exercises/02-push-to-dockerhub.md](exercises/02-push-to-dockerhub.md)
