package token

import (
	"DHBW_Golang_Project/pkg/location"
	"math/rand"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Token string
type Validator func(t Token, l location.Location) bool

var tokenMap = map[location.Location][]Token{}

//validate Token with actual location
func Validate(expToken Token, tokenLocation location.Location) bool {
	for _, actToken := range tokenMap[tokenLocation] {
		//check all token at location-map
		if actToken == expToken {
			return true
		}
	}
	return false
}

func CreateToken(tokenLocation location.Location) Token {
	//create 10 byte-array and set value to random letters
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	//invalidate first token
	if len(tokenMap[tokenLocation]) == 2 {
		tokenMap[tokenLocation] = []Token{tokenMap[tokenLocation][1]}
	}

	newToken := Token(b)
	//add new token to the validated list
	tokenMap[tokenLocation] = append(tokenMap[tokenLocation], newToken)
	return newToken
}
