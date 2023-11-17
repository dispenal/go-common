package common_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFilte(t *testing.T) {
	if _, err := CreateFile("test.txt"); err != nil {
		t.Error(err)
	}

	assert.Equal(t, true, CheckIfFileExists("test.txt"), "File created")
}

func TestFileExists(t *testing.T) {
	assert.Equal(t, true, CheckIfFileExists("test.txt"), "File exists")
}

func TestExtFile(t *testing.T) {
	assert.Equal(t, "text/plain; charset=utf-8", GetExt("test.txt"), "File extension")
}

func TestDeleteFile(t *testing.T) {
	if err := DeleteFile("test.txt"); err != nil {
		t.Error(err)
	}

	assert.Equal(t, false, CheckIfFileExists("test.txt"), "File deleted")
}

func TestDownloadFile(t *testing.T) {
	if err := DownloadFile("https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png", "google.png"); err != nil {
		t.Error(err)
	}

	assert.Equal(t, true, CheckIfFileExists("google.png"), "File downloaded")

	if err := DeleteFile("google.png"); err != nil {
		t.Error(err)
	}
}
