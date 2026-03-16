package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	modelsToken := mustEnv("GH_MODELS_TOKEN")
	telegramToken := mustEnv("TELEGRAM_BOT_TOKEN")
	chatID := mustEnv("TELEGRAM_CHAT_ID")

	versions := readVersions()

	for _, repo := range githubRepos {
		lastVersion := versions[repo]

		release, err := fetchLatestRelease(repo)
		if err != nil {
			log.Printf("[%s] fetch release: %v", repo, err)
			continue
		}

		if release.TagName == lastVersion {
			fmt.Printf("[%s] no new version (current: %s)\n", repo, lastVersion)
			continue
		}

		fmt.Printf("[%s] new version detected: %s (was: %s)\n", repo, release.TagName, lastVersion)

		summary, err := interpretRelease(modelsToken, repo, release)
		if err != nil {
			log.Printf("[%s] interpret release: %v", repo, err)
			continue
		}

		msg := formatMessage(repo, release, summary)
		if err := sendTelegram(telegramToken, chatID, msg); err != nil {
			log.Printf("[%s] send telegram: %v", repo, err)
			continue
		}

		versions[repo] = release.TagName
	}

	if err := writeVersions(versions); err != nil {
		log.Fatalf("write versions: %v", err)
	}

	fmt.Println("done")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env %s is required", key)
	}
	return v
}

func readVersions() map[string]string {
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return make(map[string]string)
	}
	var versions map[string]string
	if err := json.Unmarshal(data, &versions); err != nil {
		return make(map[string]string)
	}
	return versions
}

func writeVersions(versions map[string]string) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(versionFile, append(data, '\n'), 0644)
}

func formatMessage(repo string, r *Release, summary string) string {
	msg := fmt.Sprintf(messageTmpl, repo, r.TagName, summary, r.HTMLURL)
	if len([]rune(msg)) > telegramMsgLimit {
		runes := []rune(msg)
		msg = string(runes[:telegramMsgLimit]) + "..."
	}
	return msg
}
