package common_utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func ImageUriToBase64(uri string) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		LogError(fmt.Sprintf("Error fetching image: %s", err))
		return "", err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		LogError(fmt.Sprintf("Error reading image data: %s", err))
		return "", err
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	return base64Image, nil
}
