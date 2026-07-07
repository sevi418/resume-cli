# resume-cli

AI 简历解析命令行工具。读取 PDF 简历，调用大模型提取结构化信息，并根据岗位描述（JD）进行匹配评分。

核心流程：

```text
PDF 简历 -> 提取文本 -> AI 结构化解析 -> JD 匹配评分 -> JSON 输出
```

---

## 快速演示

无需 API Key，直接用 mock 模式跑通三个核心命令：

```bash
make build

./resume-cli parse ./testdata/01_zh_fullstack_senior.pdf
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --mock
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --mock
```

查看调试日志：

```bash
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --mock --verbose
```

保存结果到文件：

```bash
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --mock --output result.json
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --mock --output score.json
```

说明：源码构建后生成的是当前目录下的本地二进制，直接运行时使用 `./resume-cli`。如果把二进制安装到 PATH 中，也可以使用 `resume-cli`。

---

## 安装

### 前置要求

- Go 1.21+ — [下载安装](https://go.dev/dl/)
- （可选）[poppler](https://poppler.freedesktop.org/) — 启用 `pdftotext` 降级方案，提升 PDF 兼容性

**安装 poppler（可选）**

```bash
# macOS
brew install poppler

# Debian/Ubuntu
sudo apt-get install poppler-utils

# Windows (winget)
winget install --id=oschwartz10612.Poppler

# Windows (Chocolatey)
choco install poppler

# Windows (Scoop)
scoop install poppler
```

安装后确保 `pdftotext` 在 PATH 中（Windows 下通常为 `pdftotext.exe`）。

### 从源码构建

```bash
git clone https://github.com/sevi418/resume-cli.git
cd resume-cli
```

**macOS / Linux**

```bash
make build
# 或
go build -o resume-cli .
```

构建完成后使用本地二进制：

```bash
./resume-cli --help
```

**Windows (PowerShell / CMD)**

```powershell
go build -o resume-cli.exe .
# 若已安装 Make（如 Git Bash、choco install make）
make build
```

### Docker

```bash
docker build -t resume-cli .

# mock 模式，无需 API Key
docker run --rm \
  -v "$PWD/testdata:/data" \
  resume-cli extract /data/01_zh_fullstack_senior.pdf --mock

# 真实 AI 调用
docker run --rm \
  -e OPENAI_API_KEY="$OPENAI_API_KEY" \
  -e OPENAI_API_BASE="$OPENAI_API_BASE" \
  -e OPENAI_MODEL="$OPENAI_MODEL" \
  -v "$PWD/testdata:/data" \
  resume-cli extract /data/01_zh_fullstack_senior.pdf
```

---

## CLI 命令

| 命令 | 说明 |
|------|------|
| `parse` | 从 PDF 简历提取纯文本 |
| `extract` | 调用 AI 提取结构化简历信息（姓名、联系方式、教育、技能等） |
| `score` | 对比简历与 JD，输出匹配评分与面试建议 |

查看帮助：

```bash
./resume-cli --help
./resume-cli parse --help
./resume-cli extract --help
./resume-cli score --help
```

### 全局参数

| 参数 | 适用命令 | 说明 |
|------|----------|------|
| `--output <file>` | 全部 | 将结果写入文件 |
| `--mock` | `extract`、`score` | 使用预设数据，无需 API Key |
| `--verbose` / `-v` | 全部 | 输出调试日志，日志写入 stderr |
| `--help` / `-h` | 全部 | 显示帮助 |

### `parse` — 提取 PDF 文本

```bash
./resume-cli parse ./testdata/01_zh_fullstack_senior.pdf

# 保存提取文本
./resume-cli parse ./testdata/01_zh_fullstack_senior.pdf --output resume.txt
```

**输出示例：**

```
张三
电话：13800138000
邮箱：zhangsan@example.com
城市：北京
...
```

### `extract` — AI 结构化提取

```bash
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf

# mock 模式（无需 API Key）
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --mock

# 英文简历 mock 演示
./resume-cli extract ./testdata/04_en_senior_backend.pdf --mock

# 保存到文件
./resume-cli extract ./testdata/01_zh_fullstack_senior.pdf --output result.json
```

**输出示例：**

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
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt

# mock 模式
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --mock

# 保存评分结果
./resume-cli score ./testdata/01_zh_fullstack_senior.pdf --jd ./testdata/sample_jd.txt --output score.json
```

**输出示例：**

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

---

## 环境变量

| 变量 | 必需 | 说明 |
|------|------|------|
| `OPENAI_API_KEY` | 是（非 mock 模式） | OpenAI API 密钥 |
| `OPENAI_API_BASE` | 否 | 自定义 API 端点，默认 `https://api.openai.com/v1` |
| `OPENAI_MODEL` | 否 | 模型名称，默认 `gpt-4o-mini` |

推荐使用 `.env`（项目启动时会自动读取，且不覆盖你已经 export 的系统变量）：

```bash
make setup-env
# 编辑 .env 填入真实密钥
```

`.env` 示例：

```bash
OPENAI_API_KEY=sk-...
OPENAI_API_BASE=https://api.openai.com/v1
OPENAI_MODEL=gpt-4o-mini
```

也可以直接用 shell 环境变量：

```bash
export OPENAI_API_KEY="sk-..."
export OPENAI_API_BASE="https://api.openai.com/v1"   # 可选
export OPENAI_MODEL="gpt-4o-mini"                   # 可选
```

无 API Key 时可用 `--mock` 模式演示全部流程。

---

## 语言支持

- 支持中文、英文及中英混合的文本型 PDF 简历
- JSON 字段名固定为英文，字段值保留简历原文语言
- `score` 的 `comment` 和 `interview_questions` 会跟随 JD 的主要语言
- 不提供 CLI 界面 i18n（没有 `--lang` 参数）

---

## 技术选型

| 模块 | 选型 | 理由 |
|------|------|------|
| CLI 框架 | [spf13/cobra](https://github.com/spf13/cobra) | Go CLI 事实标准，原生子命令、自动 `--help`、PersistentFlags |
| PDF 解析 | [ledongthuc/pdf](https://github.com/ledongthuc/pdf)（主）+ `pdftotext`（降级） | 纯 Go、零 CGO；文本为空时尝试系统 `pdftotext` 提升兼容性 |
| AI 调用 | [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) | 社区成熟，支持自定义 BaseURL，兼容 OpenAI 及各类代理 |
| 日志 | `log/slog`（Go 1.21+） | 标准库结构化日志，零额外依赖 |
| 测试 | `testing` + [testify/assert](https://github.com/stretchr/testify) | 可读断言，社区主流 |
| 构建 | Makefile + Dockerfile | 本地开发便捷，多阶段构建镜像极小 |

---

## 项目结构

```
resume-cli/
├── cmd/                    # Cobra 命令定义
│   ├── root.go            # 根命令，全局 flags
│   ├── parse.go           # parse 子命令
│   ├── extract.go         # extract 子命令
│   └── score.go           # score 子命令
├── internal/
│   ├── pdf/               # PDF 文本提取
│   ├── ai/                # AI 客户端、extract、score、mock
│   ├── model/             # Resume、Score 等数据模型
│   └── util/              # JSON 修复与校验
├── testdata/              # 示例 PDF 与 JD
│   ├── 01_zh_fullstack_senior.pdf  # 中文高级全栈（主演示）
│   ├── 02_zh_backend_junior.pdf    # 中文初级后端（低年限样本）
│   ├── 03_zh_product_manager.pdf   # 中文产品经理（低匹配交叉测试）
│   ├── 04_en_senior_backend.pdf    # 英文高级后端（英文主演示）
│   ├── 05_en_fullstack_mid.pdf     # 英文中级全栈（英文中等匹配）
│   ├── 06_en_marketing_manager.pdf # 英文营销经理（低匹配交叉测试）
│   └── sample_jd.txt               # 中文高级全栈 JD（score 统一输入）
├── main.go
├── go.mod
├── .env.example
├── Makefile
├── Dockerfile
├── CHECKLIST.md
└── README.md
```

---

## 开发

```bash
make build    # 编译
make test     # 运行测试
make lint     # 代码检查
make setup-env # 创建 .env（若不存在）
make clean    # 清理构建产物
```

---

## 已实现功能

### 核心功能

- [x] `parse` — PDF 文本提取，含文件不存在 / 非 PDF / 无法读取 / 文本为空等错误提示
- [x] `extract` — AI 结构化信息提取，JSON 输出与基本字段校验
- [x] `score` — JD 匹配评分（0–100），含评语与面试建议
- [x] 清晰的命令参数与 `--help`
- [x] 基础单元测试与 mock 模式演示

### 扩展能力

- [x] `--output result.json` 保存结果
- [x] `--mock` 模式，无 API Key 可演示
- [x] AI 返回 JSON 自动修复（去除 markdown 代码块、尾逗号等）
- [x] `--verbose` 结构化日志
- [x] Makefile
- [x] Dockerfile

---

## 功能与实现

| 功能 | 实现方式 |
|------|----------|
| PDF 简历文本解析 | `parse <pdf_path>`，位于 `internal/pdf` |
| AI 结构化信息提取 | `extract <pdf_path>`，输出 `name`、`phone`、`email`、`city`、`education`、`skills` |
| JD 匹配评分 | `score <pdf_path> --jd <jd_path>`，输出 0-100 评分、评语和面试问题 |
| CLI 帮助 | Cobra 自动生成 `--help` |
| 示例命令和输出 | README 的 CLI 命令章节 |
| 基础测试 | `go test ./...` / `make test` |
| Mock 演示 | `--mock` |
| 保存结果 | `--output <file>` |
| 日志输出 | `--verbose` |
| 构建与容器 | `Makefile` / `Dockerfile` |

---

## 错误处理

| 场景 | 处理方式 |
|------|----------|
| 文件不存在 | 明确提示路径 |
| 非 PDF 文件 | 检查扩展名，提示格式错误 |
| PDF 无法解析 | 捕获错误，提示文件可能损坏 |
| PDF 文本为空 | 提示可能是扫描件，建议 OCR |
| API Key 未设置 | 提示设置环境变量或使用 `--mock` |
| AI 调用失败 | 显示网络 / 限流 / 格式等具体错误 |
| AI 返回非 JSON | 尝试自动修复，失败则报错 |
| 缺少 `--jd` | 提示 `--jd is required` |
| JD 文件为空 | 提示 JD 内容为空 |

---

## 已知问题与限制

1. **PDF 提取质量**：`ledongthuc/pdf` 对复杂排版效果有限，简单文本型简历 PDF 足够；扫描件无法提取文本，不在本次 Demo 范围
2. **AI 输出稳定性**：通过 JSON Mode + 自动修复缓解，极端情况下仍可能需重试
3. **降级依赖**：`pdftotext` 需系统安装 poppler-utils，未安装时仅使用纯 Go 方案

---

## 测试

```bash
make test
```

| 类型 | 覆盖内容 |
|------|----------|
| 单元测试 | PDF 解析（中英文正常 / 异常）、JSON 修复、模型校验 |
| 集成测试 | Mock 模式下 `parse`、`extract` 和 `score` 端到端 |
| 手动测试 | 真实 API 调用，验证输出格式 |
