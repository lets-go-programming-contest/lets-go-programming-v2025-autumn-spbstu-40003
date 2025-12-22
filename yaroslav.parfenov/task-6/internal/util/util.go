package util

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errMockTemplate = errors.New("mock error")

func IsEmpty(str string) bool {
	return str == ""
}

func MakeError(msg string) error {
	if IsEmpty(msg) {
		return nil
	}

	return fmt.Errorf("%s %w", msg, errMockTemplate)
}

func AssertError(t *testing.T, response interface{}, err error, msg string) {
	t.Helper()

	assert.Nil(t, response)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), msg)
	}
}

func AssertNoError(t *testing.T, expected interface{}, response interface{}, err error) {
	t.Helper()

	if assert.NoError(t, err) {
		assert.Equal(t, expected, response)
	}
}
