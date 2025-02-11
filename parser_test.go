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
		"Families": testArr(
			testPoint(Fields{
				"Name":    testToken("Family"),
				"Aliases": testTokens("Alias", "Alias2"),
				"Comments": testTokens(
					"* Family comment",
					"* Family comment 2",
					"* Family comment 3",
				),
				"Relations": testArr(
					testPoint(Fields{
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name", "Name2"),
							"Separators": testTokens("+"),
						}),
					}),
					testPoint(Fields{
						"Sources": testPoint(Fields{
							"Persons":    testPersons("Name3", "Name4"),
							"Separators": testTokens("+"),
						}),
						"Arrow":    testToken("="),
						"Label":    testToken("label"),
						"Comments": testTokens("# relation comment"),
					}),
					testPoint(Fields{
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
				"Name":     testToken("Family2"),
				"Comments": testTokens("* Family2 comment"),
				"Relations": testArr(
					testPoint(Fields{
						"Sources": testPoint(Fields{
							"Persons": testPersons("unknown?", "Name1"),
						}),
						"Comments": testTokens("* relation comment"),
					}),
				),
			}),
		),
	}))
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
