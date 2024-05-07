package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeError_ReturnsWithCorrectMessage(t *testing.T) {
	assert := assert.New(t)

	var err error
	assert.NotPanics(
		func() { err = MakeError("Test", "Error!") },
		"Should not panic when forming error!",
	)

	assert.Equal(
		"Test: Error!",
		err.Error(),
		"Error messages must be the same!",
	)
}

func TestGetErrorMessage_ReturnsCorrectMessage(t *testing.T) {
	assert := assert.New(t)

	var errMsg string
	assert.NotPanics(
		func() {
			errMsg = GetErrorMessage("Test", "Error!")
		},
		"Should not panic when phrasing an error message!",
	)

	assert.Equal(
		"Test: Error!",
		errMsg,
		"Error messages must be the same!",
	)
}
