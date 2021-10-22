package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	tok1 := Token("test")
	tok2 := CreateToken("DE")
	tok3 := CreateToken("ENG")

valid1,_ := Validate(tok1)
valid2,_ := Validate(tok2)
valid3,_ := Validate(tok3)
	assert.False(t, valid1)
	assert.True(t, valid2)
	assert.True(t, valid3)

	//create Token to invalidate tok2
	tok4 := CreateToken("IT")
	valid2, _ = Validate(tok2)
	valid3, _ = Validate(tok3)
	valid4, _ := Validate(tok4)
	assert.False(t, valid2)
	assert.True(t, valid3)
	assert.True(t, valid4)
}
