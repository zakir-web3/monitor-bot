# AI Reads

[![Go Version](https://img.shields.io/github/go-mod/go-version/zakir-web3/ai-reads)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Pages](https://img.shields.io/badge/Deep%20Analysis-GitHub%20Pages-brightgreen)](https://zakir-web3.github.io/ai-reads/)

An AI-powered release interpreter that monitors open-source repositories, generates intelligent analysis, and delivers them through two channels:

- **Telegram** — concise notification for quick awareness
- **[GitHub Pages](https://zakir-web3.github.io/ai-reads/)** — in-depth technical analysis website, auto-updated on every new release

## How It Works

```
New Release Detected
 │
 ├─→ AI Summary (concise) ──→ Telegram Push
 │
 └─→ AI Deep Analysis ──→ Markdown ──→ gh-pages branch ──→ GitHub Pages Website
```

1. Reads last known versions from `last_versions.json`
2. Fetches recent releases for each monitored repository via GitHub API
3. For each new release:
   - Calls AI to generate a **concise summary** → sends to Telegram
   - Calls AI to generate a **deep technical analysis** → publishes to [GitHub Pages](https://zakir-web3.github.io/ai-reads/)
4. Commits the updated `last_versions.json`

## Tracked Repositories

| Repository | Description |
|------------|-------------|
| [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) | Ethereum execution layer client (Geth) |
| [bnb-chain/bsc](https://github.com/bnb-chain/bsc) | BNB Smart Chain node |

Edit the `githubRepos` list in `config.go` to add or remove repositories.

## Getting Started

### Prerequisites

- Go 1.26 or later
- A [GitHub Models](https://github.com/marketplace/models) API token
- A Telegram Bot (create one via [@BotFather](https://t.me/BotFather))
- The chat ID of your target Telegram group or channel

### 1. Fork or Clone

```bash
git clone https://github.com/zakir-web3/ai-reads.git
cd ai-reads
```

### 2. Set Up Secrets

Go to your forked repository's **Settings > Secrets and variables > Actions**, and add the following secrets:

| Secret | Description |
|--------|-------------|
| `GH_PAT_CLASSIC_TOKEN` | GitHub personal access token (classic) for Models API |
| `TELEGRAM_BOT_TOKEN` | Telegram Bot token obtained from @BotFather |
| `TELEGRAM_CHAT_ID` | Target Telegram chat/group/channel ID |
| `GITHUB_TOKEN` | Automatically provided by GitHub Actions; used for API rate limiting and publishing to GitHub Pages |

> **Tip:** To find your Telegram chat ID, send a message to your bot, then visit `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates` and look for the `chat.id` field.

### 3. Enable GitHub Pages

Go to **Settings > Pages**, set source to `Deploy from a branch`, and select `gh-pages` / `/ (root)`.

### 4. Deploy via GitHub Actions

The workflow runs automatically at **08:18 CST (Beijing time)** every day. You can also trigger it manually from the **Actions** tab.

To change the schedule, edit the cron expression in `.github/workflows/ai-reads.yml`:

```yaml
on:
  schedule:
    - cron: '18 0 * * *'  # UTC time — adjust to your timezone
```

### 5. Run Locally (Optional)

```bash
export GH_PAT_CLASSIC_TOKEN="your-token"
export TELEGRAM_BOT_TOKEN="your-bot-token"
export TELEGRAM_CHAT_ID="your-chat-id"
export GITHUB_TOKEN="your-github-token"
export GITHUB_REPOSITORY="your-user/ai-reads"
go run .
```

## Project Structure

```
ai-reads/
├── config.go       # All configuration: repos, AI prompts, templates
├── main.go         # Entry point and orchestration
├── github.go       # GitHub API: fetch releases
├── ai.go           # AI model calls (concise + deep analysis)
├── telegram.go     # Telegram message delivery
├── pages.go        # Publish deep analysis to GitHub Pages
├── client.go       # Shared HTTP client and retry logic
└── .github/workflows/
    └── ai-reads.yml # Scheduled GitHub Actions workflow
```

## Customization

- **Tracked repos** — Edit the `githubRepos` list in `config.go`.
- **AI model / prompts** — Edit `modelName`, `systemPrompt`, `deepSystemPrompt` and prompt templates in `config.go`.
- **Message format** — Modify `msgHeader` / `msgFooter` in `config.go`. Telegram messages use HTML format.
- **Schedule** — Adjust the cron expression in `.github/workflows/ai-reads.yml`.

## License

MIT
