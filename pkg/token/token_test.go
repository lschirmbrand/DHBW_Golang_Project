package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	tok := Token("test")

	valid, _ := Validate(tok)

	assert.True(t, valid)
}
