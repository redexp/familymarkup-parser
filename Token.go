package parser

//go:generate stringer -type=TokenType,ErrType

type TokenType int

const (
	TokenName TokenType = iota + 1
	TokenSurname
	TokenAlias
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

type ErrType int

const (
	ErrUnexpected ErrType = iota + 1
)

type Token struct {
	Type    TokenType
	SubType TokenType
	ErrType
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
