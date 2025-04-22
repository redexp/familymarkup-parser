package parser

import (
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

type Rule struct {
	Type    TokenType
	SubType TokenType
	Regexp  *regexp.Regexp
}

var rules []Rule = []Rule{
	{
		Type:    TokenEmptyLines,
		SubType: TokenNewLine,
		Regexp:  regexp.MustCompile(`^\r?\n[\r\n\t ]*\n`),
	},
	{
		Type:    TokenNewLine,
		SubType: TokenNewLine,
		Regexp:  regexp.MustCompile(`^\r?\n`),
	},
	{
		Type:   TokenSpace,
		Regexp: regexp.MustCompile(`^ +`),
	},
	{
		Type:   TokenArrow,
		Regexp: regexp.MustCompile(`^<?-+>`),
	},
	{
		Type:   TokenArrow,
		Regexp: regexp.MustCompile(`^<-+>?`),
	},
	{
		Type:   TokenUnknown,
		Regexp: regexp.MustCompile(`^[\p{L}\d"'\-.]*\?`),
	},
	{
		Type:   TokenName,
		Regexp: regexp.MustCompile(`^\p{Lu}[\p{L}\d'.]*(-+[\p{L}\d'.]+)*`),
	},
	{
		Type:   TokenWord,
		Regexp: regexp.MustCompile(`^['"][\p{Ll}\d'".]+(-+[\p{Ll}\d'".]+)*`),
	},
	{
		Type:   TokenWord,
		Regexp: regexp.MustCompile(`^\p{Ll}[\p{Ll}\d'".]*(-+[\p{Ll}\d'".]+)*`),
	},
	{
		Type:   TokenArrow,
		Regexp: regexp.MustCompile(`^-+`),
	},
	{
		Type:    TokenArrow,
		SubType: TokenEqual,
		Regexp:  regexp.MustCompile(`^=+`),
	},
	{
		Type:    TokenPunctuation,
		SubType: TokenPlus,
		Regexp:  regexp.MustCompile(`^\+`),
	},
	{
		Type:    TokenPunctuation,
		SubType: TokenComma,
		Regexp:  regexp.MustCompile(`^,`),
	},
	{
		Type:    TokenBracket,
		SubType: TokenBracketLeft,
		Regexp:  regexp.MustCompile(`^\(`),
	},
	{
		Type:    TokenBracket,
		SubType: TokenBracketRight,
		Regexp:  regexp.MustCompile(`^\)`),
	},
	{
		Type:   TokenNum,
		Regexp: regexp.MustCompile(`^\d+\.?`),
	},
	{
		Type:   TokenComment,
		Regexp: regexp.MustCompile(`^[/#*][^\r\n]*`),
	},
}

func Lexer(src string) (list []*Token) {
	offset := 0
	length := len(src)
	line := 0
	chars := 0
	var prev *Token
	leftOpen := false
	hasFamilyName := atomic.Bool{}
	var wg sync.WaitGroup

	for offset < length {
		var token *Token
		offsetSrc := src[offset:]

		for _, rule := range rules {
			match := rule.Regexp.FindStringSubmatch(offsetSrc)

			if match == nil {
				continue
			}

			text := match[0]

			token = &Token{
				Type:     rule.Type,
				SubType:  rule.SubType,
				Offest:   offset,
				Length:   len(text),
				Text:     text,
				Line:     line,
				Char:     chars,
				CharsNum: utf8.RuneCountInString(text),
			}

			break
		}

		if token == nil {
			if prev != nil && prev.Type == TokenInvalid {
				_, size := utf8.DecodeRuneInString(src[prev.Offest:prev.End()])
				prev.Length += size
				prev.CharsNum++
				prev.Text = src[prev.Offest:prev.End()]

				token = prev
			} else {
				_, size := utf8.DecodeRuneInString(src[offset:])
				token = &Token{
					Type:     TokenInvalid,
					Offest:   offset,
					Length:   size,
					Line:     line,
					Char:     chars,
					CharsNum: 1,
					Text:     src[offset : offset+size],
				}
			}
		}

		switch token.Type {
		case TokenName:
			if leftOpen {
				token.SubType = TokenAlias
			} else {
				checkSurname(list, token)
			}

		case TokenUnknown:
			list = mergeUnknown(list, token, src)

		case TokenWord:
			list = mergeWords(list, token, src)

		case TokenBracket:
			leftOpen = token.SubType == TokenBracketLeft

		case TokenNewLine:
			line++

		case TokenEmptyLines:
			line += strings.Count(token.Text, "\n")

			wg.Add(1)
			go func(list []*Token) {
				defer wg.Done()

				has := checkFamilyName(list)

				if has {
					hasFamilyName.Store(true)
				}
			}(list)
		}

		if leftOpen {
			switch token.Type {
			case TokenName, TokenSurname, TokenBracket, TokenSpace, TokenWord, TokenInvalid:
				// ok

			case TokenPunctuation:
				if token.SubType != TokenComma {
					leftOpen = false
				}

			default:
				leftOpen = false
			}
		}

		if token != prev {
			list = append(list, token)
			prev = token
		}

		if token.Type == TokenNewLine || token.Type == TokenEmptyLines {
			chars = 0
		} else {
			chars = prev.EndChar()
		}

		offset = prev.End()
	}

	wg.Wait()

	if !hasFamilyName.Load() {
		checkFamilyName(list)
	}

	return
}

func mergeUnknown(list []*Token, token *Token, src string) []*Token {
	count := len(list)

	if count == 0 {
		return list
	}

	prevTokens, breakToken := getPrevTokens(list, -1, TokenName|TokenWord)
	prevCount := len(prevTokens)

	if prevCount == 0 {
		return list
	}

	if breakToken != nil && (breakToken.Type == TokenArrow || breakToken.Type == TokenEqual) {
		nextTokens := getNextTokens(prevTokens, 0, TokenWord)
		count := len(nextTokens)

		if count > 0 {
			if count == prevCount && prevTokens[prevCount-1].Type == TokenWord {
				count--
			}

			prevTokens = trimTokenLeft(prevTokens[count:], TokenSpace)
		}
	}

	return mergeTokens(list, prevTokens, token, src)
}

func mergeWords(list []*Token, token *Token, src string) []*Token {
	count := len(list)

	if count == 0 {
		return list
	}

	prevTokens, _ := getPrevTokens(list, -1, TokenWord)

	return mergeTokens(list, prevTokens, token, src)
}

func mergeTokens(list []*Token, prevTokens []*Token, token *Token, src string) []*Token {
	count := len(prevTokens)

	if count == 0 {
		return list
	}

	first := prevTokens[0]

	token.Length = token.End() - first.Offest
	token.Offest = first.Offest
	token.Line = first.Line
	token.Char = first.Char
	token.Text = src[token.Offest:token.End()]
	token.CharsNum = utf8.RuneCountInString(token.Text)

	return list[:len(list)-count]
}

func checkFamilyName(list []*Token) bool {
	tokens, breakToken := getPrevTokens(list, -1, TokenComment|TokenInvalid|TokenNewLine)

	if breakToken == nil || breakToken.Type == TokenEmptyLines {
		return false
	}

	total := len(tokens)

	tokens, breakToken = getPrevTokens(list, -total-1, TokenName|TokenSurname|TokenBracket|TokenPunctuation|TokenComment|TokenInvalid)

	if breakToken != nil && breakToken.Type == TokenNewLine {
		total += len(tokens)
		_, breakToken = getPrevTokens(list, -total-1, TokenComment|TokenNewLine)
	}

	if breakToken != nil && breakToken.Type != TokenEmptyLines {
		return false
	}

	for _, token := range tokens {
		if token.Type == TokenPunctuation && token.SubType != TokenComma {
			return false
		}
	}

	for _, token := range tokens {
		if token.Type != TokenName {
			continue
		}

		token.Type = TokenSurname
	}

	return true
}

func checkSurname(list []*Token, token *Token) {
	tokens, breakToken := getPrevTokens(list, -1, TokenName|TokenSurname)

	for _, token := range tokens {
		if token.Type == TokenSurname {
			token.Type = TokenName
		}
	}

	if len(tokens) > 0 {
		token.Type = TokenSurname
		return
	}

	if breakToken == nil || breakToken.SubType != TokenBracketRight {
		return
	}

	tokens, _ = getPrevTokens(cutAliasesRight(list), -1, TokenName)

	if len(tokens) > 0 {
		token.Type = TokenSurname
	}
}

func getPrevTokens(list []*Token, start int, validTokens TokenType) ([]*Token, *Token) {
	count := len(list)

	if count == 0 {
		return list, nil
	}

	if start < 0 {
		start = count + start // count + (-1)
	}

	var index int
	var breakToken *Token

	for index = start; index >= 0; index-- {
		t := list[index].Type

		if t == TokenSpace {
			continue
		}

		if validTokens&t == 0 {
			breakToken = list[index]
			index++
			break
		}
	}

	if index < 0 {
		index = 0
	}

	if index > start {
		return []*Token{}, breakToken
	}

	if list[index].Type == TokenSpace {
		index++
	}

	if index > start {
		return []*Token{}, breakToken
	}

	return list[index : start+1], breakToken
}

func getNextTokens(list []*Token, start int, validTokens TokenType) []*Token {
	count := len(list)

	if count == 0 {
		return list
	}

	var index int

	for index = start; index < count; index++ {
		t := list[index].Type

		if t == TokenSpace {
			continue
		}

		if validTokens&t == 0 {
			index--
			break
		}
	}

	if index >= count {
		index = count - 1
	}

	if index < start {
		return []*Token{}
	}

	if list[index].Type == TokenSpace {
		index--
	}

	if index < start {
		return []*Token{}
	}

	return list[start : index+1]
}

func trimTokenLeft(list []*Token, t TokenType) []*Token {
	count := len(list)

	if count > 0 && list[0].Type == t {
		return list[1:]
	}

	return list
}

func cutAliasesRight(list []*Token) []*Token {
	right := false

	for i := len(list) - 1; i >= 0; i-- {
		token := list[i]
		t := token.Type

		if t == TokenSpace {
			continue
		} else if t == TokenNewLine || t == TokenEmptyLines {
			return list
		}

		if !right {
			if token.SubType != TokenBracketRight {
				return list
			}

			right = true
			continue
		}

		if token.SubType == TokenBracketLeft {
			return list[:i]
		}
	}

	return list
}
