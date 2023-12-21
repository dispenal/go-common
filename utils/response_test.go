package common_utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUpRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestGenerateSuccessResp(t *testing.T) {
	recorder := setUpRecorder()
	data := "Test Success"
	successStatusCode := 201

	GenerateJsonResponse(recorder, data, successStatusCode, "Test Success")

	var resp Response[string]
	NewDecoder(recorder.Body).Decode(&resp)

	assert.Equal(t, "Test Success", resp.Message)
	assert.Equal(t, http.StatusText(successStatusCode), resp.Status)
	assert.Equal(t, successStatusCode, recorder.Result().StatusCode)
	assert.Equal(t, "Test Success", resp.Data)
}

func TestGenerateErrorResp(t *testing.T) {
	recorder := setUpRecorder()
	data := "Test Failed"
	errorStatusCode := 400

	GenerateJsonResponse(recorder, data, errorStatusCode, "Test Failed")

	var resp Response[string]
	NewDecoder(recorder.Body).Decode(&resp)

	assert.Equal(t, "Test Failed", resp.Message)
	assert.Equal(t, http.StatusText(errorStatusCode), resp.Status)
	assert.Equal(t, 400, recorder.Result().StatusCode)
	assert.Equal(t, "Test Failed", resp.Data)
}
