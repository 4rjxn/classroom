package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var googleExportMap = map[string]struct {
	ExportMIME string
	Extension  string
}{
	"application/vnd.google-apps.document":     {"application/pdf", ".pdf"},
	"application/vnd.google-apps.spreadsheet":  {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".xlsx"},
	"application/vnd.google-apps.presentation": {"application/pdf", ".pdf"},
	"application/vnd.google-apps.drawing":      {"application/pdf", ".pdf"},
}


type driveFileMeta struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Size     string `json:"size,omitempty"`
}


func isGoogleNativeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "application/vnd.google-apps.")
}

func DownloadDriveFile(token string, fileID string, fallbackName string, destDir string) (string, error) {
	
	metaURL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s?fields=id,name,mimeType,size", fileID)
	metaResp, err := DoGetRequest(metaURL, token)
	if err != nil {
		return "", fmt.Errorf("failed to get file metadata: %w", err)
	}
	defer metaResp.Body.Close()

	metaBytes, err := io.ReadAll(metaResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read file metadata: %w", err)
	}

	var meta driveFileMeta
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return "", fmt.Errorf("failed to parse file metadata: %w", err)
	}

	
	fileName := meta.Name
	if fileName == "" {
		fileName = fallbackName
	}
	if fileName == "" {
		fileName = "download"
	}

	
	var downloadURL string

	if exportInfo, ok := googleExportMap[meta.MimeType]; ok {
		
		downloadURL = fmt.Sprintf(
			"https://www.googleapis.com/drive/v3/files/%s/export?mimeType=%s",
			fileID, exportInfo.ExportMIME,
		)
		
		if !strings.HasSuffix(strings.ToLower(fileName), exportInfo.Extension) {
			fileName += exportInfo.Extension
		}
	} else if isGoogleNativeType(meta.MimeType) {
		
		return "", fmt.Errorf("cannot download Google %s files — open in browser instead",
			strings.TrimPrefix(meta.MimeType, "application/vnd.google-apps."))
	} else {
		
		downloadURL = fmt.Sprintf(
			"https://www.googleapis.com/drive/v3/files/%s?alt=media",
			fileID,
		)
	}

	
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create download directory %s: %w", destDir, err)
	}

	
	savePath := uniqueFilePath(filepath.Join(destDir, sanitizeFileName(fileName)))

	
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "classroom-cli/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download failed with HTTP %d", resp.StatusCode)
	}

	
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", savePath, err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		
		_ = os.Remove(savePath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return savePath, nil
}


func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}


func uniqueFilePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)

	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s (%d)%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}

	return path
}
