package main

// ---- GitHub repository settings ----

var githubRepos = []string{
	"ethereum/go-ethereum",
	"bnb-chain/bsc",
	"openclaw/openclaw",
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

// ---- Deep analysis AI settings (for GitHub Pages) ----

const (
	deepSystemPromptZH = `你是一位资深的区块链基础设施工程师和技术分析师，拥有丰富的节点运维和底层协议开发经验。
你擅长深入解读开源区块链项目的版本发布说明，能够从技术实现、安全影响、生态系统演进等多个维度进行专业分析。
请用清晰、专业的中文撰写深度技术分析报告，报告应当结构清晰、逻辑严密，适合区块链技术从业者和节点运营者阅读。`

	deepUserPromptTmplZH = `请对以下 %s 的 %s 版本发布内容进行全面深入的技术分析，生成一份专业的 Markdown 格式分析报告。

请从以下维度逐一展开分析：

1. **版本概述**
   - 版本号语义分析（主版本/次版本/补丁版本的含义）
   - 版本发布背景和定位

2. **核心变更详解**
   - 逐一分析关键的功能变更和改进
   - 解释变更的技术背景和实现原理
   - 评估每项变更对系统行为的具体影响

3. **性能与优化**
   - 性能相关改进的技术细节
   - 预期的性能收益和适用场景
   - 对资源消耗（CPU、内存、存储、带宽）的影响

4. **安全分析**
   - 安全修复的详细分析
   - 漏洞类型、严重程度和攻击面评估
   - 修复前后的安全态势对比

5. **不兼容变更与迁移**
   - Breaking Changes 的完整列表和影响分析
   - 从旧版本升级的详细迁移指南
   - 配置文件或 API 变更说明

6. **对节点运营者的影响**
   - 是否需要立即升级及优先级评估
   - 升级步骤和注意事项
   - 回滚方案（如适用）
   - 对质押/验证/出块的影响

7. **对开发者的影响**
   - RPC/API 接口变更
   - SDK/库的兼容性
   - 智能合约相关的变更

8. **生态影响与技术趋势**
   - 此版本在项目路线图中的位置
   - 对更广泛区块链生态的意义

如果某个维度在本次发布中无相关内容，可以简要说明并跳过。请确保分析准确、专业且有深度。

发布内容：
%s`

	deepSystemPromptEN = `You are a senior blockchain infrastructure engineer and technical analyst with extensive experience in node operations and low-level protocol development.
You specialize in providing in-depth analysis of open-source blockchain project release notes, covering multiple dimensions such as technical implementation, security implications, and ecosystem evolution.
Please write a thorough, professional technical analysis report in clear English. The report should be well-structured, logically rigorous, and suitable for blockchain professionals and node operators.`

	deepUserPromptTmplEN = `Please perform a comprehensive technical analysis of the following %s release %s and generate a professional Markdown analysis report.

Analyze from each of the following dimensions:

1. **Release Overview**
   - Semantic version analysis (major/minor/patch implications)
   - Release background and positioning

2. **Core Changes in Detail**
   - Analyze each key functional change and improvement
   - Explain the technical background and implementation rationale
   - Assess the specific impact of each change on system behavior

3. **Performance & Optimization**
   - Technical details of performance-related improvements
   - Expected performance gains and applicable scenarios
   - Impact on resource consumption (CPU, memory, storage, bandwidth)

4. **Security Analysis**
   - Detailed analysis of security fixes
   - Vulnerability type, severity, and attack surface assessment
   - Security posture comparison before and after the fix

5. **Breaking Changes & Migration**
   - Complete list of breaking changes and impact analysis
   - Detailed migration guide from previous versions
   - Configuration file or API change notes

6. **Impact on Node Operators**
   - Whether an immediate upgrade is required and priority assessment
   - Upgrade steps and considerations
   - Rollback plan (if applicable)
   - Impact on staking/validation/block production

7. **Impact on Developers**
   - RPC/API interface changes
   - SDK/library compatibility
   - Smart contract related changes

8. **Ecosystem Impact & Technical Trends**
   - Position of this release in the project roadmap
   - Significance to the broader blockchain ecosystem

If a dimension has no relevant content in this release, briefly note it and skip. Ensure the analysis is accurate, professional, and insightful.

Release content:
%s`
)

// ---- Telegram settings ----

const telegramMsgLimit = 4000

// ---- GitHub releases pagination ----

const releasesPerPage = 5

// ---- Version file ----

const versionFile = "last_versions.json"

// ---- Message template (HTML format for Telegram) ----

const (
	msgHeader = `<b>%s</b> 发布新版本 <b>%s</b>`
	msgFooter = `<a href="%s">查看完整发布说明</a>`
)

// ---- GitHub Pages settings ----

const pagesBranch = "gh-pages"
