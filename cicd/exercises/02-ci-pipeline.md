# 練習二：CI Pipeline 實戰練習

> **難度：** 中級 | **對應章節：** 03-Go 專案 CI Pipeline

---

## 目錄

- [練習 2-1：Multi-Job 暖身](#練習-2-1multi-job-暖身)
- [練習 2-2：擴充 CI Pipeline](#練習-2-2擴充-ci-pipeline)
- [練習 2-3：Matrix 策略練習](#練習-2-3matrix-策略練習)
- [練習 2-4：PR 檢查練習](#練習-2-4pr-檢查練習)
- [練習 2-5：Coverage 門檻挑戰](#練習-2-5coverage-門檻挑戰)
- [延伸思考](#延伸思考)

---

## 練習 2-1：Multi-Job 暖身

### 目標

在動手擴充 `ci.yml` 之前，先單獨練一次 `needs`——建立一個多 job 的 workflow，把「平行 + 合流」的依賴關係跑一遍，建立對 job DAG 的直覺。

### 要求

1. 建立一個新的 workflow 檔案 `.github/workflows/multi-job.yml`
2. **Job `info`**：印出 repository 名稱、branch、commit SHA（使用 `${{ github.xxx }}` context）
3. **Job `env-check`**：印出系統資訊（`uname -s`、`uname -m`、`date`）
4. **Job `summary`**：`needs: [info, env-check]`，在前兩個 job 都完成後執行，印出 `"All checks passed"`
5. 確認在 Actions 頁面上，`info` 和 `env-check` 是 **平行執行**，`summary` 在兩者都完成後才開始

### 提示

- 沒有 `needs` 的 job 會平行執行
- `needs: [job-a, job-b]` 可以一次等多個 job 完成
- 如果任何一個被依賴的 job 失敗，後續 job 預設不會執行

### 預期結果

在 GitHub Actions 頁面看到的依賴圖會長這樣：

```
info      ──┐
            ├──▶ summary
env-check ──┘
```

這和 03 章介紹的 `lint` + `test` → `build` 是同一種模式。

<details>
<summary>點擊查看答案</summary>

```yaml
name: Multi Job Workflow

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  info:
    name: Repo Info
    runs-on: ubuntu-latest
    steps:
      - name: Show GitHub context
        run: |
          echo "Repository: ${{ github.repository }}"
          echo "Branch: ${{ github.ref_name }}"
          echo "Commit SHA: ${{ github.sha }}"
          echo "Triggered by: ${{ github.actor }}"

  env-check:
    name: Environment Check
    runs-on: ubuntu-latest
    steps:
      - name: Show system info
        run: |
          echo "OS: $(uname -s)"
          echo "Architecture: $(uname -m)"
          echo "Current Time: $(date)"

  summary:
    name: Summary
    needs: [info, env-check]
    runs-on: ubuntu-latest
    steps:
      - run: echo "All checks passed"
```

**重點說明：**

- `info` 和 `env-check` 都沒設 `needs`，所以會 **平行** 開始跑
- `summary` 用 `needs: [info, env-check]` 等兩個都成功才執行
- 如果其中一個失敗，`summary` 預設不會執行（可以用 `if: always()` 覆蓋這個行為）
- 這正是 03 章 `ci.yml` 裡 `lint` + `test` → `build` 的 DAG 模式，接下來的練習 2-2 就會直接擴充這支 pipeline

</details>

---

## 練習 2-2：擴充 CI Pipeline

### 目標

在課程教材中的 `ci.yml` 基礎上，加入更多程式碼品質檢查，建立更完整的 CI pipeline。

### 要求

在 `ci.yml` 的 `test` job 中，在 `Run tests` 步驟 **之前** 加入以下三個檢查：

1. **`go vet ./...`** — Go 內建的靜態分析工具，可以檢查常見的程式碼錯誤（例如 `fmt.Printf` 的格式化字串與參數不匹配）
2. **`go mod verify`** — 驗證 `go.sum` 中記錄的 module hash 是否正確，確保 dependencies 沒有被篡改
3. **程式碼格式檢查** — 使用 `gofmt -l .` 列出所有未格式化的檔案。如果有任何檔案未格式化，workflow 就應該 fail

### 提示

- `gofmt -l .` 會列出所有格式不正確的檔案。如果輸出為空，代表所有檔案都已正確格式化。
- 可以用 shell 的條件判斷來檢查 `gofmt` 的輸出是否為空。
- 這些檢查步驟都放在 `test` job 中，並且在 `Run tests` 步驟之前。

### 預期結果

- 如果有 `go vet` 發現的問題，CI 會在 "Run go vet" 步驟失敗。
- 如果 `go.sum` 不一致，CI 會在 "Verify dependencies" 步驟失敗。
- 如果有未格式化的檔案，CI 會在 "Check code formatting" 步驟失敗，並列出哪些檔案需要格式化。

<details>
<summary>點擊查看答案</summary>

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

      # === New checks start here ===
      - name: Run go vet
        run: go vet ./...

      - name: Verify dependencies
        run: go mod verify

      - name: Check code formatting
        run: |
          UNFORMATTED=$(gofmt -l .)
          if [ -n "$UNFORMATTED" ]; then
            echo "The following files are not properly formatted:"
            echo "$UNFORMATTED"
            echo ""
            echo "Please run 'gofmt -w .' to fix formatting."
            exit 1
          fi
          echo "All files are properly formatted."
      # === New checks end here ===

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
        run: go build -o bin/app .
      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: app-binary
          path: bin/app
```

**重點說明：**

- `go vet` 是 Go 內建的靜態分析工具，與 `golangci-lint` 不同，它不需要額外安裝。
- `go mod verify` 會比對本地 module cache 與 `go.sum` 的 hash 值，確保沒有人篡改依賴套件。
- `gofmt -l .` 只列出檔案，不修改。如果輸出不為空，代表有檔案需要格式化。使用 `exit 1` 讓 step 失敗。
- 這三個檢查放在 `Run tests` 前面，如果基本品質不過關，就不用浪費時間跑測試。

</details>

---

## 練習 2-3：Matrix 策略練習

### 目標

使用 **matrix strategy** 同時在多個 Go 版本和作業系統上執行測試，確保程式碼的跨環境相容性。

> **提示**：matrix 的完整語法請參考 [GitHub Actions 參考手冊](../02-reference.md)，主教材沒有展開這部分。

### 要求

1. 修改（或建立新的）workflow，在 test job 中使用 matrix strategy
2. 測試 **Go 1.24** 和 **Go 1.25** 兩個版本（sample-app 的 `go.mod` 需要 Go 1.24+）
3. 測試 **ubuntu-latest** 和 **macos-latest** 兩個作業系統
4. 使用 `include` 額外加入一個 **Go 1.25 + windows-latest** 的組合
5. 在測試步驟中印出目前的 **Go 版本** 和 **作業系統**

### 提示

- `matrix` 定義在 `strategy` 底下。
- `${{ matrix.go-version }}` 和 `${{ matrix.os }}` 可以在 step 中取得目前的 matrix 值。
- `runs-on` 要設為 `${{ matrix.os }}`，讓每個組合在不同的 Runner 上跑。
- `include` 可以在原本的組合之外，額外加入特定的組合。
- 預設的 2 x 2 = 4 個組合，加上 `include` 的 1 個，總共會有 **5 個** job。

### 預期結果

在 GitHub Actions 頁面上，你應該會看到 5 個 test job 同時平行執行：

| Job | Go Version | OS |
|-----|-----------|-----|
| Test (Go 1.24, ubuntu-latest) | 1.24 | ubuntu-latest |
| Test (Go 1.24, macos-latest) | 1.24 | macos-latest |
| Test (Go 1.25, ubuntu-latest) | 1.25 | ubuntu-latest |
| Test (Go 1.25, macos-latest) | 1.25 | macos-latest |
| Test (Go 1.25, windows-latest) | 1.25 | windows-latest |

<details>
<summary>點擊查看答案</summary>

```yaml
name: Matrix Test

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  test:
    name: Test (Go ${{ matrix.go-version }}, ${{ matrix.os }})
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        go-version: ['1.24', '1.25']
        os: [ubuntu-latest, macos-latest]
        include:
          - go-version: '1.25'
            os: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Show Environment
        run: |
          echo "Go Version: $(go version)"
          echo "OS: ${{ matrix.os }}"
          echo "Runner OS: ${{ runner.os }}"
          echo "Runner Arch: ${{ runner.arch }}"

      - name: Run tests
        run: go test -v -race ./...
```

**重點說明：**

- `fail-fast: false` 設定讓所有 matrix job 都會執行完畢，即使其中一個失敗。這樣你可以看到完整的測試結果，知道是哪些環境有問題。
- `include` 不會覆蓋原本的組合，而是**額外增加**一個 Go 1.25 + windows-latest 的組合。
- 如果某些環境的測試會失敗（例如路徑分隔符在 Windows 上不同），你可以用 `exclude` 排除特定組合。
- job 名稱中使用 `${{ matrix.go-version }}` 和 `${{ matrix.os }}` 讓每個 job 的名稱都不一樣，方便辨識。
- **版本一定要加引號**：YAML 會把 `1.20` 解析成浮點數 `1.2`，結果你以為指定了 Go 1.20，實際下載的是 1.2。永遠用 `'1.24'` 這種字串格式。
- **Runner 計價注意**：在 GitHub 免費額度裡，macOS Runner 是 Linux 的 10 倍計價、Windows 是 2 倍。真實專案如果不需要真的測那些平台，就不要隨便加進 matrix。

</details>

---

## 練習 2-4：PR 檢查練習

### 目標

實際走一遍完整的 **PR 檢查流程**：建立分支、修改程式碼、撰寫測試、開 PR、觀察 CI 結果。

### 要求

1. 從 `main` 建立一個新的分支 `feature/add-info-endpoint`
2. 在範例專案中新增一個 `/info` endpoint，回傳應用程式的資訊
3. 為新的 endpoint 撰寫對應的測試
4. 開一個 Pull Request 到 `main`，並觀察 CI 檢查的結果

> 這個練習 **不提供完整答案**，而是提供 step by step 的引導。目的是讓你走一遍真實的開發流程。

### Step by Step 引導

#### Step 1：建立新分支

```bash
git checkout main
git pull origin main
git checkout -b feature/add-info-endpoint
```

#### Step 2：新增 `/info` endpoint

在 `handler.go` 中新增一個 handler 函式（注意維持與既有 handler 相同的命名風格與 method 路由風格）：

```go
// handleInfo returns application metadata.
func handleInfo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "app":     "sample-app",
        "version": version,
    })
}
```

然後在 `main.go` 中註冊這個 handler：

```go
mux.HandleFunc("GET /info", handleInfo)
```

#### Step 3：撰寫測試

在 `handler_test.go` 中新增測試：

```go
func TestHandleInfo(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/info", nil)
    w := httptest.NewRecorder()
    handleInfo(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }

    // Verify response contains expected fields
    var result map[string]string
    if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if result["app"] != "sample-app" {
        t.Errorf("expected app to be 'sample-app', got '%s'", result["app"])
    }
}
```

#### Step 4：在本地驗證

```bash
# Run tests locally (same flags as CI)
go test -v -race ./...

# Run go vet
go vet ./...

# Check formatting
gofmt -l .
```

確認所有檢查都通過後再繼續。

#### Step 5：Commit 並 Push

```bash
git add handler.go main.go handler_test.go
git commit -m "feat: add /info endpoint"
git push origin feature/add-info-endpoint
```

#### Step 6：開 Pull Request

1. 到 GitHub repository 頁面
2. 你應該會看到一個提示，讓你建立 PR
3. 填寫 PR 標題和描述
4. 點擊 **Create pull request**

#### Step 7：觀察 CI 結果

1. 在 PR 頁面下方，你會看到 CI 檢查的狀態
2. 點擊 **Details** 查看詳細的執行記錄
3. 確認所有檢查都通過（綠色勾勾）

### 思考問題

- 如果 CI 檢查失敗了，你會怎麼做？
- 如果你在 push 之後又發現了一個 bug，你會怎麼處理？（提示：直接在同一個分支上修復並 push，CI 會自動重跑）

---

## 練習 2-5：Coverage 門檻挑戰

### 目標

在 CI pipeline 中加入 **測試覆蓋率門檻檢查**，確保測試覆蓋率不低於 70%。如果低於門檻，CI 就會 fail。

### 要求

1. 在 `ci.yml` 的 test job 中加入覆蓋率門檻檢查
2. 設定門檻為 **70%**
3. 如果覆蓋率低於 70%，workflow 要 **fail** 並顯示目前的覆蓋率
4. 將覆蓋率報告（`coverage.out`）和 HTML 版本的覆蓋率報告 **上傳為 artifact**

### 提示

- `go test -coverprofile=coverage.out ./...` 會產生覆蓋率報告。
- `go tool cover -func=coverage.out` 的最後一行會顯示 **總覆蓋率**，格式類似 `total: (statements) 85.7%`。
- 你可以用 `grep`、`awk` 等工具從輸出中提取覆蓋率數字。
- `go tool cover -html=coverage.out -o coverage.html` 可以產生 HTML 版本的覆蓋率報告。

### 預期結果

- 如果覆蓋率 >= 70%：CI 通過，log 中顯示 "Coverage is above threshold"。
- 如果覆蓋率 < 70%：CI 失敗，log 中顯示 "Coverage is below threshold" 和目前的覆蓋率。
- Artifact 頁面可以下載 HTML 覆蓋率報告。

<details>
<summary>點擊查看答案</summary>

> 為了聚焦，以下只展示 `test` job 的改動。`lint` 和 `build` job 維持 2-2 的結果即可——也就是 `test` job 前面的 `go vet` / `go mod verify` / `gofmt` 三個檢查步驟都保留。

```yaml
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      # ... (2-2 加的 go vet / go mod verify / gofmt 檢查)

      - name: Run tests with coverage
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Show coverage summary
        run: go tool cover -func=coverage.out

      # Generate HTML report BEFORE the threshold check, so it's always
      # available in the artifact even when coverage fails.
      - name: Generate HTML coverage report
        run: go tool cover -html=coverage.out -o coverage.html

      - name: Upload coverage reports
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: coverage-reports
          path: |
            coverage.out
            coverage.html
          retention-days: 14

      - name: Check coverage threshold
        run: |
          THRESHOLD=70
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Total coverage: ${COVERAGE}%"
          echo "Threshold: ${THRESHOLD}%"
          echo ""

          # Use awk for floating point comparison
          PASS=$(echo "$COVERAGE $THRESHOLD" | awk '{print ($1 >= $2)}')
          if [ "$PASS" = "1" ]; then
            echo "Coverage is above threshold. (${COVERAGE}% >= ${THRESHOLD}%)"
          else
            echo "::error::Coverage is below threshold! (${COVERAGE}% < ${THRESHOLD}%)"
            exit 1
          fi
```

**重點說明：**

- `go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//'` 這行指令會：
  1. 顯示每個函式的覆蓋率
  2. `grep total` 取出最後一行（總覆蓋率）
  3. `awk '{print $3}'` 取出第三個欄位（覆蓋率百分比，如 `85.7%`）
  4. `sed 's/%//'` 去掉百分號（得到 `85.7`）
- `awk '{print ($1 >= $2)}'` 用來做浮點數比較，因為 shell 原生不支援浮點數比較。
- `::error::` 是 GitHub Actions 的特殊語法，會在 workflow summary 中顯示一個紅色的錯誤訊息。
- **步驟順序很重要**：Generate HTML 和 Upload 要放在 Check threshold **之前**。因為 threshold 檢查失敗時會 `exit 1`，後面的 step 預設不會執行——如果 HTML 放在後面，低於門檻時 artifact 就沒有 HTML 報告，正好是你最需要它的時候。
- `if: always()` 確保即使後面的 threshold check 失敗，artifact 仍然會被上傳。
- `retention-days: 14` 設定 artifact 保留 14 天。

</details>

---

## 延伸思考

以下問題沒有標準答案，供進階學生思考和討論：

1. **覆蓋率的迷思**：100% 的測試覆蓋率是否代表程式沒有 bug？為什麼？你認為合理的覆蓋率目標應該是多少？

2. **Matrix 的成本**：如果你的 matrix 有 3 個 Go 版本 x 3 個 OS x 2 個資料庫版本 = 18 個 job，這樣好嗎？在 CI 速度和測試完整性之間，你會怎麼取捨？

3. **Flaky Tests**：有些測試有時候通過、有時候失敗（稱為 flaky test）。這會對 CI pipeline 造成什麼問題？你會怎麼處理？

4. **CI 效能最佳化**：如果你的 CI pipeline 跑了 15 分鐘，你會用什麼方法來加速？列出至少三種可能的最佳化策略。

---

[← 練習一：GitHub Actions 基礎練習](01-basics.md) ｜ [回到教材目錄 →](../README.md)
