package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func mustCreateTempDir() string {
	dir, err := ioutil.TempDir("", "unit-testing")
	if err != nil {
		panic("Unable to create temp directory for testing: " + err.Error())
	}
	return dir
}

func createFiles(paths ...string) (string, []string) {
	result := make([]string, len(paths))
	rootDir := mustCreateTempDir()

	for i, path := range paths {
		dummyArtifact := rootDir + "/" + path
		mustMkdirAll(dummyArtifact)
		result[i] = dummyArtifact
	}

	return rootDir, result
}

var dateIncrement = 0

func mustMkdirAll(path string) {
	e := os.MkdirAll(path, 0755)
	if e != nil {
		panic("Unable to create dummy artifact directory for testing: " + e.Error())
	}
}

func touchFile(path string) {
	dateIncrement++
	t := time.Now().AddDate(0, 0, 1)
	e := os.Chtimes(path, t, t)
	if e != nil {
		panic("Unable to touch (increment date on file): " + e.Error())
	}
}

func TestFindLatestSnapshot(t *testing.T) {
	//example: /Users/sfaber/local-repo/com/linkedin/ligradle-core/core/4.0.67-SNAPSHOT
	//rootDir := "/Users/sfaber/local-repo/com/linkedin"
	rootDir, artifacts := createFiles(
		"mp1/moduleA/0.0.0-SNAPSHOT/0.zip",
		"mp1/moduleA/1.0.0-SNAPSHOT/1.zip",
		"mp1/moduleB/2.0.0-SNAPSHOT/2.zip",
		"mp2/moduleC/3.0.0-SNAPSHOT/3.zip",
		"mp2/moduleC/4.0.0-SNAPSHOT/4.zip",
		"mp2/moduleD/5.0.0-SNAPSHOT/5.zip",
		"mp3/moduleE/5.0.0-SNAPSHOT",
		"mp4/moduleE",
		"mp5",
	)

	//expect
	touchFile(artifacts[0])

	assertTrue(snapshot{"mp1", "0.0.0-SNAPSHOT"}.Equal(*findSnapshot(rootDir)), t)

	touchFile(artifacts[2])
	assertTrue(snapshot{"mp1", "2.0.0-SNAPSHOT"}.Equal(*findSnapshot(rootDir)), t)

	touchFile(artifacts[4])
	assertTrue(snapshot{"mp2", "4.0.0-SNAPSHOT"}.Equal(*findSnapshot(rootDir)), t)

	assertTrue(findSnapshot("some dir that does not exist") == nil, t)
}

func assertTrue(predicate bool, t *testing.T) {
	if !predicate {
		t.Errorf("Expected true!")
	}
}
