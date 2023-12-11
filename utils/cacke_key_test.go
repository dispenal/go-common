package common_utils

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const TEST_KEY = "TEST"
const TESTS_KEY = "TESTS"

type testArgs struct {
	Name   string
	ID     int
	Status string
}

func TestCacheKey(t *testing.T) {
	id := uuid.New()
	arg := testArgs{
		ID:     1,
		Name:   "Testing",
		Status: "active",
	}
	key := BuildCacheKey(TEST_KEY, id.String(), "funcName", arg, arg, testArgs{})

	assert.Equal(t, fmt.Sprintf("%s-%s-funcName|Name:%s,ID:%d,Status:%s|Name:%s,ID:%d,Status:%s|ID:0", TEST_KEY, id, arg.Name,
		arg.ID, arg.Status, arg.Name, arg.ID, arg.Status), key)
}

func TestBuildPrefixKey(t *testing.T) {
	id := uuid.NewString()
	prefixKey := BuildPrefixKey(TEST_KEY, id, "TESTING")

	assert.Equal(t, fmt.Sprintf("%s-%s-TESTING", TEST_KEY, id), prefixKey)
}
