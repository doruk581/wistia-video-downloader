package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func GetVideoIDFromCourseURL(courseURL string) (string, error) {
	re := regexp.MustCompile(`wvideo=([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(courseURL)
	if len(matches) < 2 {
		return "", errors.New("Wistia video ID not found")
	}
	return matches[1], nil
}

func balanceBraces(jsonStr string) string {
	openCount := strings.Count(jsonStr, "{")
	closeCount := strings.Count(jsonStr, "}")
	for openCount > closeCount {
		jsonStr += "}"
		closeCount++
	}
	return jsonStr
}

func GetIframeJSON(bodyString string) (string, error) {
	re := regexp.MustCompile(`(?s)W\.iframeInit\(\s*(\{.*\})\s*,\s*\{`)
	matches := re.FindStringSubmatch(bodyString)
	if len(matches) < 2 {
		return "", errors.New("JSON data not found (regex match failed)")
	}
	jsonStr := matches[1]
	jsonStr = balanceBraces(jsonStr)
	return jsonStr, nil
}

type WistiaResponse struct {
	Assets []struct {
		DisplayName string `json:"display_name"`
		URL         string `json:"url"`
		Height      int    `json:"height"`
	} `json:"assets"`
}

func Get1080pURL(videoID string) (string, error) {
	iframeURL := fmt.Sprintf("http://fast.wistia.net/embed/iframe/%s", videoID)
	resp, err := http.Get(iframeURL)
	if err != nil {
		return "", fmt.Errorf("error fetching iframe page: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading iframe page: %w", err)
	}
	bodyString := string(bodyBytes)

	jsonStr, err := GetIframeJSON(bodyString)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var wResponse WistiaResponse
	if err := json.Unmarshal([]byte(jsonStr), &wResponse); err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	for _, asset := range wResponse.Assets {
		if asset.Height == 1080 {
			url := asset.URL
			if strings.HasSuffix(url, ".bin") {
				url = strings.TrimSuffix(url, ".bin") + ".mp4"
			}
			return url, nil
		}
	}
	return "", errors.New("1080p video not found")
}

func DownloadFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("file download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("file download failed, HTTP status code: %d", resp.StatusCode)
	}

	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("file could not be created: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
