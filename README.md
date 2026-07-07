# resume-cli

AI 简历解析命令行工具。读取 PDF 简历，调用大模型提取结构化信息，并根据岗位描述（JD）进行匹配评分。

```text
PDF 简历 -> 提取文本 -> AI 结构化解析 -> JD 匹配评分 -> JSON 输出
```

## 已实现功能

| 命令 | 说明 | 实现位置 |
|------|------|----------|
| `parse` | 从 PDF 提取纯文本 | `internal/pdf` |
| `extract` | AI 提取结构化简历信息（姓名、联系方式、教育、技能等） | `internal/ai` |
| `score` | 对比简历与 JD，输出 0–100 匹配评分、评语与面试建议 | `internal/ai` |

支持能力：`--mock`（无 API Key 演示）、`--output`（保存结果）、`--verbose`（调试日志）、AI 返回 JSON 自动修复、Makefile、Dockerfile。

## 快速开始

无需 API Key，用 mock 模式跑通三个命令：

```bash
make build

./resume-cli parse   ./testdata/01_zh_fullstack_senior.pdf
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --mock
./resume-cli score   ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --mock
```

源码构建后是当前目录下的本地二进制，直接运行用 `./resume-cli`；若装入 PATH 则可直接用 `resume-cli`。

## 安装

前置依赖：

| 依赖 | Release 二进制 | 源码构建 | 说明 |
|------|:------------:|:--------:|------|
| Go 1.24.1+ | 不需要 | 需要 | 仅编译时需要 |
| [poppler](https://poppler.freedesktop.org/) | 可选 | 可选 | 提供 `pdftotext` 降级，提升复杂 PDF 兼容性 |
| `OPENAI_API_KEY` | 非 mock 时需要 | 非 mock 时需要 | 可用 `--mock` 跳过 |
| PDF / JD 文件 | 需要 | 需要 | 自备简历 PDF；`score` 另需 JD 文本文件 |

### 方式一：下载 Release 二进制

从 [GitHub Releases](https://github.com/sevi418/resume-cli/releases) 下载对应平台文件，**无需安装 Go**：

| 平台 | 文件 |
|------|------|
| macOS (Apple Silicon) | `resume-cli-darwin-arm64` |
| macOS (Intel) | `resume-cli-darwin-amd64` |
| Linux (amd64) | `resume-cli-linux-amd64` |
| Linux (arm64) | `resume-cli-linux-arm64` |
| Windows | `resume-cli-windows-amd64.exe` |

```bash
# macOS / Linux，按平台替换文件名
chmod +x resume-cli-darwin-arm64
./resume-cli-darwin-arm64 extract ./my-resume.pdf --mock
```

```powershell
# Windows PowerShell
.\resume-cli-windows-amd64.exe extract .\my-resume.pdf --mock
```

需要示例数据时可单独 clone 仓库获取 `testdata/`。

### 方式二：从源码构建

需要 Go 1.24.1+：

```bash
git clone https://github.com/sevi418/resume-cli.git
cd resume-cli
make build              # 或 go build -o resume-cli .
./resume-cli --help
```

Windows：`go build -o resume-cli.exe .`

### 方式三：Docker

```bash
docker build -t resume-cli .

# mock 模式
docker run --rm -v "$PWD/testdata:/data" \
  resume-cli extract /data/01_zh_fullstack_senior.pdf --mock

# 真实 AI 调用
docker run --rm \
  -e OPENAI_API_KEY -e OPENAI_API_BASE -e OPENAI_MODEL \
  -v "$PWD/testdata:/data" \
  resume-cli extract /data/01_zh_fullstack_senior.pdf
```

### poppler（可选）

```bash
brew install poppler                      # macOS
sudo apt-get install poppler-utils        # Debian/Ubuntu
winget install --id=oschwartz10612.Poppler  # Windows
```

未安装时仅使用纯 Go PDF 解析，部分复杂 PDF 可能提取为空。

## 命令用法

全局参数：

| 参数 | 适用命令 | 说明 |
|------|----------|------|
| `--output <file>` / `-o` | 全部 | 将结果写入文件（默认 stdout） |
| `--mock` | `extract`、`score` | 使用预设数据，无需 API Key |
| `--verbose` / `-v` | 全部 | 输出调试日志（写入 stderr） |
| `--help` / `-h` | 全部 | 显示帮助 |

### `parse` — 提取 PDF 文本

```bash
./resume-cli parse ./testdata/01_zh_fullstack_senior.pdf
./resume-cli parse ./testdata/01_zh_fullstack_senior.pdf --output resume.txt
```

```text
张三
电话：13800138000
邮箱：zhangsan@example.com
城市：北京
...
```

### `extract` — AI 结构化提取

```bash
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --mock
./resume-cli extract ./testdata/04_en_senior_backend.pdf --mock   # 英文简历
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --output result.json
```

```json
{
  "name": "张三",
  "phone": "13800138000",
  "email": "zhangsan@example.com",
  "city": "北京",
  "education": [
    {
      "school": "清华大学",
      "major": "计算机科学与技术",
      "degree": "本科",
      "graduation_time": "2022-06"
    }
  ],
  "skills": ["Go", "React", "PostgreSQL"]
}
```

### `score` — JD 匹配评分

```bash
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --mock
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --output score.json
```

```json
{
  "overall_score": 82,
  "skill_score": 88,
  "experience_score": 80,
  "education_score": 75,
  "comment": "候选人具备较好的全栈开发基础，技能与岗位要求较匹配，但缺少明确的大模型应用经验。",
  "interview_questions": [
    "请介绍一个你主导过的全栈项目。",
    "你是否有调用大模型 API 的实际经验？"
  ]
}
```

## 环境变量

| 变量 | 必需 | 默认值 | 说明 |
|------|------|--------|------|
| `OPENAI_API_KEY` | 非 mock 模式 | — | OpenAI API 密钥 |
| `OPENAI_API_BASE` | 否 | `https://api.openai.com/v1` | 自定义 API 端点 |
| `OPENAI_MODEL` | 否 | `gpt-4o-mini` | 模型名称 |

推荐用 `.env`（启动时自动读取，不覆盖已 export 的系统变量）：

```bash
make setup-env      # 从 .env.example 创建 .env，再填入密钥
```

也可直接 export：

```bash
export OPENAI_API_KEY="sk-..."
export OPENAI_API_BASE="https://api.openai.com/v1"   # 可选
export OPENAI_MODEL="gpt-4o-mini"                    # 可选
```

语言说明：支持中文 / 英文 / 中英混合的文本型 PDF；JSON 字段名固定英文，字段值保留原文语言；`score` 的 `comment` 与 `interview_questions` 跟随 JD 主要语言。

## 错误处理

| 场景 | 处理方式 |
|------|----------|
| 文件不存在 / 非 PDF | 校验路径与扩展名，明确提示 |
| PDF 无法解析 | 提示文件可能损坏 |
| PDF 文本为空 | 提示可能是扫描件，建议 OCR |
| API Key 未设置 | 提示设置环境变量或使用 `--mock` |
| AI 调用失败 | 显示网络 / 限流 / 格式等具体错误 |
| AI 返回非 JSON | 尝试自动修复（去 markdown 围栏、尾逗号），失败则报错 |
| 缺少 `--jd` / JD 为空 | 提示 `--jd is required` 或 JD 内容为空 |

## 技术选型

| 模块 | 选型 | 理由 |
|------|------|------|
| CLI 框架 | [spf13/cobra](https://github.com/spf13/cobra) | Go CLI 事实标准，原生子命令、自动 `--help` |
| PDF 解析 | [ledongthuc/pdf](https://github.com/ledongthuc/pdf) + `pdftotext` 降级 | 纯 Go、零 CGO；文本为空时降级提升兼容性 |
| AI 调用 | [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) | 支持自定义 BaseURL，兼容 OpenAI 及各类代理 |
| 日志 | `log/slog` | 标准库结构化日志，零额外依赖 |
| 测试 | `testing` + [testify](https://github.com/stretchr/testify) | 可读断言 |
| 构建 | Makefile + Dockerfile | 本地开发便捷，多阶段构建镜像小 |

## 项目结构

```text
resume-cli/
├── cmd/                # Cobra 命令定义（root / parse / extract / score）
├── internal/
│   ├── pdf/            # PDF 文本提取（Go 解析 + pdftotext 降级）
│   ├── ai/             # AI 客户端、extract、score、mock
│   ├── model/          # Resume、Score 数据模型与校验
│   └── util/           # JSON 修复、输出、.env 加载
├── testdata/           # 示例 PDF（中英文多档位）与 sample_jd.txt
├── cli_test.go         # CLI 端到端测试
├── main.go
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── .env.example
└── README.md
```

## 开发

```bash
make build       # 编译
make test        # 运行测试（go test ./...）
make lint        # gofmt + go vet
make setup-env   # 创建 .env
make clean       # 清理构建产物
```

测试覆盖：PDF 解析（中英文正常 / 异常）、JSON 修复、模型校验，以及 mock 模式下 `parse` / `extract` / `score` 端到端。

发布多平台二进制：

```bash
mkdir -p dist
GOOS=darwin  GOARCH=arm64 go build -o dist/resume-cli-darwin-arm64 .
GOOS=darwin  GOARCH=amd64 go build -o dist/resume-cli-darwin-amd64 .
GOOS=linux   GOARCH=amd64 go build -o dist/resume-cli-linux-amd64 .
GOOS=linux   GOARCH=arm64 go build -o dist/resume-cli-linux-arm64 .
GOOS=windows GOARCH=amd64 go build -o dist/resume-cli-windows-amd64.exe .

gh release create v0.1.0 --title "v0.1.0" --notes "AI 简历解析 CLI Demo" dist/*
```

## 已知问题与限制

1. **PDF 提取**：`ledongthuc/pdf` 对复杂排版效果有限；扫描件无法提取文本，不在 Demo 范围。
2. **AI 稳定性**：通过 JSON Mode + 自动修复缓解，极端情况仍可能需重试。
3. **降级依赖**：`pdftotext` 需系统安装 poppler，未安装时仅用纯 Go 方案。
