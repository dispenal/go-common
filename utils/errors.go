package common_utils

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type AppError struct {
	Message    string
	StatusCode int
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("app error: status code %d, message %s", ae.StatusCode, ae.Message)
}

func CustomError(message string, statusCode int) error {
	return fmt.Errorf("|%s<->%d", message, statusCode)
}

func CustomErrorWithTrace(err error, message string, statusCode int) error {
	return fmt.Errorf("%s|%s<->%d", err.Error(), message, statusCode)
}

func PanicIfError(err error) {
	if err != nil {
		customError := strings.Split(err.Error(), "<->")
		message := customError[0]
		statusCode := 500

		if len(customError) > 1 {
			statusCode, _ = strconv.Atoi(customError[1])
		}

		appErr := AppError{
			Message:    message,
			StatusCode: statusCode,
		}
		panic(appErr)
	}
}

func PanicIfAppError(err error, message string, statusCode int) {
	if err != nil {
		customErr := CustomErrorWithTrace(err, message, statusCode)
		PanicIfError(customErr)
	}
}

func PanicAppError(message string, statusCode int) {
	customErr := CustomError(message, statusCode)
	PanicIfError(customErr)
}

func DeferCheck(function func() error) {
	if err := function(); err != nil {
		LogError("defer error", zap.Error(err))
	}
}

func LogIfError(err error) {
	if err != nil {
		LogError("error occurred", zap.Error(err))
	}
}

func LogAndPanicIfError(err error, message string) {
	if err != nil {
		errMsg := fmt.Sprintf("%s :%v", message, err)
		LogError(errMsg, zap.Error(err))
		panic(err)
	}
}
