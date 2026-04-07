# 02 — GitHub Actions

## 目錄

- [事前準備](#事前準備)
- [先跑一個 Workflow](#先跑一個-workflow)
- [GitHub Actions 介面](#github-actions-介面)
- [回頭看：剛剛的 Workflow 在做什麼](#回頭看剛剛的-workflow-在做什麼)
- [手動觸發與失敗處理](#手動觸發與失敗處理)
- [常見問題排解](#常見問題排解)
- [小結](#小結)
- [參考手冊](02-reference.md) — Events、Actions、Context、環境變數、Runners、YAML 語法


## 事前準備

### 確認你的 GitHub 帳號

請確認你已經擁有一個 GitHub 帳號。如果還沒有，請前往 [github.com](https://github.com) 註冊一個。

### 建立一個新的 GitHub Repository

我們需要一個全新的 repository 來練習。請依照以下步驟操作：

1. 登入 GitHub 後，點擊右上角的 **「+」** 按鈕，選擇 **「New repository」**
2. 填寫 repository 資訊：
   - **Repository name**：輸入 `github-actions-lab`
   - **Description**（選填）：輸入 `My first GitHub Actions workflow`
   - **Visibility**：選擇 **Public**
3. 點擊 **「Create repository」** 按鈕

### 將 Repository Clone 到本機

打開終端機，執行以下指令：

```bash
# Replace <your-username> with your GitHub username
git clone https://github.com/<your-username>/github-actions-lab.git
cd github-actions-lab
```


## 先跑一個 Workflow

先不用管每一行在幹嘛，照著貼就好，我們等一下會回來解釋。

### 建立目錄結構

GitHub Actions 的 workflow 檔案必須放在 `.github/workflows/` 目錄下。我們先來建立這個目錄：

```bash
mkdir -p .github/workflows
```

### 建立 Workflow 檔案

在 `.github/workflows/` 目錄下建立一個名為 `hello.yml` 的檔案：

```bash
# You can use any editor you prefer
vim .github/workflows/hello.yml
```

將以下內容貼入編輯器，完成後打 `:wq` 離開 vim：

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
        run: echo "Hello, GitHub Actions! 🚀"

      - name: Show System Info
        run: |
          echo "OS: $(uname -s)"
          echo "Architecture: $(uname -m)"
          echo "Go Version: $(go version)"
          echo "Docker Version: $(docker --version)"
          echo "Current Time: $(date)"

      - name: Show GitHub Context
        run: |
          echo "Repository: ${{ github.repository }}"
          echo "Branch: ${{ github.ref_name }}"
          echo "Commit SHA: ${{ github.sha }}"
          echo "Actor: ${{ github.actor }}"
          echo "Event: ${{ github.event_name }}"
```

### Push 到 GitHub

```bash
git add .github/workflows/hello.yml
git commit -m "Add first GitHub Actions workflow"
git push origin main
```


## GitHub Actions 介面

打開瀏覽器，進入你的 GitHub repository 頁面，點擊上方的 **「Actions」** 分頁（如果看不到，請確認檔案已成功 push）。你應該會看到一個正在執行（或已完成）的 workflow run，名稱就是你剛剛的 commit message：

![GitHub Actions 介面](assets/github-actions-tab.png)

點進去之後，你會看到這次 workflow run 的詳細資訊：

![Workflow Run 頁面](assets/github-actions-workflow-run.png)

- **Summary**：這次 run 的摘要，包括觸發方式、commit SHA、狀態、執行時間
- **Job 列表**：左側會顯示這次 run 包含的所有 job（我們只有一個 `greeting` job）

點開 **greeting** job，你會看到這個 job 裡每個 step 的執行紀錄：

![Job Steps 列表](assets/github-actions-job-steps.png)

每個步驟旁邊都有打勾符號代表執行成功。點擊任何一個步驟可以展開查看詳細的輸出內容：

![Step Logs 展開畫面](assets/github-actions-step-logs.png)

試著看一下 **Show System Info**，你會發現 Runner 上已經預裝了 Go 和 Docker，並可以從中得知他們的版本。

## 回頭看：剛剛的 Workflow 在做什麼

OK，你已經成功跑完第一個 workflow 了。我們將在這個章節中，了解其中出現的重要概念。

### 1. Workflow

**Workflow** 是一個定義自動化流程的 **YAML 檔案**。剛剛那個 `hello.yml` 就是一個 workflow。

- 檔案位置：必須放在 repository 的 **`.github/workflows/`** 目錄下
- 副檔名：`.yml` 或 `.yaml`
- 一個 repository 可以有 **多個** workflow

```bash
my-repo/
├── .github/
│   └── workflows/
│       ├── ci.yml # workflow 1    
│       ├── pr-check.yml # workflow 2 
│       └── release.yml # workflow 3
├── main.go
└── go.mod
```

### 2. Event

**Event** 是觸發 workflow 執行的條件。在剛剛的 `hello.yml` 裡，我們定義了兩個 event：

```yaml
on:
  push:
    branches: [main]
  workflow_dispatch:
```

- `push: branches: [main]` — 當有程式碼被 push 到 `main` 分支時自動觸發
- `workflow_dispatch` — 允許在 GitHub UI 上手動觸發

常見事件還包括：

- 建立或更新 Pull Request（`pull_request`）
- 定時排程（`schedule`）

### 3. Job

**Job** 是 workflow 中的 **執行單元**。在剛剛的 `hello.yml` 裡，我們只有一個 job 叫做 `greeting`：

```yaml
jobs:
  greeting:
    runs-on: ubuntu-latest
```

一個 workflow 可以有多個 job，預設會平行執行。在下一章我們會看到一個 workflow 裡面同時有多個 job 的例子。

### 4. Step

**Step** 就是 job 裡面的執行步驟。在剛剛的 `hello.yml` 裡，`greeting` job 有三個 step：`Say Hello`、`Show System Info`、`Show GitHub Context`，他們會依序執行（同一個 job 內的 step 不能平行運作），並共用相同的檔案系統。在 每一個 step，你可以使用 `run` 來執行 shell 指令，以及 `uses` 引用別人已經寫好的 Action。

### 5. Action

**Action** 是別人已經寫好的可重用步驟，例如把程式碼 checkout 下來、安裝 Go 環境這些常見的步驟。社群都有現成的 Action 可以用，在 `hello.yml` 裡我們還沒用到，下一章開始會大量使用。

### 6. Runner

**Runner** 是實際執行 job 的 **機器環境**。在剛剛的 `hello.yml` 裡，`runs-on: ubuntu-latest` 代表我們使用 GitHub 提供的 Ubuntu 虛擬機。

- **GitHub-hosted Runner**：GitHub 提供的雲端虛擬機，用完就會自動釋放資源。使用 GitHub-hosted Runner 的好處是自己不需要管理伺服器，並且有 Ubuntu、Windows、macOS 等作業系統的選擇。
- **Self-hosted Runner**：自己管理的機器，適用於需要特殊硬體、軟體或網路環境的場景，但需要自行維護與更新。

### 概念關係圖

用剛剛的 `hello.yml` 來看這些概念的階層關係：

```
Workflow (.github/workflows/hello.yml)
├── Event
│   ├── on: push (branches: [main])
│   └── on: workflow_dispatch
├── Job: greeting
│   ├── Runner: ubuntu-latest
│   ├── Step 1 — run: echo "Hello, GitHub Actions! 🚀"
│   ├── Step 2 — run: echo "OS: ..." (Show System Info)
│   └── Step 3 — run: echo "Repository: ..." (Show GitHub Context)
```

## 手動觸發與失敗處理

### 手動觸發 Workflow

#### 什麼是 `workflow_dispatch`？

還記得我們在 `on:` 中定義了 `workflow_dispatch` 嗎？這個觸發器允許你**不需要 push 程式碼**，就能直接在 GitHub UI 上手動觸發 workflow，常用於手動重新執行失敗的 pipeline、測試 workflow 設定是否正確，或是觸發不需要程式碼變更的維運作業（例如手動清理暫存檔案）。

#### 手動觸發步驟

1. 前往你的 repository 的 **Actions** 分頁
2. 在左側邊欄找到 **「Hello GitHub Actions」** workflow
3. 點擊它，你會看到右上方出現一個 **「Run workflow」** 按鈕
4. 點擊 **「Run workflow」**，選擇 branch 為 `main`，然後按下綠色的 **「Run workflow」** 按鈕
5. 等待幾秒鐘，頁面會出現新的 workflow run

觀察這次手動觸發的結果，點開 **Show GitHub Context** 步驟，你會發現 `Event` 欄位顯示的是 `workflow_dispatch`，而不是上次的 `push`。

### 製造一個失敗

#### 為什麼要故意製造失敗？

在實際開發中，CI pipeline 失敗是常態。了解失敗時的表現可以幫助你快速排除問題。

#### 修改 Workflow

編輯 `.github/workflows/hello.yml`，在 `steps` 最後面加入一個 **一定會失敗** 的步驟：

```yaml
      - name: This Step Will Fail
        run: |
          echo "About to fail..."
          exit 1
```

`exit 1` 代表以 **非零的退出碼** 結束，shell 會將其視為錯誤。

#### Push 並觀察失敗結果

```bash
git add .github/workflows/hello.yml
git commit -m "Add a failing step to test error handling"
git push origin main
```

前往 Actions 分頁，你會看到：

- Workflow run 旁邊會出現一個 **紅色叉叉** ❌
- 點進去查看，前面的步驟仍然是 ✅，但最後一個步驟會顯示 ❌
- 點開失敗的步驟，可以看到錯誤訊息：`Error: Process completed with exit code 1.`

**重要觀念**：當某個 step 失敗時，**後續的 step 預設不會執行**，整個 job 會被標記為失敗。

#### 修復並重新 Push

了解失敗的表現後，讓我們把失敗的步驟移除（或註解掉），恢復正常：

```yaml
      # - name: This Step Will Fail
      #   run: |
      #     echo "About to fail..."
      #     exit 1
```

或者直接刪除那個步驟，然後 push：

```bash
git add .github/workflows/hello.yml
git commit -m "Remove failing step"
git push origin main
```

回到 Actions 分頁，確認新的 workflow run 又回到 **綠色勾勾** ✅。


> 更多 Events、Actions、Context、環境變數、Runner 規格、YAML 語法等詳細資訊，請參考 [GitHub Actions 參考手冊](02-reference.md)。


## 常見問題排解

在第一次撰寫 workflow 時，以下是最常見的幾個問題：

### 1. YAML 縮排錯誤

YAML 對縮排非常敏感，必須使用 **空格（space）**，**不能使用 Tab**。

**常見錯誤訊息**：

```
Invalid workflow file: .github/workflows/hello.yml
  - yaml: line X: mapping values are not allowed in this context
```

**排解方式**：

- 確認你的編輯器設定為使用空格縮排（建議 2 個空格）
- 在 VS Code 中，右下角可以看到目前的縮排設定（`Spaces: 2`），點擊可以切換
- 使用線上 YAML 驗證工具（如 [yamllint.com](https://www.yamllint.com/)）檢查語法

```yaml
# Wrong — mixed indentation
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Test
       run: echo "wrong indent"    # ← only 1 space, should be 2

# Correct
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Test
        run: echo "correct indent"  # ← 2 spaces from the dash
```

### 2. 檔案路徑錯誤

Workflow 檔案 **必須** 放在 `.github/workflows/` 目錄下，路徑大小寫也要正確。

**常見錯誤**：

| 錯誤路徑 | 說明 |
|----------|------|
| `github/workflows/hello.yml` | 少了前面的「.」 |
| `.github/workflow/hello.yml` | `workflows` 少了 `s` |
| `.Github/Workflows/hello.yml` | 大小寫錯誤 |
| `.github/workflows/hello.yaml` | 副檔名 `.yaml` 其實可以，但要確認一致 |

### 3. 權限問題

如果你的 repository 是新建的或剛 fork 的，可能需要確認 Actions 權限：

1. 前往 repository 的 **Settings** → **Actions** → **General**
2. 在 **Actions permissions** 區塊，選擇 **「Allow all actions and reusable workflows」**
3. 在 **Workflow permissions** 區塊，確認權限設定為 **「Read and write permissions」**（後面章節會用到寫入權限）

### 4. 分支名稱不符

如果你的 repository 預設分支是 `master` 而非 `main`，workflow 的 `branches: [main]` 不會被觸發。

**排解方式**：確認你的預設分支名稱，並修改 workflow 中的 `branches` 設定：

```yaml
on:
  push:
    branches: [master]  # ← change to match your default branch
```


## 小結

- **GitHub Actions** 是 GitHub 內建的 CI/CD 平台，用 YAML 檔案定義自動化流程
- 六個核心概念：**Workflow → Event → Job → Step → Action → Runner**
- Workflow 檔案放在 **`.github/workflows/`** 目錄下
- 透過 **Event** 決定何時觸發，透過 **Job** 與 **Step** 定義要做什麼
- `workflow_dispatch` 讓你可以 **手動觸發** workflow
- 當 step 失敗（exit code 非零）時，**後續步驟不會執行**，job 標記為失敗
- 更多 Events、Actions、Context、環境變數等詳細資訊請參考 [參考手冊](02-reference.md)

### 練習題

完成以下練習來鞏固本章所學：

👉 [練習一：GitHub Actions 基礎練習](exercises/exercise-01-basics.md)（練習 1-1 至 1-3）


[← 上一章：CI/CD 概念介紹](01-cicd-intro.md) ｜ [下一章：Go 專案 CI Pipeline →](03-go-ci-pipeline.md)
