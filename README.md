# 🤖 Model Arena — AI 多模型实时对比擂台

> **同一个 prompt，同时发给多个大模型，看谁回答好、谁速度快、谁更便宜**

[![Go](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Release](https://img.shields.io/github/v/release/18296023612/model-arena)](https://github.com/18296023612/model-arena/releases)

---

## 📺 这是做什么的？

```
你是一个 AI 应用开发者，你需要选一个模型：
  DeepSeek? 通义千问? 智谱GLM? 百度文心?

每个模型都说自己好，你想：
  ❓ 同一个问题谁答得好？
  ❓ 谁速度快？
  ❓ 谁更便宜？

以前：一条一条 curl 手动测试 → 累
现在：model-arena run --prompt "你的问题"
```

**Model Arena 一条命令，同时测试多个模型，自动出对比报告。**

---

## ✨ 功能

| 功能 | 说明 |
|------|------|
| ✅ **多模型并发测试** | 同时调用 DeepSeek / 千问 / 智谱 / 百度 / 火山引擎 |
| ✅ **详细对比报告** | 表格展示：状态、延迟、Token 消耗 |
| ✅ **完整输出对比** | 每个模型的回答一字不差展示出来 |
| ✅ **流式输出** | `--stream` 实时看到每个模型的输出过程 |
| ✅ **配置文件** | 支持 YAML 配置，自定义任意模型 |
| ✅ **JSON 输出** | `--json` 集成 CI/CD 或监控 |
| ✅ **彩色终端** | 绿色✅成功 / 红色❌失败 / 黄色⚠️慢 |
| ✅ **零依赖** | 单文件二进制，下载即用 |

---

## 🚀 30 秒上手

### 1️⃣ 下载

从 [Releases](https://github.com/18296023612/model-arena/releases) 下载最新版。

### 2️⃣ 设置 API Key

```bash
# 至少设置一个 key 就能用
export DEEPSEEK_API_KEY=sk-your-deepseek-key
export QWEN_API_KEY=sk-your-qwen-key
export ZHIPU_API_KEY=sk-your-zhipu-key
```

### 3️⃣ 运行对比

```bash
model-arena run --prompt "介绍一下你自己"
```

输出示例：
```
╔══════════════════════════════════════════════════════════════╗
║             🤖  Model Arena — 模型对比报告               ║
╚══════════════════════════════════════════════════════════════╝

  Model                Status     Latency      Tokens   Provider
  ────────────────────────────────────────────────────────────────────────
  deepseek-chat        ✅         1,234ms      156      deepseek
  qwen-plus            ✅         2,345ms      203      alibaba
  glm-4-flash          ✅         1,567ms      178      zhipu
  ernie-speed          ❌         HTTP 401      -       baidu

  ✓ deepseek-chat  1234ms  156tokens

    你好！我是 DeepSeek，一个由深度求索公司开发的 AI 助手……
```

### 4️⃣ 进阶用法

```bash
# 只测试特定模型
model-arena run --prompt "写一首关于秋天的诗" --models deepseek-chat,qwen-plus

# 流式输出（实时看每个模型一个字一个字地输出）
model-arena run --prompt "讲个笑话" --stream

# 使用配置文件
model-arena run --prompt "Hello" --config arena.yaml

# JSON 输出
model-arena run --prompt "1+1=?" --json

# 查看配置模板
model-arena config
```

---

## 📦 配置文件

如果想自定义供应商或添加更多模型，创建 `arena.yaml`：

```yaml
models:
  - name: deepseek-chat
    alias: ds
    provider: deepseek
    base_url: https://api.deepseek.com
    api_key: ${DEEPSEEK_API_KEY}

  - name: qwen-plus
    alias: qwen
    provider: alibaba
    base_url: https://dashscope.aliyuncs.com/compatible-mode/v1
    api_key: ${QWEN_API_KEY}

  - name: glm-4-flash
    alias: glm
    provider: zhipu
    base_url: https://open.bigmodel.cn/api/paas/v4
    api_key: ${ZHIPU_API_KEY}

  - name: ernie-speed
    alias: baidu
    provider: baidu
    base_url: https://qianfan.baidubce.com/v2
    api_key: ${BAIDU_API_KEY}

  - name: doubao-pro-32k
    alias: volc
    provider: volcengine
    base_url: https://ark.cn-beijing.volces.com/api/v3
    api_key: ${VOLC_API_KEY}
```

用法：`model-arena run --prompt "你好" --config arena.yaml`

---

## 🔧 环境变量

| 变量 | 用途 |
|------|------|
| `DEEPSEEK_API_KEY` | DeepSeek API 密钥 |
| `QWEN_API_KEY` | 通义千问 API 密钥 |
| `ZHIPU_API_KEY` | 智谱 GLM API 密钥 |
| `BAIDU_API_KEY` | 百度文心 API 密钥 |
| `VOLC_API_KEY` | 火山引擎 Doubao API 密钥 |
| `NO_COLOR` | 设为任意值禁用彩色输出 |

---

## 💡 使用场景

### 场景 1：选模型

> 新项目要选模型，不知道用哪个好？
> 跑一次 model-arena，对比后再决定。

### 场景 2：验证供应商

> 供应商说升级了，真的变快了？
> 跑一次 model-arena --json，把数据存下来对比。

### 场景 3：CI/CD 检测

```bash
# 在 CI 中跑，JSON 输出，检测延迟是否异常
model-arena run --prompt "ping" --json
```

---

## 🏗️ 自行编译

```bash
git clone https://github.com/18296023612/model-arena.git
cd model-arena
go build -o model-arena .
```

需要 Go 1.21+。

---

## 📜 许可证

MIT License

---

## ⭐ 支持这个项目

如果这个工具帮到了你，欢迎 Star ⭐ 支持！
