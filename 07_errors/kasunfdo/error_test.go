package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrCodeString(t *testing.T) {
	assert.Equal(t, "invalid input: %v", ErrInvalid.String())
	assert.Equal(t, "not found: %v", ErrNotFound.String())
	assert.Equal(t, "internal error", ErrInternal.String())

	var ErrFoo ErrCode = 900
	assert.Equal(t, "error occurred", ErrFoo.String())
}
