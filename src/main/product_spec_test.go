package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func mustReadFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot read file '%v' due to: %v", path, err.Error()))
	}
	s := string(b)
	return s
}

func TestProductSpecUpdate(t *testing.T) {
	//given
	testFiles := []string{"toolchains-and-product", "toolchains-only", "product-only", "none"}

	for _, testFile := range testFiles {
		t.Run("Update version in "+testFile, func(t *testing.T) {
			//given
			initial := mustReadFile("testdata/" + testFile + ".json")
			updated := mustReadFile("testdata/" + testFile + "_updated.json")

			//when
			was_updated := updateVersion(&initial, "dummy-mp-for-testing", "155.0.0-SNAPSHOT")

			//then
			expected_updated := testFile != "none"
			assertEquals(expected_updated, was_updated, t)
			assertEquals(initial, updated, t)
		})
	}
}

func TestFindLineNo(t *testing.T) {
	content := []string{"1", "2", "  1", "  2", "3"}

	assertEquals(1, findLineIdx(content, "2"), t)
	assertEquals(3, findLineIdx(content, "  2"), t)

	assertEquals(1, findLineIdx(content, "1", "2"), t)
	assertEquals(3, findLineIdx(content, "  1", "  2"), t)

	assertEquals(4, findLineIdx(content, "  1", "  2", "3"), t)

	assertEquals(1, findLineIdx(content, "^1$", "^2$"), t)
	assertEquals(3, findLineIdx(content, " +1", " +2"), t)
}

func assertEquals(one interface{}, two interface{}, t *testing.T) {
	if one != two {
		t.Errorf("%v\n--- not equal to ---\n%v", one, two)
	}
}
