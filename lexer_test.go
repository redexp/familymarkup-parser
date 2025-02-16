package parser

import (
	"errors"
	"iter"
	"os"
	"path/filepath"
	"runtime"
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
			Type:     TokenSpace,
			Text:     " ",
			Line:     0,
			Char:     3,
			CharsNum: 1,
		},
		{
			Type:     TokenSurname,
			Text:     "Fam",
			Line:     0,
			Char:     4,
			CharsNum: 3,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketLeft,
			Text:     "(",
			Line:     0,
			Char:     7,
			CharsNum: 1,
		},
		{
			Type:     TokenSurname,
			SubType:  TokenAlias,
			Text:     "Ali",
			Line:     0,
			Char:     8,
			CharsNum: 3,
		},
		{
			Type:     TokenEmptyLines,
			Text:     "\n \n",
			Line:     0,
			Char:     11,
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
			Char:     8,
			CharsNum: 1,
		},
		{
			Type:     TokenName,
			SubType:  TokenAlias,
			Text:     `P2`,
			Line:     2,
			Char:     9,
			CharsNum: 2,
		},
		{
			Type:     TokenBracket,
			SubType:  TokenBracketRight,
			Text:     `)`,
			Line:     2,
			Char:     11,
			CharsNum: 1,
		},
		{
			Type:     TokenArrow,
			Text:     `-->`,
			Line:     2,
			Char:     12,
			CharsNum: 3,
		},
		{
			Type:     TokenSpace,
			Text:     ` `,
			Line:     2,
			Char:     15,
			CharsNum: 1,
		},
		{
			Type:     TokenArrow,
			SubType:  TokenEqual,
			Text:     `=`,
			Line:     2,
			Char:     16,
			CharsNum: 1,
		},
		{
			Type:     TokenSpace,
			Text:     `  `,
			Line:     2,
			Char:     17,
			CharsNum: 2,
		},
		{
			Type:     TokenWord,
			Text:     `text`,
			Line:     2,
			Char:     19,
			CharsNum: 4,
		},
		{
			Type:     TokenSpace,
			Text:     ` `,
			Line:     2,
			Char:     23,
			CharsNum: 1,
		},
		{
			Type:     TokenUnknown,
			Text:     `слово?`,
			Line:     2,
			Char:     24,
			CharsNum: 6,
		},
		{
			Type:     TokenComment,
			Text:     `* sdf`,
			Line:     2,
			Char:     30,
			CharsNum: 5,
		},
		{
			Type:     TokenNewLine,
			Text:     "\n",
			Line:     2,
			Char:     35,
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
		} else if token.Line != item.Line {
			t.Errorf("Token %d - Invalid line %d, expect %d", i, token.Line, item.Line)
		} else if token.Char != item.Char {
			t.Errorf("Token %d - Invalid char %d, expect %d", i, token.Char, item.Char)
		} else if token.CharsNum != item.CharsNum {
			t.Errorf("Token %d - Invalid charsnum %d, expect %d", i, token.CharsNum, item.CharsNum)
		}
	}
}

func TestLexerFiles(t *testing.T) {
	for str, path := range testFilesIter(t) {
		tokens := Lexer(string(str))

		for _, token := range tokens {
			if token.Type == TokenInvalid {
				t.Errorf("File: %s, line: %d, char: %d, text: %s", path, token.Line, token.Char, token.Text)
			}
		}
	}
}

func TestAlias(t *testing.T) {
	str := testFile("alias.fml")

	tokens := Lexer(str)

	for _, token := range tokens {
		hasAlias := strings.HasSuffix(token.Text, "Alias")
		if (hasAlias && token.SubType != TokenAlias) || (!hasAlias && token.SubType == TokenAlias) {
			t.Errorf("token %d:%d, text: %s, type: %s, subtype: %s", token.Line, token.Char, token.Text, token.Type.String(), token.SubType.String())
		}
	}
}

func TestNameOnly(t *testing.T) {
	list := []Token{
		{
			Type: TokenSurname,
			Text: "Name",
		},
		{
			Type: TokenSpace,
			Text: " ",
		},
		{
			Type:    TokenBracket,
			SubType: TokenBracketLeft,
			Text:    "(",
		},
		{
			Type:    TokenSurname,
			SubType: TokenAlias,
			Text:    "Alias",
		},
		{
			Type:    TokenBracket,
			SubType: TokenBracketRight,
			Text:    ")",
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
		} else if token.Text != item.Text {
			t.Errorf("Token %d - Invalid text %s, expect %s", i, token.Text, item.Text)
		}
	}
}

func testFilesIter(t *testing.T) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		_, currentFile, _, ok := runtime.Caller(0)

		if !ok {
			t.Errorf("runtime.Caller !ok")
			return
		}

		currentDir := filepath.Dir(currentFile)
		root := filepath.Join(currentDir, "tests")
		count := 0

		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".fml") {
				return err
			}

			str, err := os.ReadFile(path)

			if err != nil {
				return err
			}

			count++

			if !yield(string(str), path) {
				return errors.New("stop iter")
			}

			return nil
		})

		if err != nil {
			t.Error(err)
		}

		if count == 0 {
			t.Errorf("count: %d", count)
		}
	}
}

func testFile(name string) string {
	_, currentFile, _, ok := runtime.Caller(0)

	if !ok {
		panic("runtime.Caller !ok")
	}

	currentDir := filepath.Dir(currentFile)
	path := filepath.Join(currentDir, "tests", name)

	str, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(str)
}
