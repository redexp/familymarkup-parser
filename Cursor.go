package parser

import (
	"iter"
	"slices"
)

type Cursor struct {
	Tokens []*Token
	Count  int
	Index  int
}

func NewCursor(tokens []*Token) *Cursor {
	count := len(tokens)

	return &Cursor{
		Tokens: tokens,
		Count:  count,
		Index:  0,
	}
}

func (c *Cursor) Cur() *Token {
	return c.Tokens[c.Index]
}

func (c *Cursor) Iter() iter.Seq[*Token] {
	return func(yield func(*Token) bool) {
		for _ = c.Index; c.Index < c.Count; c.Index++ {
			if !yield(c.Cur()) {
				break
			}
		}
	}
}

func (c *Cursor) GetAllNext(validTokens []TokenType) (tokens []*Token) {
	defer c.StepBackIfNotEnd()

	for c.Index++; c.Index < c.Count; c.Index++ {
		token := c.Tokens[c.Index]

		if !slices.Contains(validTokens, token.Type) && !slices.Contains(validTokens, token.SubType) {
			return
		}

		tokens = append(tokens, token)
	}

	return
}

func (c *Cursor) PickClosestTextToken(dir int) *Token {
	index := c.Index

	if index >= c.Count && dir < 0 {
		index = c.Count - 1
	}

	if dir < -1 {
		dir = -1
		index--
	}

	for i := index; i >= 0 && i < c.Count; i += dir {
		token := c.Tokens[i]

		switch token.Type {
		case TokenSpace, TokenNewLine, TokenEmptyLines, TokenInvalid:
			continue
		default:
			return token
		}
	}

	return nil
}

func (c *Cursor) PickNext() *Token {
	return c.PickClosestTextToken(+1)
}

func (c *Cursor) PickPrev() *Token {
	return c.PickClosestTextToken(-1)
}

func (c *Cursor) PickPrevSkipCur() *Token {
	return c.PickClosestTextToken(-2)
}

func (c *Cursor) IsNext(t TokenType) bool {
	for i := c.Index + 1; i < c.Count; i++ {
		token := c.Tokens[i]

		if token.Type == TokenSpace {
			continue
		}

		return token.Type == t
	}

	return false
}

func (c *Cursor) IsStartOfNewLine() bool {
	for i := c.Index - 1; i >= 0; i-- {
		t := c.Tokens[i].Type

		if t == TokenSpace {
			continue
		}

		return t == TokenNewLine || t == TokenEmptyLines
	}

	return true
}

func (c *Cursor) StepBackIfNotEnd() {
	if c.Index > 0 && c.Index < c.Count {
		c.Index--
	}
}
