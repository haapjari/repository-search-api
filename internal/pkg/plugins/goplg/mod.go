package goplg

// type GoMod struct {
// 	Module  string
// 	Require []string
// 	Replace []string
// }
//
// func NewGoMod() *GoMod {
// 	return &GoMod{
// 		Module:  "",
// 		Require: []string{},
// 		Replace: []string{},
// 	}
// }
//
// // Give a path to "go.mod" - file, and return out a "GoMod" -structure including the
// // parsed data from the "go.mod" file.
// func parseGoMod(path string) (*GoMod, error) {
// 	goModFile := NewGoMod()
//
// 	pathToFile := filepath.Join(utils.GetProcessDirPath(), path)
//
// 	bytes, err := ioutil.ReadFile(pathToFile)
// 	if err != nil {
// 		fmt.Println("error, while reading the modfile: ", err)
// 		return goModFile, err
// 	}
//
// 	file, err := modfile.Parse(utils.GetProcessDirPath()+"/"+path, bytes, nil)
// 	if err != nil {
// 		fmt.Println("error, while parsing modfile: ", err)
// 		return goModFile, err
// 	}
//
// 	requirementsSlice := make([]string, len(file.Require))
// 	replacementsSlice := make([]string, len(file.Replace))
//
// 	for i := 0; i < len(file.Require); i++ {
// 		requirementsSlice[i] = file.Require[i].Mod.Path + " " + file.Require[i].Mod.Version
// 	}
//
// 	for i := 0; i < len(file.Replace); i++ {
// 		replacementsSlice[i] = file.Replace[i].New.Path
// 	}
//
// 	goModFile.Require = requirementsSlice
// 	goModFile.Replace = replacementsSlice
//
// 	return goModFile, nil
// }
