package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type createFileRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch"`
}

// publishToPages pushes a deep-analysis markdown file to the gh-pages branch
// via the GitHub Contents API. The file is placed under _posts/ so that Jekyll
// automatically includes it in the site index.
// lang should be "zh" or "en".
func publishToPages(ctx context.Context, token, pagesRepo, repo string, r *Release, analysis, lang string) error {
	now := time.Now()
	slug := sanitizeSlug(repo, r.TagName)
	path := fmt.Sprintf("_posts/%s-%s-%s.md", now.Format("2006-01-02"), slug, lang)

	md := buildPostMarkdown(repo, r, analysis, now, lang)

	reqBody := createFileRequest{
		Message: fmt.Sprintf("docs: add %s analysis for %s %s", lang, repo, r.TagName),
		Content: base64.StdEncoding.EncodeToString([]byte(md)),
		Branch:  pagesBranch,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", pagesRepo, path)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github contents api status %d: %s", resp.StatusCode, respBody)
	}
	return nil
}

func sanitizeSlug(repo, tag string) string {
	s := strings.ToLower(repo + "-" + tag)
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return s
}

func buildPostMarkdown(repo string, r *Release, analysis string, now time.Time, lang string) string {
	var b strings.Builder
	slug := sanitizeSlug(repo, r.TagName)
	date := now.Format("2006-01-02")

	// Jekyll front matter
	b.WriteString("---\n")
	b.WriteString("layout: default\n")
	fmt.Fprintf(&b, "lang: %s\n", lang)
	if lang == "en" {
		fmt.Fprintf(&b, "title: \"%s %s Deep Dive\"\n", repo, r.TagName)
	} else {
		fmt.Fprintf(&b, "title: \"%s %s 版本深度解读\"\n", repo, r.TagName)
	}
	fmt.Fprintf(&b, "date: %s\n", date)
	fmt.Fprintf(&b, "repo: %s\n", repo)
	fmt.Fprintf(&b, "tag: %s\n", r.TagName)
	b.WriteString("---\n\n")

	// Language switch link
	if lang == "en" {
		altPath := fmt.Sprintf("/%s-%s-zh", date, slug)
		fmt.Fprintf(&b, "> [中文版](%s)\n\n", altPath)
	} else {
		altPath := fmt.Sprintf("/%s-%s-en", date, slug)
		fmt.Fprintf(&b, "> [English Version](%s)\n\n", altPath)
	}

	// Page content
	if lang == "en" {
		fmt.Fprintf(&b, "# %s %s Deep Dive\n\n", repo, r.TagName)
		fmt.Fprintf(&b, "> Analysis date: %s | [View original release notes](%s)\n\n", date, r.HTMLURL)
	} else {
		fmt.Fprintf(&b, "# %s %s 版本深度解读\n\n", repo, r.TagName)
		fmt.Fprintf(&b, "> 分析日期: %s | [查看原始发布说明](%s)\n\n", date, r.HTMLURL)
	}
	b.WriteString("---\n\n")
	b.WriteString(analysis)
	b.WriteString("\n")

	return b.String()
}
