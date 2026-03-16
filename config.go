package main

// ---- GitHub 仓库配置 ----

var githubRepos = []string{
	"ethereum/go-ethereum",
}

// ---- AI 模型配置 ----

const (
	modelsEndpoint = "https://models.inference.ai.azure.com/chat/completions"
	modelName      = "gpt-4o-mini"
	systemPrompt   = "你是一个区块链和开源技术专家，擅长解读技术发布说明。请用简洁清晰的中文回答，总字数控制在800字以内。"
	userPromptTmpl = `请用中文解读以下 %s 版本发布内容，简明扼要，分以下几点：
1. 版本概述
2. 主要新特性
3. 重要变更或不兼容改动
4. 安全修复（如有）
5. 对运营者的建议

版本：%s
发布内容：
%s`
)

// ---- Telegram 配置 ----

const telegramMsgLimit = 4000

// ---- 版本文件 ----

const versionFile = "last_versions.json"

// ---- 消息模板 ----

const messageTmpl = "%s 发布新版本 %s\n\n%s\n\n查看完整发布说明：%s"
