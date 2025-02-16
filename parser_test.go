package parser

import (
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	types "github.com/onsi/gomega/types"
)

type M = types.GomegaMatcher

func TestParser(t *testing.T) {
	g := NewWithT(t)

	text := testFile("main.fml")

	root := Parser(text)

	g.Expect(root).To(testPoint(Fields{
		"Loc": testLoc(2, 0, 23, 18),
		"Families": testArr(
			testPoint(Fields{
				"Loc":     testLoc(5, 0, 18, 18),
				"Name":    testToken("Family"),
				"Aliases": testTokens("Alias", "Alias2"),
				"Comments": testTokens(
					"* Family comment",
					"* Family comment 2",
					"* Family comment 3",
				),
				"Relations": testArr(
					testPoint(Fields{
						"Loc": testLoc(10, 0, 10, 12),
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name", "Name2"),
							"Separators": testTokens("+"),
						}),
					}),
					testPoint(Fields{
						"Loc": testLoc(12, 2, 12, 42),
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name3", "Name4"),
							"Separators": testTokens("+"),
						}),
						"Arrow":    testToken("="),
						"Label":    testToken("label"),
						"Comments": testTokens("# relation comment"),
					}),
					testPoint(Fields{
						"Loc": testLoc(14, 0, 16, 5),
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name5", "mother?"),
							"Separators": testTokens("+"),
						}),
						"Arrow": testToken("="),
						"Targets": testPoint(Fields{
							"Persons": testArr(
								testPoint(Fields{
									"Name":     testToken("Name6"),
									"Aliases":  testTokens("NameAlias"),
									"Surname":  testToken("Surname"),
									"Comments": testTokens("# person comment"),
								}),
								testPoint(Fields{
									"Name": testToken("Name7"),
								}),
							),
						}),
					}),
				),
			}),
			testPoint(Fields{
				"Loc":      testLoc(20, 2, 23, 18),
				"Name":     testToken("Family2"),
				"Comments": testTokens("* Family2 comment"),
				"Relations": testArr(
					testPoint(Fields{
						"Sources": testPoint(Fields{
							"Persons": testArr(
								testPoint(Fields{
									"Loc":     testLoc(22, 0, 22, 8),
									"Unknown": testToken("unknown?"),
								}),
								testPoint(Fields{
									"Loc":  testLoc(22, 11, 22, 16),
									"Name": testToken("Name1"),
								}),
							),
							"Separators": testTokens("+"),
						}),
						"Arrow":    BeNil(),
						"Label":    BeNil(),
						"Targets":  BeNil(),
						"Comments": testTokens("* relation comment"),
					}),
				),
			}),
		),
	}))
}

func TestNameless(t *testing.T) {
	g := NewWithT(t)

	text := testFile("nameless.family")

	root := Parser(text)

	g.Expect(root).To(testPoint(Fields{
		"Loc": testLoc(0, 0, 2, 13),
		"Families": testArr(
			testPoint(Fields{
				"Loc":  testLoc(0, 0, 2, 13),
				"Name": BeNil(),
				"Relations": testArr(
					testPoint(Fields{
						"Loc": testLoc(0, 0, 0, 12),
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name", "Name2"),
							"Separators": testTokens("+"),
						}),
					}),
					testPoint(Fields{
						"Loc": testLoc(2, 0, 2, 13),
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name3", "Name4"),
							"Separators": testTokens("+"),
						}),
					}),
				),
			}),
		),
	}))

	root = Parser("Nameless (Alias, Alias2)")

	g.Expect(root).To(testPoint(Fields{
		"Loc": testLoc(0, 0, 0, 24),
		"Families": testArr(
			testPoint(Fields{
				"Loc":     testLoc(0, 0, 0, 24),
				"Name":    testToken("Nameless"),
				"Aliases": testTokens("Alias", "Alias2"),
			}),
		),
	}))
}

func testLoc(startLine, startChar, endLine, endChar int) M {
	return testProps(Fields{
		"Start": testProps(Fields{
			"Line": BeNumerically("==", startLine),
			"Char": BeNumerically("==", startChar),
		}),
		"End": testProps(Fields{
			"Line": BeNumerically("==", endLine),
			"Char": BeNumerically("==", endChar),
		}),
	})
}

func testProps(props Fields) M {
	return MatchFields(IgnoreExtras, props)
}

func testPoint(p Fields) M {
	return PointTo(testProps(p))
}

func testArr(arr ...M) M {
	e := Elements{}

	for i, item := range arr {
		e[strconv.Itoa(i)] = item
	}

	return MatchAllElementsWithIndex(IndexIdentity, e)
}

func testToken(text string) M {
	return testPoint(Fields{
		"Text": BeIdenticalTo(text),
	})
}

func testTokens(texts ...string) M {
	tokens := make([]M, len(texts))

	for i, text := range texts {
		tokens[i] = testToken(text)
	}

	return testArr(tokens...)
}

func testPersons(names ...string) M {
	persons := make([]M, len(names))

	for i, name := range names {
		p := Fields{}

		if strings.HasSuffix(name, "?") {
			p["Unknown"] = testToken(name)
		} else {
			p["Name"] = testToken(name)
		}

		persons[i] = testPoint(p)
	}

	return testArr(persons...)
}
