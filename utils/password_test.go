package common_utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestComparePassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	err = ComparedPassword(hashedPassword, password)
	require.NoError(t, err)
}
