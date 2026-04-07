# 03 — Go 專案 CI Pipeline

Ocean 成功跑完第一個 Hello World workflow 後信心大增。Andrew 看到後問他：「那你能幫我的 Go 專案也設一個自動化測試嗎？每次 push 完都要手動跑 `go test`，我已經受夠了。」Ocean 想了想，決定動手試試。

## Table of Contents

- [學習目標](#學習目標)
- [範例專案介紹](#範例專案介紹)
- [完整 CI Workflow](#完整-ci-workflow)
- [逐段解說](#逐段解說)
- [PR 觸發的 CI](#pr-觸發的-ci)
- [常見問題排解](#常見問題排解)
- [小結與練習題](#小結與練習題)


## 學習目標

完成本章節後，你將能夠：

- 為 Go 專案建立一套包含 **lint、test、build** 的完整 CI pipeline
- 理解 **job 之間的依賴關係** 以及平行與序列執行的差異
- 知道如何用 **artifacts** 在 job 之間傳遞檔案
- 理解 `on: push` 和 `on: pull_request` 的差異


## 範例專案介紹

本章使用一個簡單的 Go HTTP API 伺服器作為範例，放在 `cicd/examples/sample-app/`：

```
cicd/examples/sample-app/
├── main.go            # HTTP server entry point
├── handler.go         # HTTP handler functions
├── handler_test.go    # Unit tests
├── go.mod             # Go module definition
├── Dockerfile         # Container image definition
└── .golangci.yml      # golangci-lint configuration
```

它提供三個 endpoint：

| Endpoint | Method | 功能 | 回傳範例 |
|----------|--------|------|---------|
| `/` | GET | 首頁歡迎訊息 | `{"message": "Hello, World!"}` |
| `/health` | GET | 健康檢查 | `{"status": "ok"}` |
| `/version` | GET | 版本資訊 | `{"version": "1.0.0"}` |

並在 `handler_test.go` 中用 Go 標準函式庫的 `testing` 和 `net/http/httptest` 寫了單元測試。程式碼細節可以直接打開 `cicd/examples/sample-app/` 查看，這章的重點是幫它建立 CI pipeline。

## 完整 CI Workflow

接下來我們要為範例專案建立一個完整的 Go CI workflow。這個 workflow 會在每次 `push` 或 `pull_request` 到 `main` 時自動觸發，並依序完成三件事：

1. **Lint**：用 `golangci-lint` 檢查程式碼風格與潛在問題
2. **Test**：執行單元測試，產生覆蓋率報告並上傳成 artifact
3. **Build**：在 lint 與 test 都通過後，才編譯出 binary 並上傳

其中 `lint` 和 `test` 會平行執行以節省時間，`build` 則透過 `needs` 等待前兩個 job 成功後才啟動。整個流程會是這樣：

```
        ┌──────┐
        │ push │
        └───┬──┘
     ┌──────┴──────┐
     ▼             ▼
  ┌──────┐     ┌──────┐
  │ lint │     │ test │
  └───┬──┘     └───┬──┘
     └──────┬──────┘
            ▼
        ┌───────┐
        │ build │
        └───────┘
```

請在你的專案中建立 `.github/workflows/ci.yml`：

```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Show coverage
        run: go tool cover -func=coverage.out
      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.out

  build:
    name: Build
    needs: [lint, test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build binary
        run: go build -o bin/app ./...
      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: app-binary
          path: bin/app
```

> 完整檔案也可以在 `cicd/examples/sample-app/.github/workflows/ci.yml` 找到，直接複製過去即可。

## 逐段解說

### 觸發條件

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

這個 workflow 會在 push 到 `main` 或有 PR 指向 `main` 時執行。`pull_request` 在合併前先檢查一次，`push` 作為合併後的最後防線。兩者的細節差異與進階設定（activity types、permissions、fork 限制）見下方「[PR 觸發的 CI](#pr-觸發的-ci)」。


### Lint Job — 程式碼品質檢查

```yaml
lint:
  name: Lint
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.24'
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v6
      with:
        version: latest
```

**Linter** 是一種靜態分析工具，它不會執行你的程式碼，而是檢查程式碼的風格、潛在錯誤和不良寫法——就像文章的文法檢查器，但針對的是程式碼。[golangci-lint](https://golangci-lint.run/) 是 Go 生態系中最受歡迎的 linter，它整合了數十個 linter（如 `govet`、`errcheck`、`staticcheck` 等），只需要一個指令就能執行所有檢查。

這個 job 裡每個 step 做的事：

| 步驟 | 做什麼 | 為什麼需要 |
|------|--------|-----------|
| `actions/checkout@v4` | 下載 repository 程式碼 | Linter 需要讀取程式碼才能分析 |
| `actions/setup-go@v5` | 安裝 Go 1.24 | golangci-lint 需要 Go 環境 |
| `golangci-lint-action@v6` | 執行 golangci-lint | 檢查程式碼品質 |

其中 [golangci/golangci-lint-action](https://github.com/golangci/golangci-lint-action) 是社群維護的 Action，會自動下載並安裝 golangci-lint、快取執行結果，並在偵測到 `.golangci.yml` 時自動套用設定。


### Test Job — 自動化測試

```yaml
test:
  name: Test
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.24'
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    - name: Show coverage
      run: go tool cover -func=coverage.out
    - name: Upload coverage
      uses: actions/upload-artifact@v4
      with:
        name: coverage-report
        path: coverage.out
```

#### `go test` 的各個 Flag

| Flag | 用途 | 說明 |
|------|------|------|
| `-v` | **Verbose** 詳細輸出 | 印出每個測試函式的名稱與結果，而不只是通過/失敗的摘要 |
| `-race` | **Race Detector** 競爭偵測 | 在測試時偵測 goroutine 之間的 data race，這是 Go 並行程式中常見的 bug |
| `-coverprofile=coverage.out` | **Coverage Profile** 覆蓋率報告 | 將測試覆蓋率資料輸出到 `coverage.out` 檔案 |
| `./...` | **所有套件** | 遞迴執行所有子目錄中的測試 |

#### 測試覆蓋率（Test Coverage）

**測試覆蓋率** 是衡量「你的測試涵蓋了多少程式碼」的指標。

```bash
# Example output of go tool cover -func
github.com/user/app/handler.go:10:  HomeHandler     100.0%
github.com/user/app/handler.go:17:  HealthHandler   100.0%
github.com/user/app/handler.go:24:  VersionHandler  80.0%
total:                              (statements)    93.3%
```

- **100%** 表示該函式的每一行都被測試到了
- **80%** 表示有 20% 的程式碼路徑沒有被測試覆蓋
- 一般來說，核心業務邏輯建議達到 **80% 以上** 的覆蓋率

#### Artifact 上傳

```yaml
- name: Upload coverage
  uses: actions/upload-artifact@v4
  with:
    name: coverage-report
    path: coverage.out
```

這是我們第一次用到 **artifact**。

不同 job 在不同的 Runner 上執行，檔案系統互不相通，所以 job 之間如果要傳檔案就要靠 artifact。`upload-artifact` 把檔案上傳到 GitHub，後面的 job 可以用 `download-artifact` 取回，或是在 Actions 頁面手動下載。

```yaml
# Upload (in one job)
- uses: actions/upload-artifact@v4
  with:
    name: coverage-report      # Artifact name
    path: coverage.out         # File or directory to upload
    retention-days: 7          # How long to keep (default: 90 days)

# Download (in another job)
- uses: actions/download-artifact@v4
  with:
    name: coverage-report      # Must match the upload name
    path: ./downloaded/        # Where to save
```

Artifact 預設保存 90 天，可以用 `retention-days` 自訂（最短 1 天，最長 90 天，Pro/Team 方案最長 400 天）。


### Build Job — 建置可執行檔

```yaml
build:
  name: Build
  needs: [lint, test]
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.24'
    - name: Build binary
      run: go build -o bin/app ./...
    - name: Upload binary
      uses: actions/upload-artifact@v4
      with:
        name: app-binary
        path: bin/app
```

#### `needs` 關鍵字 — Job 依賴關係

```yaml
needs: [lint, test]
```

`needs` 用來定義 **job 之間的執行順序**：

- **沒有 `needs`** 的 job（如 `lint` 和 `test`）會 **平行執行**，同時開始
- **有 `needs`** 的 job（如 `build`）會 **等待** 指定的 job 全部完成後才開始
- 如果 `needs` 中的任何一個 job **失敗**，這個 job **不會執行**

#### 為什麼 Build 要等 Lint 和 Test 通過？

想想看：如果程式碼有 lint 錯誤或測試失敗，我們還需要花時間去 build 嗎？

答案是 **不需要**。既然程式碼品質有問題，就應該先修好再建置。這樣可以：

1. **節省 CI 資源**：不浪費時間在必然無用的建置上
2. **明確的失敗訊號**：開發者知道問題出在 lint 或 test，而不是 build
3. **邏輯上的正確性**：只有品質通過檢查的程式碼才值得建置

Build 產出的 binary 同樣用 `upload-artifact` 上傳，後續的 deploy job 或 release workflow 就可以直接取用這份 binary，不用重新編譯。

> **關於快取**：`actions/setup-go@v5` 偵測到 `go.sum` 時會自動快取 Go modules，後續執行只要 `go.sum` 沒變就會直接用快取，不用額外設定。

## PR 觸發的 CI

前面的 CI workflow 同時設了 `on: push` 和 `on: pull_request`，兩者的用途不一樣：

- **`pull_request`**：在 PR 階段就先檢查一次，確認合併後不會壞掉。它會在**模擬合併後的 merge commit** 上執行，能及早發現合併衝突或不相容。
- **`push`**：合併到 `main` 之後的最後防線，確認分支上的程式碼真的正常。

兩個一起設可以兼顧「合併前預防」和「合併後驗證」。如果你的 workflow 只想在 PR 階段跑，也可以只設 `pull_request`：

```yaml
on:
  pull_request:
    branches: [main]
```

### `pull_request` 的觸發類型

`pull_request` 事件可以指定 activity types，控制什麼情況下要重跑 workflow：

| Type | 觸發時機 | 說明 |
|------|---------|------|
| `opened` | PR 剛建立時 | 第一次提交 PR |
| `synchronize` | 有新 commit push 到 PR 時 | 更新程式碼後重新檢查 |
| `reopened` | 關閉後重新開啟 PR 時 | 重新開啟需要再次檢查 |
| `ready_for_review` | 從 Draft 變成 Ready 時 | 只在正式 review 時才跑的檢查 |

預設會包含 `opened`、`synchronize`、`reopened`，多數情況不用自己指定。

### Permissions 設定

PR workflow 建議明確設定 `permissions`，遵循最小權限原則：

```yaml
permissions:
  contents: read
  pull-requests: write
  checks: write
```

| 權限 | 用途 |
|------|------|
| `contents: read` | 讀取 repository 的程式碼 |
| `pull-requests: write` | 允許在 PR 上留言或更新 status |
| `checks: write` | 允許建立 check run 結果 |

### Fork PR 的安全限制

當有人從 **fork** 的 repository 提交 PR 時，GitHub 會有一些安全限制，避免惡意的 fork PR 竊取你的 secrets 或修改 repository：

- **無法存取 secrets**：fork PR 的 workflow 無法讀取 repository 的 secrets
- **寫入權限受限**：預設只有 read 權限
- **程式碼檢查仍可執行**：lint、test、build 等不需要 secrets 的檢查不受影響


## 常見問題排解

### 1. `go test` 在 CI 通過但本地失敗（或反過來）

最常見的原因是 **Race Detector** 的行為差異。CI 中使用 `-race` flag，但你本地可能沒加。另一個原因是環境差異，例如 Go 版本不同、作業系統不同。

**排解方式**：在本地也用 `go test -race ./...` 來跑測試，確保行為一致。

### 2. golangci-lint 報錯但本地沒事

可能是 golangci-lint 版本不同。CI 中用 `version: latest` 會拿到最新版，但你本地可能是舊版。

**排解方式**：在 `.golangci.yml` 中鎖定你要的規則，或在 CI 中指定 golangci-lint 的版本號而非 `latest`。

### 3. Cache 沒有生效

`setup-go` 用 `go.sum` 的 hash 作為 cache key，只要 `go.sum` 有改動 cache 就會失效。第一次跑也一定沒有 cache。在 Actions log 中搜尋 "cache" 關鍵字可以看到 cache hit 或 miss 的訊息。


## 小結與練習題

這一章你為一個真正的 Go 專案建起了完整的 CI pipeline：用 `golangci-lint` 檢查程式碼、用 `go test -race` 跑單元測試並收集 coverage、再編譯出 binary，並透過 `needs` 讓 lint 和 test 平行執行、build 等前兩者都通過才啟動。你也認識了 `upload-artifact` 如何在 job 之間傳遞產物，並理解 `on: push` 搭配 `on: pull_request` 如何兼顧「合併前預防」與「合併後驗證」。下一章我們會把這支 binary 實際部署到雲端，讓整個 CI/CD 流程完整串起來。

### 本章重點回顧

- 一個完整的 Go CI pipeline 通常包含 **Lint → Test → Build** 三個階段
- 使用 `needs` 定義 job 之間的依賴關係，`lint` 和 `test` 平行執行，`build` 等它們都通過才執行
- **golangci-lint** 是 Go 最受歡迎的 linter 工具，可透過 `golangci-lint-action` 在 CI 中使用
- `go test -v -race -coverprofile=coverage.out ./...` 是標準的 CI 測試指令
- **Artifacts** 用來在不同 job 之間傳遞檔案（例如把 binary 從 build job 傳給 deploy job）
- 同時設 `on: push` 和 `on: pull_request` 可以兼顧「合併前預防」和「合併後驗證」

### 關於 Release 自動化

當專案穩定後，你會希望把某個特定版本「標記」起來發佈給使用者。這就是 **Release** 的概念：用 **Semantic Versioning**（例如 `v1.2.3`，MAJOR 代表破壞性變更、MINOR 代表新功能、PATCH 代表 bug 修正）搭配 **Git Tag** 來標記版本，並透過 `on: push: tags: ['v*']` 觸發一個 release workflow，自動建置產物並建立 GitHub Release。本工作坊不展開細節，有興趣可以查閱 [softprops/action-gh-release](https://github.com/softprops/action-gh-release)。

### 練習題

完成以下練習來鞏固本章所學：

[練習二：CI Pipeline 實戰練習](exercises/02-ci-pipeline.md)

[← 上一章：GitHub Actions 基礎](02-github-actions-basics.md) ｜ [下一章：部署到雲端平台 →](04-deployment.md)
