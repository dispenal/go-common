package common_utils

import (
	"encoding/base64"
	"io"
	"net/http"
)

func ImageUriToBase64(uri string) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	return base64Image, nil
}
