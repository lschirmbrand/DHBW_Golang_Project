package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	tok1 := Token("test")
	tok2 := CreateToken("DE")
	tok3 := CreateToken("ENG")

	assert.False(t, VerifyToken(tok1))
	assert.True(t, VerifyToken(tok2))
	assert.True(t, VerifyToken(tok3))

	//create Token to invalidate tok2
	tok4 := CreateToken("IT")

	assert.False(t, VerifyToken(tok2))
	assert.True(t, VerifyToken(tok3))
	assert.True(t, VerifyToken(tok4))
}
