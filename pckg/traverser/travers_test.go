package traverser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TraverserProperlyWorks(t *testing.T) {
}

func Test_IsValidPath(t *testing.T) {
	assert.True(t, isValidPath("/", &[]string{""}))
}

func Test_IsValidPathProperlyUsesIgnores(t *testing.T) {
	assert.False(t, isValidPath("/some/string", &[]string{""}))
}
