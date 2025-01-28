package parser

//go:generate stringer -type=TokenType

type TokenType int

const (
	TokenName TokenType = iota
	TokenUnknown
	TokenWord
	TokenNum
	TokenArrow
	TokenEqual
	TokenPunctuation
	TokenPlus
	TokenComma
	TokenBracket
	TokenBracketLeft
	TokenBracketRight
	TokenComment
	TokenSpace
	TokenEmptyLines
	TokenNewLine
	TokenInvalid
)

type Token struct {
	Type     TokenType
	SubType  TokenType
	Offest   int
	Length   int
	Line     int
	Char     int
	CharsNum int
	Text     string
}

func (token *Token) End() int {
	return token.Offest + token.Length
}

func (token *Token) EndChar() int {
	return token.Char + token.CharsNum
}
