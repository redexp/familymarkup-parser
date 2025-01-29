package parser

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	list := []Token{
		{
			Type:     TokenSurname,
			Text:     "Fam",
			Line:     0,
			Char:     0,
			CharsNum: 3,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketLeft,
			Text:     "(",
			Line:     0,
			Char:     3,
			CharsNum: 1,
		},
		{
			Type:     TokenSurname,
			SubType:  TokenAlias,
			Text:     "Ali",
			Line:     0,
			Char:     4,
			CharsNum: 3,
		},
		{
			Type:     TokenEmptyLines,
			Text:     "\n \n",
			Line:     0,
			Char:     7,
			CharsNum: 3,
		},
		{
			Type:     TokenName,
			Text:     `Тест`,
			Line:     2,
			Char:     0,
			CharsNum: 4,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketLeft,
			Text:     `(`,
			Line:     2,
			Char:     4,
			CharsNum: 1,
		},
		{
			Type:     TokenName,
			SubType:  TokenAlias,
			Text:     `P-1`,
			Line:     2,
			Char:     5,
			CharsNum: 3,
		},
		{
			Type:     TokenPunctuation,
			SubType:  TokenComma,
			Text:     `,`,
			Line:     2,
			Char:     6,
			CharsNum: 1,
		},
		{
			Type:     TokenName,
			SubType:  TokenAlias,
			Text:     `P2`,
			Line:     2,
			Char:     7,
			CharsNum: 2,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketRight,
			Text:     `)`,
			Line:     2,
			Char:     9,
			CharsNum: 1,
		},
		{
			Type:     TokenArrow,
			Text:     `-->`,
			Line:     2,
			Char:     10,
			CharsNum: 3,
		},
		{
			Type:     TokenSpace,
			Text:     ` `,
			Line:     2,
			Char:     13,
			CharsNum: 1,
		},
		{
			Type:     TokenArrow,
			SubType:  TokenEqual,
			Text:     `=`,
			Line:     2,
			Char:     14,
			CharsNum: 1,
		},
		{
			Type:     TokenSpace,
			Text:     `  `,
			Line:     2,
			Char:     15,
			CharsNum: 2,
		},
		{
			Type:     TokenWord,
			Text:     `text`,
			Line:     2,
			Char:     17,
			CharsNum: 4,
		},
		{
			Type:     TokenSpace,
			Text:     ` `,
			Line:     2,
			Char:     21,
			CharsNum: 1,
		},
		{
			Type:     TokenUnknown,
			Text:     `слово?`,
			Line:     2,
			Char:     22,
			CharsNum: 6,
		},
		{
			Type:     TokenComment,
			Text:     `* sdf`,
			Line:     2,
			Char:     28,
			CharsNum: 5,
		},
		{
			Type:     TokenNewLine,
			Text:     "\n",
			Line:     2,
			Char:     33,
			CharsNum: 1,
		},
		{
			Type:     TokenName,
			Text:     "Test-Name",
			Line:     3,
			Char:     0,
			CharsNum: 9,
		},
		{
			Type:     TokenSpace,
			Text:     " ",
			Line:     3,
			Char:     9,
			CharsNum: 1,
		},
		{
			Type:     TokenArrow,
			Text:     "-",
			Line:     3,
			Char:     10,
			CharsNum: 1,
		},
		{
			Type:     TokenSpace,
			Text:     " ",
			Line:     3,
			Char:     11,
			CharsNum: 1,
		},
		{
			Type:     TokenWord,
			Text:     "word-word слово",
			Line:     3,
			Char:     12,
			CharsNum: 15,
		},
		{
			Type:     TokenSpace,
			Text:     " ",
			Line:     3,
			Char:     27,
			CharsNum: 1,
		},
		{
			Type:     TokenUnknown,
			Text:     "Имя word-word?",
			Line:     3,
			Char:     28,
			CharsNum: 14,
		},
		{
			Type:     TokenSpace,
			Text:     " ",
			Line:     3,
			Char:     42,
			CharsNum: 1,
		},
		{
			Type:     TokenUnknown,
			Text:     "?",
			Line:     3,
			Char:     43,
			CharsNum: 1,
		},
		{
			Type:     TokenNewLine,
			Text:     "\n",
			Line:     3,
			Char:     44,
			CharsNum: 1,
		},
		{
			Type:     TokenName,
			Text:     "Name",
			Line:     4,
			Char:     0,
			CharsNum: 4,
		},
		{
			Type:     TokenSpace,
			Text:     " ",
			Line:     4,
			Char:     4,
			CharsNum: 1,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketLeft,
			Text:     "(",
			Line:     4,
			Char:     5,
			CharsNum: 1,
		},
		{
			Type:     TokenName,
			SubType:  TokenAlias,
			Text:     "Al",
			Line:     4,
			Char:     6,
			CharsNum: 2,
		},
		{
			Type:     TokenPunctuation,
			SubType:  TokenComma,
			Text:     ",",
			Line:     4,
			Char:     8,
			CharsNum: 1,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketRight,
			Text:     ")",
			Line:     4,
			Char:     9,
			CharsNum: 1,
		},
		{
			Type:     TokenSurname,
			Text:     "Sur",
			Line:     4,
			Char:     10,
			CharsNum: 3,
		},
		{
			Type:     TokenNewLine,
			Text:     "\n",
			Line:     4,
			Char:     13,
			CharsNum: 1,
		},
		{
			Type:     TokenNum,
			Text:     "1.",
			Line:     5,
			Char:     0,
			CharsNum: 2,
		},
		{
			Type:     TokenName,
			Text:     "Name",
			Line:     5,
			Char:     2,
			CharsNum: 4,
		},
		{
			Type:     TokenSpace,
			Text:     " ",
			Line:     5,
			Char:     6,
			CharsNum: 1,
		},
		{
			Type:     TokenSurname,
			Text:     "Sur",
			Line:     5,
			Char:     7,
			CharsNum: 3,
		},
		{
			Type:     TokenInvalid,
			Text:     "@%",
			Line:     5,
			Char:     10,
			CharsNum: 2,
		},
	}

	var b strings.Builder

	for _, token := range list {
		b.WriteString(token.Text)
	}

	tokens := Lexer(b.String())

	if len(tokens) != len(list) {
		t.Errorf("Invalid tokens length %d, expect %d", len(tokens), len(list))
		return
	}

	for i, token := range tokens {
		item := list[i]

		if token.Type != item.Type {
			t.Errorf("Token %d - Invalid type %s, expect %s", i, token.Type.String(), item.Type.String())
		} else if token.SubType != item.SubType {
			t.Errorf("Token %d - Invalid subtype %s, expect %s", i, token.SubType.String(), item.SubType.String())
		} else if token.Text != item.Text {
			t.Errorf("Token %d - Invalid text %s, expect %s", i, token.Text, item.Text)
		}
	}
}
