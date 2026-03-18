package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"os"
	"slices"
	"time"
)

func main() {
	modelsToken := mustEnv("GH_PAT_CLASSIC_TOKEN")
	telegramToken := mustEnv("TELEGRAM_BOT_TOKEN")
	chatID := mustEnv("TELEGRAM_CHAT_ID")
	githubToken := os.Getenv("GITHUB_TOKEN")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	versions := readVersions()
	var hasError bool

	for _, repo := range githubRepos {
		knownSet := toSet(versions[repo])

		releases, err := retry(ctx, "fetch-releases", func(ctx context.Context) ([]*Release, error) {
			return fetchRecentReleases(ctx, repo, githubToken, releasesPerPage)
		})
		if err != nil {
			slog.Error("fetch releases failed", "repo", repo, "error", err)
			hasError = true
			continue
		}

		// All fetched tags form the new baseline; failed notifications
		// will be removed so they are retried on the next run.
		updatedTags := make(map[string]bool, len(releases))
		for _, r := range releases {
			updatedTags[r.TagName] = true
		}

		var newReleases []*Release
		for _, r := range releases {
			if !knownSet[r.TagName] {
				newReleases = append(newReleases, r)
			}
		}
		slices.Reverse(newReleases)

		if len(newReleases) == 0 {
			slog.Info("no new releases", "repo", repo)
		}

		for _, release := range newReleases {
			slog.Info("new release detected", "repo", repo, "tag", release.TagName, "prerelease", release.Prerelease)

			summary, err := retry(ctx, "interpret-release", func(ctx context.Context) (string, error) {
				return interpretRelease(ctx, modelsToken, repo, release)
			})
			if err != nil {
				slog.Error("interpret release failed", "repo", repo, "tag", release.TagName, "error", err)
				delete(updatedTags, release.TagName)
				hasError = true
				continue
			}

			msg := formatMessage(repo, release, summary)
			if err := retryDo(ctx, "send-telegram", func(ctx context.Context) error {
				return sendTelegram(ctx, telegramToken, chatID, msg)
			}); err != nil {
				slog.Error("send telegram failed", "repo", repo, "tag", release.TagName, "error", err)
				delete(updatedTags, release.TagName)
				hasError = true
				continue
			}
		}

		versions[repo] = sortedKeys(updatedTags)
	}

	if err := writeVersions(versions); err != nil {
		slog.Error("write versions failed", "error", err)
		os.Exit(1)
	}

	if hasError {
		slog.Error("completed with errors")
		os.Exit(1)
	}

	slog.Info("done")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("required env not set", "key", key)
		os.Exit(1)
	}
	return v
}

// readVersions supports both the legacy format ("repo": "tag") and
// the current format ("repo": ["tag1", "tag2", ...]).
func readVersions() map[string][]string {
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return make(map[string][]string)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return make(map[string][]string)
	}

	versions := make(map[string][]string, len(raw))
	for k, v := range raw {
		var arr []string
		if json.Unmarshal(v, &arr) == nil {
			versions[k] = arr
			continue
		}
		var s string
		if json.Unmarshal(v, &s) == nil {
			versions[k] = []string{s}
		}
	}
	return versions
}

func writeVersions(versions map[string][]string) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(versionFile, append(data, '\n'), 0644)
}

func formatMessage(repo string, r *Release, summary string) string {
	header := fmt.Sprintf(msgHeader, html.EscapeString(repo), html.EscapeString(r.TagName))
	if r.Prerelease {
		header += " <i>[Pre-release]</i>"
	}
	footer := fmt.Sprintf(msgFooter, r.HTMLURL)
	suffix := "\n\n" + footer

	escaped := html.EscapeString(summary)
	maxLen := telegramMsgLimit - len([]rune(header)) - len([]rune(suffix)) - 5
	runes := []rune(escaped)
	if len(runes) > maxLen {
		escaped = string(runes[:maxLen]) + "..."
	}
	return header + "\n\n" + escaped + suffix
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func sortedKeys(m map[string]bool) []string {
	ss := make([]string, 0, len(m))
	for k := range m {
		ss = append(ss, k)
	}
	slices.Sort(ss)
	return ss
}
