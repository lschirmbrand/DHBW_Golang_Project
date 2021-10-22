package token

type Token string

type Validator func(t Token) (bool, string)

func Validate(t Token) (bool, string) {
	return true, "lol"
}
