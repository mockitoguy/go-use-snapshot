package main

import (
	"os"
	"testing"
)

func createArtifactInLocalRepo(mp string, version string, t *testing.T) string {
	s := findSnapshotInLocalRepo()
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

func assertSafeVersion(version string, t *testing.T) string {
	s := findSnapshotInLocalRepo()
	if s != nil {
		assertTrue(s.version != version, t)
	}
	return version
}

func TestRun(t *testing.T) {
	mpDir := createArtifactInLocalRepo("dummy-mp-for-testing", "155.0.0-SNAPSHOT", t)
	defer os.RemoveAll(mpDir)

	dir := mustCreateTempDir()
	specFile := dir + "/product-spec.json"
	mustWriteFile(specFile, mustReadFile("testdata/product-only.json"))

	//when
	result := run(specFile)

	//then
	assertEquals("product-spec.json was updated! Stay happy!", result, t)

	//and
	expected := mustReadFile("testdata/product-only_updated.json")
	actual := mustReadFile(specFile)
	assertEquals(expected, actual, t)
}
