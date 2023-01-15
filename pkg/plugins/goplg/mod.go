package goplg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/haapjari/glass/pkg/utils"
	"golang.org/x/mod/modfile"
)

type GoMod struct {
	Module  string
	Require []string
	Replace []string
}

func parseGoMod(path string) (*GoMod, error) {
	var goModFile GoMod

	pathToFile := filepath.Join(utils.GetProcessDirPath(), path)

	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		fmt.Println("error, while reading the modfile: ", err)
		return &goModFile, err
	}

	file, err := modfile.Parse(utils.GetProcessDirPath()+"/"+path, bytes, nil)
	if err != nil {
		fmt.Println("error, while parsing modfile: ", err)
		return &goModFile, err
	}

	requirementsSlice := make([]string, len(file.Require))
	replacementsSlice := make([]string, len(file.Replace))

	for i := 0; i < len(file.Require); i++ {
		requirementsSlice[i] = file.Require[i].Mod.Path + " " + file.Require[i].Mod.Version
	}

	for i := 0; i < len(file.Replace); i++ {
		replacementsSlice[i] = file.Replace[i].New.Path
	}

	goModFile.Require = requirementsSlice
	goModFile.Replace = replacementsSlice

	return &goModFile, nil
}
