# 練習 3：把你的 Image 推上 Docker Hub

> 搭配章節：[03-dockerfile.md](../03-dockerfile.md)

Ocean 把 `my-app` build 成 Image 之後很得意，想請在外地實習的 Andrew 幫忙測試。Andrew：「我手上又沒有原始碼，怎麼 build？」Ocean：「那我把 Image 傳給你？」Andrew：「推 Docker Hub 啊，我直接 `docker run` 就好。」

這次的目標是把 3.3 建好的 `my-app` 推到 Docker Hub，然後模擬另一個人從 Docker Hub 拉回來執行。

## Prerequisite

到 [hub.docker.com](https://hub.docker.com) 註冊一個免費帳號，記住你的 username，或是直接把這份文件中的 <你的 username> 尋找並取代為你的 username。

## Step 1：確認本機已經有 my-app

3.3 應該已經 build 過 `my-app`，檢查一下：

```bash
docker images my-app
```

如果看不到，重新 build 一次：

```bash
cd Docker/exercises/my-app
docker build -t my-app .
```

## Step 2：登入 Docker Hub

```bash
docker login
```

完成登入流程會看到 `Login Succeeded`。

## Step 3：為 Image 加上帶 namespace 的 tag

Docker Hub 上的 Image 必須用 `<username>/<repository>:<tag>` 的格式，例如 `ocean1029/my-app:v1`。原本的 `my-app` 沒有 namespace，推不上去，要先重新 tag：

```bash
docker tag my-app <你的 username>/my-app:v1
```

> `docker tag` 不會產生新的 Image，只是幫同一張 Image 多加一個名字。用 `docker images` 可以看到兩個 tag 共用同一個 `IMAGE ID`。

## Step 4：Push 到 Docker Hub

```bash
docker push <你的 username>/my-app:v1
```

Push 完之後，打開瀏覽器連到 `https://hub.docker.com/r/<你的 username>/my-app`，會看到 Image 已經出現在你的 Docker Hub 頁面上：

![my-app 的 Docker Hub 頁面](../assets/docker-hub-my-app.png)

## Step 5：模擬從別台機器拉下來執行

先把本機的 Image 刪掉，模擬你身在外地、手上沒有任何本地 Image 的 Andrew：

```bash
docker rmi my-app <你的 username>/my-app:v1
docker images | grep my-app     # 應該什麼都沒有
```

> 如果跳 conflict error 說 container xxxxx is using ......，請先把 container 刪掉。

然後直接 `docker run`，Docker 會自動從 Docker Hub 拉回來再跑起來：

```bash
docker run --rm -p 8080:8080 <你的 username>/my-app:v1
```

打開 `http://localhost:8080` 應該看到 `Hello from Go! Hostname: ...`。整個流程：本機 build → push → 別人 pull → run，跟你前面用 `nginx`、`postgres` 的方式一模一樣，差別只在這次倉庫裡放的是你自己的服務。

## 延伸練習

1. 修改 `main.go` 的回傳訊息，重新 build 並推 `:v2`，觀察 Docker Hub 上多了一個版本。
