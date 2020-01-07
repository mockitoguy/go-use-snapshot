package main

import (
	"os"
	"testing"
)

//creates dummy artifact in local repository
// and pushes it's modification date few days ahead
func createArtifactInLocalRepo(mp string, version string, t *testing.T) string {
	s := findSnapshot(getLocalRepoMpPath())
	if s != nil {
		assertTrue(s.version != version, t)
	}
	dir := getLocalRepoMpPath()
	mpDir := dir + "/" + mp
	dummyArtifact := mpDir + "/some-module/" + version + "/some-artifact"
	mustMkdirAll(dummyArtifact)
	touchFile(dummyArtifact)
	return mpDir
}

//creates tmp file with content from supplied test file path
//returns path to the newly created product spec file
func createProductSpec(testFilePath string) string {
	dir := mustCreateTempDir()
	specFile := dir + "/product-spec.json"
	mustWriteFile(specFile, mustReadFile(testFilePath))
	return specFile
}

func TestE2E(t *testing.T) {
	mpDir := createArtifactInLocalRepo("dummy-mp-for-testing", "155.0.0-SNAPSHOT", t)
	defer os.RemoveAll(mpDir)

	specFile := createProductSpec("testdata/product-only.json")

	//when
	result := run(specFile, getLocalRepoMpPath())

	//then
	assertEquals(resultUpdated, result, t)

	//and
	expected := mustReadFile("testdata/product-only_updated.json")
	actual := mustReadFile(specFile)
	assertEquals(expected, actual, t)
}

func TestNoSnapshotFound(t *testing.T) {
	//when
	result := run("product-spec.json", "not existing dir")

	//then
	assertEquals(resultNoSnapshot, result, t)
}
