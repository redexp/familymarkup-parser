package parser

func Parse(src string) *Root {
	return ParseTokens(Lexer(src))
}

func ParseTokens(tokens []*Token) *Root {
	return visitRoot(NewCursor(tokens))
}

func visitRoot(c *Cursor) *Root {
	root := &Root{}
	start := c.PickNext()

	if start == nil {
		return root
	}

	root.Start = toPos(start)

	for token := range c.Iter() {
		switch token.Type {
		case TokenComment:
			root.Comments = append(root.Comments, token)

		case TokenNewLine, TokenEmptyLines, TokenSpace, TokenInvalid:
			continue

		default:
			f := visitFamily(c)
			root.Families = append(root.Families, f)
			c.StepBackIfNotEnd()
		}
	}

	root.End = toEndPos(c.PickPrev())

	return root
}

func visitFamily(c *Cursor) (family *Family) {
	family = &Family{}
	family.Start = toPos(c.PickNext())

	defer func() {
		for _, rel := range family.Relations {
			rel.Family = family
		}
	}()

	for token := range c.Iter() {
		switch token.Type {
		case TokenSurname:
			if family.Name != nil {
				family.End = toEndPos(c.PickPrevSkipCur())
				return
			}

			family.Name = token

			tokens := c.GetAllNext(TokenSurname | TokenBracket | TokenPunctuation | TokenSpace | TokenInvalid)

			for _, token := range tokens {
				if token.SubType == TokenAlias {
					family.Aliases = append(family.Aliases, token)
					continue
				}

				if token.Type == TokenSurname {
					token.ErrType = ErrUnexpected
					continue
				}
			}

		case TokenComment:
			family.Comments = append(family.Comments, token)

		case TokenNewLine, TokenEmptyLines, TokenSpace, TokenInvalid:
			continue

		case TokenName, TokenUnknown, TokenNum:
			rel := visitRelation(c)
			family.Relations = append(family.Relations, rel)
			c.StepBackIfNotEnd()

		default:
			token.ErrType = ErrUnexpected
		}
	}

	family.End = toEndPos(c.PickPrev())

	return
}

func visitRelation(c *Cursor) (rel *Relation) {
	rel = &Relation{
		Sources: &RelList{},
	}

	rel.Start = toPos(c.PickNext())

	defer func() {
		if rel.Targets != nil && len(rel.Targets.Persons) == 0 && len(rel.Targets.Separators) == 0 {
			rel.Targets = nil
		}

		for side, list := range []*RelList{rel.Sources, rel.Targets} {
			if list == nil {
				continue
			}

			persons := list.Persons
			count := len(persons)

			if count == 0 {
				continue
			}

			first := persons[0]
			last := persons[count-1]

			list.Start = first.Start
			list.End = last.End

			for i, person := range persons {
				person.Side = SideType(side)
				person.Index = i
				person.Relation = rel
			}
		}

		rel.IsFamilyDef = rel.Arrow != nil && rel.Arrow.SubType == TokenEqual

		if rel.Targets != nil && rel.IsFamilyDef {
			for _, person := range rel.Targets.Persons {
				person.IsChild = true
			}
		}
	}()

	list := rel.Sources

	for token := range c.Iter() {
		switch token.Type {
		case TokenName, TokenUnknown, TokenNum:
			p := visitPerson(c)
			list.Persons = append(list.Persons, p)
			c.StepBackIfNotEnd()

		case TokenWord, TokenPunctuation:
			list.Separators = append(list.Separators, token)

		case TokenArrow:
			if rel.Arrow != nil {
				token.ErrType = ErrUnexpected
				continue
			}

			rel.Arrow = token
			rel.Targets = &RelList{}
			list = rel.Targets

			tokens := c.GetAllNext(TokenSpace | TokenWord)

			for _, token := range tokens {
				if token.Type == TokenWord {
					rel.Label = token
				}
			}

		case TokenComment:
			rel.Comments = append(rel.Comments, token)

		case TokenEmptyLines:
			if list == rel.Targets && c.IsNext(TokenNum) {
				token.ErrType = ErrUnexpected
				continue
			}
			rel.End = toEndPos(c.PickPrev())
			return

		case TokenSpace, TokenInvalid, TokenNewLine:
			continue

		default:
			rel.End = toEndPos(c.PickPrevSkipCur())
			return
		}
	}

	rel.End = toEndPos(c.PickPrev())
	return
}

func visitPerson(c *Cursor) (p *Person) {
	p = &Person{}
	p.Start = toPos(c.PickNext())

	isStartOfLine := c.IsStartOfNewLine()

	for token := range c.Iter() {
		switch token.Type {
		case TokenUnknown:
			if p.Unknown == nil {
				p.Unknown = token
				continue
			}

			p.End = toEndPos(c.PickPrevSkipCur())
			return

		case TokenNum:
			p.Num = token

		case TokenName:
			if p.Name == nil {
				p.Name = token
				continue
			}

			token.ErrType = ErrUnexpected

		case TokenSurname:
			if p.Surname == nil {
				p.Surname = token
				continue
			}

			token.ErrType = ErrUnexpected

		case TokenBracket:
			if token.SubType == TokenBracketRight {
				continue
			}

			tokens := c.GetAllNext(TokenAlias | TokenComma | TokenSpace)

			for _, token := range tokens {
				if token.SubType == TokenAlias {
					p.Aliases = append(p.Aliases, token)
				}
			}

		case TokenComment:
			p.Comments = append(p.Comments, token)

		case TokenNewLine:
			if !isStartOfLine {
				p.End = toEndPos(c.PickPrevSkipCur())
				return
			}

			tokens := c.GetAllNext(TokenComment | TokenNewLine | TokenSpace)

			p.End = toEndPos(c.PickPrev())

			if len(tokens) > 0 {
				c.Index++
			}

			for _, token := range tokens {
				if token.Type == TokenComment {
					p.Comments = append(p.Comments, token)
				}
			}

			return

		case TokenSpace, TokenInvalid:
			continue

		default:
			p.End = toEndPos(c.PickPrevSkipCur())
			return
		}
	}

	p.End = toEndPos(c.PickPrev())

	return
}

func toPos(token *Token) Position {
	return Position{
		Line: token.Line,
		Char: token.Char,
	}
}

func toEndPos(token *Token) Position {
	return Position{
		Line: token.Line,
		Char: token.EndChar(),
	}
}
