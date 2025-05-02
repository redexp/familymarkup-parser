package parser

type Loc struct {
	Start Position
	End   Position
}

type Position struct {
	Line int
	Char int
}

type OverlapType int

const (
	OverlapBefore OverlapType = iota - 2
	OverlapAfter
	OverlapByStart
	OverlapByEnd
	OverlapInner
	OverlapOuter
)

type PosCompare int

const (
	PosLt PosCompare = iota - 1
	PosEq
	PosGt
)

func (loc *Loc) OverlapType(other Loc) OverlapType {
	if loc.End.Compare(other.Start) <= PosEq {
		return OverlapBefore
	}

	if other.End.Compare(loc.Start) <= PosEq {
		return OverlapAfter
	}

	if loc.Start.Compare(other.Start) <= PosEq && loc.End.Compare(other.End) >= PosEq {
		return OverlapOuter
	}

	if loc.Start.Compare(other.Start) >= PosEq && loc.End.Compare(other.End) <= PosEq {
		return OverlapInner
	}

	if loc.Start.Compare(other.Start) == PosLt {
		return OverlapByEnd
	} else {
		return OverlapByStart
	}
}

func (loc *Loc) Overlaps(other Loc) bool {
	return loc.OverlapType(other) >= OverlapByStart
}

func (pos *Position) Compare(other Position) PosCompare {
	if pos.Line < other.Line {
		return PosLt
	}

	if pos.Line > other.Line {
		return PosGt
	}

	if pos.Char < other.Char {
		return PosLt
	}

	if pos.Char > other.Char {
		return PosGt
	}

	return PosEq
}
