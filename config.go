package main

// ---- GitHub repository settings ----

var githubRepos = []string{
	"ethereum/go-ethereum",
	"bnb-chain/bsc",
}

// ---- AI model settings ----

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

// ---- Telegram settings ----

const telegramMsgLimit = 4000

// ---- GitHub releases pagination ----

const releasesPerPage = 20

// ---- Version file ----

const versionFile = "last_versions.json"

// ---- Message template (HTML format for Telegram) ----

const (
	msgHeader = `<b>%s</b> 发布新版本 <b>%s</b>`
	msgFooter = `<a href="%s">查看完整发布说明</a>`
)
