package token

var TokenMap map[string] string

type Token string

func VerifyToken(token Token) bool {
	return true
}

func CreateToken(location string)  map[string] string{

}

func (Token) getMap() map[string] string{
	return TokenMap
}
