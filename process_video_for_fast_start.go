package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	outputFilePath := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFilePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("ffprobe failed: %w: %s", err, stderr.String())
		}
		return "", fmt.Errorf("ffprobe failed: %w", err)
	}

	return outputFilePath, nil
}