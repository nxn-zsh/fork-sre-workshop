# 練習一：GitHub Actions 基礎練習

> **難度：** 入門 | **對應章節：** 02-GitHub Actions 基礎

---

## 目錄

- [練習 1-1：自訂 Hello World](#練習-1-1自訂-hello-world)
- [延伸思考](#延伸思考)

---

## 練習 1-1：自訂 Hello World

### 目標

修改 `.github/workflows/hello.yml`，加入你自己的自訂步驟，熟悉 GitHub Actions step 的基本寫法。

### 要求

1. 加入一個步驟，印出你自己的名字和學號
2. 加入一個步驟，使用 `actions/checkout@v4` 把 repo 的程式碼 checkout 下來，然後用 `ls -la` 列出工作目錄的檔案
3. 加入一個步驟，顯示目前的 Git 資訊（分支名稱、最近一次 commit 的訊息和作者）

### 提示

- 要能看到 repository 裡的檔案，**必須先使用 `actions/checkout@v4`** 來 checkout 程式碼。如果沒有 checkout，Runner 上的工作目錄是空的。
- 多行指令可以使用 `|`（literal block）語法。
- Git 資訊可以用 `git log -1` 來查看最近一次 commit。

### 預期結果

- 在 GitHub Actions 的 log 中可以看到你的名字和學號
- 可以看到 repo 中所有檔案的列表
- 可以看到 Git 分支名稱和 commit 資訊

<details>
<summary>點擊查看答案</summary>

```yaml
name: Hello GitHub Actions

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  greeting:
    runs-on: ubuntu-latest
    steps:
      - name: Say Hello
        run: echo "Hello, GitHub Actions!"

      - name: Show My Info
        run: |
          echo "Name: <your-name>"
          echo "Student ID: <your-student-id>"

      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: List Files in Working Directory
        run: ls -la

      - name: Show Git Info
        run: |
          echo "Branch: $(git branch --show-current)"
          echo "Latest Commit:"
          git log -1 --pretty=format:"  SHA: %H%n  Author: %an <%ae>%n  Date: %ad%n  Message: %s"
```

**重點說明：**

- `actions/checkout@v4` 必須放在 `ls -la` 和 `git log` **之前**，否則工作目錄是空的，也沒有 `.git` 資料夾。
- `git log -1` 的 `--pretty=format:` 可以自訂輸出格式，`%H` 是完整 SHA、`%an` 是作者名稱、`%s` 是 commit 訊息。

</details>

---

## 延伸思考

以下問題沒有標準答案，供進階學生思考和討論：

1. **環境變數安全性**：如果你需要在 workflow 中使用一個 API Key，應該用環境變數還是 GitHub Secrets？為什麼？兩者的差異是什麼？

2. **Workflow 除錯**：如果你的 workflow 一直失敗，你會用什麼方法來除錯？列出至少三種方法。

---

[回到教材目錄 →](../README.md)
