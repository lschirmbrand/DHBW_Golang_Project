package token

import (
	"DHBW_Golang_Project/internal/location"
	"math/rand"
)

/*
	Erstellt von: 	9514094
	Created by:		9514094

	also: 4775194, 8864957
*/

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Token string
type Validator func(t Token, l location.Location) bool

var tokenMap = map[location.Location][]Token{}

//validate Token
func Validate(expToken Token, tokenLocation location.Location) bool {
	for _, actToken := range tokenMap[tokenLocation] {
		if actToken == expToken {
			return true
		}
	}
	return false
}

//create Token
func CreateToken(tokenLocation location.Location) Token {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	if len(tokenMap[tokenLocation]) == 2 {
		tokenMap[tokenLocation] = []Token{tokenMap[tokenLocation][1]}
	}

	newToken := Token(b)
	tokenMap[tokenLocation] = append(tokenMap[tokenLocation], newToken)
	return newToken
}
