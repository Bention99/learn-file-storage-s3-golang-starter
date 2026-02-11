package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os/exec"
)

type stream struct {
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type probe struct {
	Streams []stream `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
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

	var out probe
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return "", fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	var w, h int
	for _, s := range out.Streams {
		if s.CodecType == "video" && s.Width > 0 && s.Height > 0 {
			w, h = s.Width, s.Height
			break
		}
	}
	if w <= 0 || h <= 0 {
		return "", errors.New("no video stream with width/height found in ffprobe output")
	}

	ratio := float64(w) / float64(h)

	const eps = 0.02
	if math.Abs(ratio-(16.0/9.0)) <= eps {
		return "16:9", nil
	}
	if math.Abs(ratio-(9.0/16.0)) <= eps {
		return "9:16", nil
	}
	return "other", nil
}