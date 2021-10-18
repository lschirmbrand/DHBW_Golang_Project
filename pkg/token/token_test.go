package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	tok := Token("test")

	assert.True(t, VerifyToken(tok))
}
