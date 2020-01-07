package main

import (
	"io/ioutil"
	"os"
	"path"
)

func readDir(dir string) []os.FileInfo {
	info, e := os.Stat(dir)
	if e != nil || !info.IsDir() {
		return []os.FileInfo{}
	}

	children, err := ioutil.ReadDir(dir)
	if err == nil {
		return children
	} else {
		return []os.FileInfo{}
	}
}

type snapshot struct {
	mp      string
	version string
}

func (self snapshot) Equal(other snapshot) bool {
	return self.mp == other.mp && self.version == other.version
}

/*
Finds MP name + version that was most recently released to ~/local-repo

	Structure of the local-repo dir:
		com/linkedin 				//standard group name of MPs, the correct root dir for serching
			ligradle-core 			//MP name
				core				//module name
					1.0.0-SNAPSHOT  //version
						artifact1   //artifact
						artifact2   //another artifact
						...
					1.0.1-SNAPSHOT	//another version
					...
				utils				//another module
				...
			ligradle-jvm			//another mp
			...

- this method finds latest artifact by file modification time and returns the MP name and version.
- nil is returned if there are no artifacts
- rootDir must be "~/local-repo/com/linkedin" unless we are running tests
*/
func findSnapshot(rootDir string) *snapshot {
	mpDirs := readDir(rootDir)
	var latestMp, latestVersion, latestArtifact os.FileInfo
	for _, mpDir := range mpDirs {
		mpDirPath := rootDir + "/" + mpDir.Name()
		mpModules := readDir(mpDirPath)
		for _, mpModule := range mpModules {
			mpModulePath := mpDirPath + "/" + mpModule.Name()
			versions := readDir(mpModulePath)
			for _, version := range versions {
				//TODO: regex for snapshot
				//re := regexp.MustCompile(`\d+\.\d+\.\d+\.*-SNAPSHOT`)
				//re.MatchString(info.Name()
				versionPath := mpModulePath + "/" + version.Name()
				artifacts := readDir(versionPath)
				for _, artifact := range artifacts {
					if latestArtifact == nil || artifact.ModTime().After(latestArtifact.ModTime()) {
						latestArtifact = artifact
						latestVersion = version
						latestMp = mpDir
					}
				}
			}
		}
	}
	if latestMp != nil {
		return &snapshot{latestMp.Name(), latestVersion.Name()}
	} else {
		return nil
	}
}

func getLocalRepoMpPath() string {
	dir, e := os.UserHomeDir()
	if e != nil {
		panic("Unable to read home directory: " + e.Error())
	}

	rootDir := path.Join(dir, "local-repo", "com", "linkedin")
	return rootDir
}
