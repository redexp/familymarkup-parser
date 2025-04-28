package parser

//go:generate stringer -type=TokenType,ErrType

type TokenType int

const (
	TokenName         TokenType = 1 << iota // 1
	TokenSurname                            // 2
	TokenAlias                              // 4
	TokenUnknown                            // 8
	TokenWord                               // 16
	TokenNum                                // 32
	TokenArrow                              // 64
	TokenEqual                              // 128
	TokenPunctuation                        // 256
	TokenPlus                               // 512
	TokenComma                              // 1024
	TokenBracket                            // 2048
	TokenBracketLeft                        // 4096
	TokenBracketRight                       // 8192
	TokenComment                            // 16384
	TokenSpace                              // 32768
	TokenEmptyLines                         // 65536
	TokenNewLine                            // 131072
	TokenInvalid                            // 262144
)

type ErrType int

const (
	ErrUnexpected ErrType = 1 << iota
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

func (token *Token) Loc() Loc {
	return Loc{
		Start: Position{
			Line: token.Line,
			Char: token.Char,
		},
		End: Position{
			Line: token.Line,
			Char: token.EndChar(),
		},
	}
}
