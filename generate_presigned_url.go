package main

import (
	"context"
	"time"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s3Client)

	out, err := presigner.PresignGetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
		s3.WithPresignExpires(expireTime),
	)
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
    	return video, nil
	}
	parts := strings.Split(*video.VideoURL, ",")
	if len(parts) != 2 {
		return video, nil
	}

	bucket := parts[0]
	key := parts[1]

	signedURL, err := generatePresignedURL(
		cfg.s3Client,
		bucket,
		key,
		time.Hour,
	)
	if err != nil {
		return video, err
	}

	video.VideoURL = &signedURL
	return video, nil
}