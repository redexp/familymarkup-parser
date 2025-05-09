package parser

import "iter"

type Root struct {
	Loc
	Families []*Family
	Comments []*Token
}

type Family struct {
	Loc
	Name      *Token
	Aliases   []*Token
	Relations []*Relation
	Comments  []*Token
}

type Relation struct {
	Loc
	Sources     *RelList
	Arrow       *Token
	IsFamilyDef bool
	Label       *Token
	Comments    []*Token
	Targets     *RelList

	Family *Family
}

type RelList struct {
	Loc
	Persons    []*Person
	Separators []*Token
}

type Person struct {
	Loc
	Unknown  *Token
	Num      *Token
	Name     *Token
	Aliases  []*Token
	Surname  *Token
	Comments []*Token

	Side    SideType
	Index   int
	IsChild bool

	Relation *Relation
}

type SideType int

const (
	SideSources SideType = iota
	SideTargets
)

func (rel *Relation) PersonsIter() iter.Seq[*Person] {
	return func(yield func(*Person) bool) {
		for _, list := range []*RelList{rel.Sources, rel.Targets} {
			if list == nil {
				continue
			}

			for _, person := range list.Persons {
				if !yield(person) {
					return
				}
			}
		}
	}
}
