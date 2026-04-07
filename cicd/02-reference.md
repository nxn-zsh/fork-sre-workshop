# 02 — GitHub Actions 參考手冊

本頁整理了 GitHub Actions 常用的 Events、Actions、Context、環境變數、Runner 規格、免費額度、YAML 語法和 GITHUB_TOKEN 等資訊，供需要時查閱。

## 目錄

- [常見 Events（觸發器）](#常見-events觸發器)
- [常見 Actions](#常見-actions)
- [Context 與表達式](#context-與表達式)
- [環境變數](#環境變數)
- [GitHub-hosted Runners](#github-hosted-runners)
- [免費額度與限制](#免費額度與限制)
- [YAML 語法速查](#yaml-語法速查)
- [GITHUB_TOKEN](#github_token)

## 常見 Events（觸發器）

| Event | 觸發時機 | 範例 |
|-------|---------|------|
| `push` | 推送程式碼到指定分支時 | `on: push` / `on: push: branches: [main]` |
| `pull_request` | 建立或更新 PR 時 | `on: pull_request: branches: [main]` |
| `schedule` | 定時排程（cron 語法） | `on: schedule: - cron: '0 0 * * *'` |
| `workflow_dispatch` | 手動觸發（在 Actions 頁面按按鈕） | `on: workflow_dispatch` |
| `release` | 建立 Release 時 | `on: release: types: [published]` |
| `issues` | Issue 被建立或修改時 | `on: issues: types: [opened]` |
| `workflow_run` | 另一個 workflow 完成後觸發 | `on: workflow_run: workflows: [CI]` |

### 事件篩選

你可以進一步篩選觸發條件，只在特定情況下才執行：

```yaml
on:
  push:
    # Only trigger on specific branches
    branches: [main, develop]
    # Only trigger when specific files change
    paths:
      - 'src/**'
      - '*.go'
    # Ignore specific paths
    paths-ignore:
      - 'docs/**'
      - '*.md'

  pull_request:
    # Only trigger on specific event types
    types: [opened, synchronize, reopened]
    branches: [main]
```

## 常見 Actions

以下是你會經常用到的官方 Actions：

### actions/checkout

**用途**：將 repository 的程式碼下載到 Runner 上。幾乎每個 workflow 的第一步都是這個。

```yaml
- uses: actions/checkout@v4
```

### actions/setup-go

**用途**：安裝指定版本的 Go 語言環境。

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.24'
```

### actions/cache

**用途**：快取依賴套件，加速後續的建置。

```yaml
- uses: actions/cache@v4
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

### actions/upload-artifact

**用途**：將建置產出的檔案上傳為 artifact，供下載或傳給其他 job 使用。

```yaml
- uses: actions/upload-artifact@v4
  with:
    name: my-binary
    path: ./build/app
```

### actions/download-artifact

**用途**：下載先前上傳的 artifact（通常在不同的 job 中使用）。

```yaml
- uses: actions/download-artifact@v4
  with:
    name: my-binary
    path: ./build/
```

### 如何尋找更多 Actions？

前往 [GitHub Marketplace](https://github.com/marketplace?type=actions) 搜尋你需要的功能。在搜尋時，注意以下幾點：

- 優先選擇有 **Verified Creator** 標章的 Action
- 查看 **星星數** 和 **使用人數**
- 確認 Action 的 **最後更新時間**，避免使用已停止維護的 Action
- 閱讀 Action 的 README，了解所有可用參數


## Context 與表達式

### `${{ }}` 表達式語法

在 GitHub Actions 中，`${{ }}` 是一種特殊語法，用來存取 **context 物件** 的資訊。GitHub 在 workflow 執行時會自動提供這些資訊。

```yaml
# Basic usage
run: echo "Hello, ${{ github.actor }}!"

# Use in conditionals
if: ${{ github.ref == 'refs/heads/main' }}

# Use in environment variables
env:
  REPO_NAME: ${{ github.repository }}
```

### 常用的 Context

#### `github` context

提供 workflow run 的相關資訊，是最常用的 context：

| 表達式 | 說明 | 範例值 |
|--------|------|--------|
| `github.repository` | Repository 全名 | `octocat/hello-world` |
| `github.ref` | 完整的 Git ref | `refs/heads/main` |
| `github.ref_name` | 分支或 tag 名稱 | `main` |
| `github.sha` | 觸發 commit 的完整 SHA | `abc123def456...` |
| `github.actor` | 觸發事件的使用者 | `octocat` |
| `github.event_name` | 觸發的事件名稱 | `push`、`pull_request` |
| `github.run_number` | Workflow 的執行序號 | `42` |
| `github.workspace` | Runner 上的工作目錄路徑 | `/home/runner/work/repo/repo` |

#### `env` context

用來存取在 workflow 中定義的環境變數：

```yaml
env:
  MY_VAR: hello

steps:
  - run: echo "${{ env.MY_VAR }}"
```

#### `secrets` context

用來存取在 repository 中設定的 **機密資訊**（如 API Key、Token 等）。Secrets 的內容在 log 中會自動被遮蔽，不會外洩。

```yaml
steps:
  - name: Deploy
    run: ./deploy.sh
    env:
      API_KEY: ${{ secrets.API_KEY }}
```

#### `runner` context

提供 Runner 環境的相關資訊：

| 表達式 | 說明 | 範例值 |
|--------|------|--------|
| `runner.os` | 作業系統 | `Linux`、`Windows`、`macOS` |
| `runner.arch` | 架構 | `X64`、`ARM64` |
| `runner.temp` | 暫存目錄路徑 | `/home/runner/work/_temp` |


## 環境變數

### 環境變數的三個層級

在 GitHub Actions 中，你可以在三個不同的層級設定環境變數，每個層級的 **作用範圍** 不同：

```yaml
name: Environment Variables Demo

on:
  push:
    branches: [main]

# Workflow-level: available to ALL jobs and ALL steps
env:
  APP_NAME: my-awesome-app
  ENVIRONMENT: production

jobs:
  demo:
    runs-on: ubuntu-latest

    # Job-level: available to ALL steps in THIS job only
    env:
      LOG_LEVEL: debug
      DATABASE_HOST: localhost

    steps:
      - name: Show All Variables
        # Step-level: available to THIS step only
        env:
          STEP_VAR: "I only exist in this step"
        run: |
          echo "App Name: $APP_NAME"
          echo "Environment: $ENVIRONMENT"
          echo "Log Level: $LOG_LEVEL"
          echo "Database Host: $DATABASE_HOST"
          echo "Step Var: $STEP_VAR"

      - name: Step Var Is Gone
        run: |
          echo "App Name: $APP_NAME"
          echo "Log Level: $LOG_LEVEL"
          echo "Step Var: $STEP_VAR"
          # STEP_VAR will be empty here because it was defined in the previous step
```

### 三個層級的作用範圍比較

| 層級 | 設定位置 | 作用範圍 |
|------|---------|---------|
| **Workflow level** | `on:` 同層的 `env:` | 所有 job 的所有 step 都能使用 |
| **Job level** | `jobs.<job_id>:` 底下的 `env:` | 該 job 內的所有 step 能使用 |
| **Step level** | `steps[*]:` 底下的 `env:` | 僅限該 step 能使用 |

**優先順序**：如果不同層級定義了相同名稱的環境變數，**越小範圍的優先**（Step > Job > Workflow）。

### 範例：環境變數的覆蓋

```yaml
env:
  MESSAGE: "I'm from workflow level"

jobs:
  demo:
    runs-on: ubuntu-latest
    env:
      MESSAGE: "I'm from job level"
    steps:
      - name: Check Message
        env:
          MESSAGE: "I'm from step level"
        run: echo "$MESSAGE"
        # Output: I'm from step level
```


## GitHub-hosted Runners

GitHub 提供免費的雲端 Runner，不需要自行架設伺服器。

### 可用的作業系統

| Runner 標籤 | 作業系統 | 說明 |
|-------------|---------|------|
| `ubuntu-latest` | Ubuntu (最新 LTS) | 最常用，建議預設選擇 |
| `ubuntu-24.04` | Ubuntu 24.04 | 指定特定 Ubuntu 版本 |
| `windows-latest` | Windows Server | 需要 Windows 環境時使用 |
| `macos-latest` | macOS | 需要 macOS 環境時使用（例如 iOS 開發） |

> `*-latest` 標籤對應的實際版本會隨時間更新，請參考 [GitHub 官方文件](https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners) 確認目前的版本。

### 預裝的軟體

GitHub-hosted Runner 已預裝大量常用軟體，包括：

- **語言**：Go、Node.js、Python、Java、Ruby、Rust 等
- **套件管理器**：npm、pip、gem 等
- **工具**：Git、Docker、docker compose、curl、wget、jq 等
- **雲端 CLI**：AWS CLI、Azure CLI、Google Cloud SDK 等

> 完整清單請參考 [GitHub 官方文件](https://github.com/actions/runner-images)。

### 資源限制

| 項目 | 限制 |
|------|------|
| vCPU | 2 核（Linux/Windows）/ 3 核（macOS） |
| 記憶體 | 7 GB（Linux/Windows）/ 14 GB（macOS） |
| 磁碟空間 | 14 GB（SSD） |


## 免費額度與限制

### 費用方案

| 方案 | 免費額度 |
|------|---------|
| **公開（Public）Repository** | 完全免費，無分鐘數限制 |
| **私人（Private）Repository — Free 方案** | 每月 2,000 分鐘 |
| **私人（Private）Repository — Pro 方案** | 每月 3,000 分鐘 |
| **私人（Private）Repository — Team 方案** | 每月 3,000 分鐘 |

> **注意**：不同作業系統的分鐘數消耗倍率不同。Linux = 1x、Windows = 2x、macOS = 10x。

### 執行時間限制

| 項目 | 限制 |
|------|------|
| 單次 job 最長執行時間 | **6 小時** |
| 單次 workflow 最長執行時間 | **35 天**（通常用於含有 approval 等待的 workflow） |
| API 請求頻率 | 每個 repository 每小時 1,000 次 |
| 同時執行的 job 數量 | Free 方案最多 20 個 |

## YAML 語法速查

GitHub Actions 的 workflow 使用 YAML 格式。如果你不熟悉 YAML，以下是快速入門：

### 縮排規則

YAML 使用 **空格**（space）進行縮排，**不能用 Tab**。通常使用 **2 個空格** 為一層縮排。

```yaml
# Good — using spaces
parent:
  child:
    grandchild: value

# Bad — using tabs (will cause errors!)
parent:
	child:           # ← This is a tab, YAML will reject this
```

### Key-Value Pairs（鍵值對）

```yaml
name: My Workflow
version: 1.0
enabled: true
count: 42
```

### 列表（List）

```yaml
# Block style
fruits:
  - apple
  - banana
  - cherry

# Inline style
fruits: [apple, banana, cherry]
```

### 巢狀結構

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Build
        run: go build ./...
```

### 多行字串

```yaml
# Literal block (preserves newlines)
description: |
  This is line 1.
  This is line 2.
  This is line 3.

# Folded block (joins lines with spaces)
description: >
  This is a long sentence
  that will be joined
  into a single line.
```


## GITHUB_TOKEN

`secrets.GITHUB_TOKEN` 是 GitHub 在每次 workflow 執行時自動產生的臨時 token，不需要你手動設定。它的預設權限取決於 repository 的設定，但你可以在 workflow 中用 `permissions` 關鍵字明確限縮權限。

重點：

- 這個 token 在 workflow 結束後就會失效
- 它的權限只限於觸發 workflow 的 repository
- 如果你需要存取其他 repository，就需要使用 Personal Access Token (PAT)
- 建議在每個 workflow 中都明確設定 `permissions`，遵循最小權限原則


[← 回到 GitHub Actions 基礎](02-github-actions-basics.md)
