package token

import (
	"math/rand"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
type Token string
var tokenMap = map[Token]string{}
func Validator(token Token) (bool, String) {
	if location, ok := tokenMap[token]; ok {
		return true, location
	} else {
		return false, ""
	}
}

func CreateToken(location string) Token {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	if len(tokenMap) == 2 {
		deleteFirstToken()
	}

	newToken := Token(b)
	tokenMap[newToken] = location
	return newToken
}

func deleteFirstToken() {
	for k := range tokenMap {
		delete(tokenMap, k)
		break
	}
}
