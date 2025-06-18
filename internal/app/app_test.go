package app_test

import (
	"errors"
	"in-memory-storage/internal/app"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplication_New(t *testing.T) {
	t.Run("it should return an error if initializing server fails", func(t *testing.T) {
		app, err := app.New("")
		expectedError := errors.New("missing port")
		assert.Equal(t, expectedError, err)
		assert.Nil(t, app)
	})

	t.Run("it should create a new Application instance", func(t *testing.T) {
		app, err := app.New("8080")
		assert.NoError(t, err)
		assert.NotNil(t, app)
	})
}
