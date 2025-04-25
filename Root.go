package parser

type Root struct {
	Loc
	Families []*Family
	Comments []*Token
}

type Loc struct {
	Start Position
	End   Position
}

type Position struct {
	Line int
	Char int
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
	Sources  *RelList
	Arrow    *Token
	Label    *Token
	Comments []*Token
	Targets  *RelList
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
	IsChild  bool
}
