# 練習 1：執行你的第一個容器

> 搭配章節：[01-docker-concepts.md](../01-docker-concepts.md)

Nginx 是一個常見的網頁伺服器，很多網站背後都用它來處理請求。這次的目標是啟動一個 Nginx 容器，把自己寫的 HTML 頁面透過 Bind Mount 掛進去，從瀏覽器看到結果。

## Step 1：建立工作目錄和網頁

為這次練習開一個乾淨的工作目錄：

```bash
mkdir -p ~/docker-lab/lab1
cd ~/docker-lab/lab1
```

接著在目錄下建立 `index.html`，等下會把這個檔案掛進容器讓 Nginx 去 serve：

```bash
cat << 'EOF' > index.html
<!DOCTYPE html>
<html>
<head>
    <title>Docker Lab 1</title>
</head>
<body>
    <h1>Hello from Docker!</h1>
    <p>This page is served from a Docker container.</p>
</body>
</html>
EOF
```

> `cat << 'EOF' > index.html` 是 shell 的 heredoc 語法，意思是「把到下一行 `EOF` 為止的內容寫進 `index.html`」。也可以直接用編輯器開檔案寫。

## Step 2：啟動 Nginx 容器

```bash
docker run -d --name lab1-nginx -p 8080:80 -v $(pwd):/usr/share/nginx/html nginx:1.27-alpine
```

> - `-d`：背景執行，不佔用當前 terminal。
> - `--name lab1-nginx`：給容器一個好記的名字，後面 `stop`、`rm` 就不用查 Container ID。
> - `-p 8080:80`：把 host 的 `8080` 轉發到容器內 Nginx 監聽的 `80`。
> - `-v $(pwd):/usr/share/nginx/html`：把當前目錄掛到 Nginx 預設放靜態檔案的位置，這樣 `index.html` 就會被 serve 出去。

## Step 3：確認容器在跑，然後測試

```bash
docker ps
```

應該會看到 `lab1-nginx` 這一列，狀態是 `Up`、PORTS 欄有 `0.0.0.0:8080->80/tcp`。接著打開瀏覽器連 `http://localhost:8080`，應該就能看到剛才寫的 HTML 頁面。

## Step 4：試試除錯指令

看容器的 log（每次有請求進來 Nginx 就會寫一行 access log）：

```bash
docker logs lab1-nginx
```

進到容器裡看看環境：

```bash
docker exec -it lab1-nginx sh
```

> 進去之後可以用 `ls /usr/share/nginx/html` 看到剛剛掛進去的 `index.html`。打 `exit` 離開容器。

## Step 5：清理

容器用完要記得清掉，不然會一直留在電腦上：

```bash
docker stop lab1-nginx
docker rm lab1-nginx
```

> `docker stop` 只是讓容器停止運作，容器本體還在；真的要讓它消失要 `docker rm`。想一步到位可以用 `docker rm -f lab1-nginx` 強制刪除執行中的容器。

## 延伸練習

1. 修改 `index.html` 的內容，不重啟容器，直接重新整理瀏覽器查看變化（思考：為什麼修改會立即生效？）

Snow 看著 Ocean 的練習成果，點了點頭：「不錯嘛，基本操作都會了。剛好 SDC 缺了一個 SRE Team，你要不要來做做看？」

Ocean：「...... 那要做什麼？」

Snow：「用 Docker Compose 把社團服務們的設定步驟寫成檔案。」

Ocean：「那是啥？」

詳情請見下一章：[02-docker-compose.md](02-docker-compose.md)
