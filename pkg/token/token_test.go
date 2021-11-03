package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateToken(t *testing.T) {
	//create 3 Token and check validation
	tok1 := Token("test")
	tok2 := CreateToken("DE")
	tok3 := CreateToken("ENG")

	valid1 := Validate(tok1, "test")
	valid2 := Validate(tok2, "DE")
	valid3 := Validate(tok3, "ENG")

	assert.False(t, valid1)
	assert.True(t, valid2)
	assert.True(t, valid3)

	//create Token to invalidate tok2
	tok4 := CreateToken("DE")
	tok5 := CreateToken("DE")
	valid2 = Validate(tok2, "DE")
	valid3 = Validate(tok3, "ENG")
	valid4 := Validate(tok4, "DE")
	valid5 := Validate(tok5, "DE")
	assert.False(t, valid2)
	assert.True(t, valid3)
	assert.True(t, valid4)
	assert.True(t, valid5)
}
