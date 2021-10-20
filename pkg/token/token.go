package token

var TokenMap map[string] string

type Token string

type Validator func(t Token) (bool, string)

func Validate(t Token) (bool, string) {
	return true, "lol"
}

func CreateToken(location string)  map[string] string{

}

func (Token) getMap() map[string] string{
	return TokenMap
}
