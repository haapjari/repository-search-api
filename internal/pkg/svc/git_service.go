package svc

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/logger"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/utils"
	"github.com/hhatto/gocloc"
)

type RepoOptions struct {
	URL       string
	Directory string
	Name      string
	Token     string
	Logger    logger.Logger
}

type Repo struct {
	opt   *RepoOptions
	tasks chan string
}

func NewRepo(opt *RepoOptions) *Repo {
	return &Repo{
		opt:   opt,
		tasks: make(chan string),
	}
}

func (r *Repo) Clone() error {
	auth := &http.BasicAuth{
		Username: "username",
		Password: r.opt.Token,
	}

	_, err := git.PlainClone(r.opt.Directory, false, &git.CloneOptions{
		URL:      r.opt.URL,
		Progress: os.Stdout,
		Auth:     auth,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Delete() error {
	if err := os.RemoveAll(r.opt.Directory); err != nil {
		return err
	}

	return nil
}

func (r *Repo) SelfWrittenLOC() (int, error) {
	return r.LinesOfCode(r.opt.Directory)
}

func (r *Repo) LinesOfCode(dir string) (int, error) {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	paths := []string{
		dir,
	}

	processor := gocloc.NewProcessor(languages, options)

	result, err := processor.Analyze(paths)
	if err != nil {
		return -1, nil
	}

	selfWrittenGoCode, ok := result.Languages["Go"]
	if ok {
		return int(selfWrittenGoCode.Code), nil
	}

	return 0, fmt.Errorf("projects have no go code")
}

func (r *Repo) CalculateLibraryLOC() (int, error) {
	// TODO
	modFilePath, err := utils.FindFile(r.opt.Directory, "go.mod")
	if err != nil {
		return -1, err
	}

	modFile := utils.NewModFile(modFilePath)
	err = modFile.Parse()
	if err != nil {
		return -1, err
	}

	// libraries are in format of golang.org/x/mod v0.15.0

	// libraryCount := 0

	goPath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		r.opt.Logger.Errorf("unable to exec command: %s", err.Error())
		return 0, err
	}

	for _, lib := range modFile.Replace {
		l := strings.Split(lib, " ")
		name := l[0]
		version := l[1]

		before, _, ok := strings.Cut(string(goPath), "\n")
		if ok {
			libraryPath := before + "/" + "pkg" + "/" + "mod" + "/" + name + "@" + version

			// r.log.Debugf("go get -v %s@%s", name, version)
			r.opt.Logger.Debugf("Library Path: %s", libraryPath)
		}
	}

	for _, lib := range modFile.Require {
		l := strings.Split(lib, " ")
		name := l[0]
		version := l[1]

		before, _, ok := strings.Cut(string(goPath), "\n")
		if ok {
			libraryPath := before + "/" + "pkg" + "/" + "mod" + "/" + name + "@" + version

			// r.log.Debugf("go get -v %s@%s", name, version)
			r.opt.Logger.Debugf("Library Path | %s", libraryPath)
		}
	}

	// mkdir tmp
	// cd tmp && go mod init temp
	//

	// $(go env GOPATH)/pkg/mod

	// TODO
	// parse libraries from go.mod file
	// -> clone
	// -> calculate
	// -> append
	// -> delete

	return -1, nil
}

// func (r *Repo) CalculateService() {
// 	go func() {
// 		for lib := range r.tasks {
// 			fmt.Println(lib)
// 		}
// 	}()
// }
