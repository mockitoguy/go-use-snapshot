package main

import (
	"fmt"
	"io/ioutil"
)

const resultUpdated = "product-spec.json was updated! Stay happy!"
const resultNoSnapshot = "Unable to find any snapshot inside of ~/local-repo. Please inspect this directory."
const resultNoProductSpec = "Unable to read 'product-spec.json' file in current directory. Does it exist?"

func run(productSpecPath string, localRepoMpDir string) string {
	snapshot := findSnapshot(localRepoMpDir)
	if snapshot == nil {
		return resultNoSnapshot
	}
	fmt.Printf("Found most recent snapshot of '%v' at version '%v'\n", snapshot.mp, snapshot.version)

	productSpec, e := ioutil.ReadFile(productSpecPath)
	if e != nil {
		return resultNoProductSpec
	}

	content := string(productSpec)

	updated := updateVersion(&content, snapshot.mp, snapshot.version)

	if !updated {
		return fmt.Sprintf("product-spec.json was not updated! Does it declare dependency on '%v'?", snapshot.mp)
	}

	mustWriteFile(productSpecPath, content)

	return resultUpdated
}

func mustWriteFile(filePath string, content string) {
	e := ioutil.WriteFile(filePath, []byte(content), 0644) //TODO: don't change the file mode
	if e != nil {
		panic(fmt.Sprintf("Unable to update file: '%v' due to: %v\n", filePath, e.Error()))
	}
}

func main() {
	message := run("product-spec.json", getLocalRepoMpPath())
	fmt.Print(message)
}
