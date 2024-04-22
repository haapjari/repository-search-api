package svc

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/models"
	"log/slog"
	"os"
)

type GitHubService struct {
	Token           string
	QueryParameters *models.QueryParameters
}

func NewGitHubService(token string, params *models.QueryParameters) *GitHubService {
	return &GitHubService{
		Token:           token,
		QueryParameters: params,
	}
}

func (g *GitHubService) Query() []*models.Repository {
	// TODO

	// TODO: Fetch a List of URLs

	url := "https://github.com/go-git/go-git.git"

	path, err := g.clone(url)
	if err != nil {
		return nil
	}

	// TODO: Create Repository Object
	// TODO: Fetch the Data from Repository Object and Construct a Response

	fmt.Println("")
	fmt.Println(path)

	return nil
}

func (g *GitHubService) clone(url string) (string, error) {
	dir, err := os.MkdirTemp("", "clone-")
	if err != nil {
		slog.Error("unable to create directory: " + err.Error())
		return "", err
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     nil,
	})
	if err != nil {
		slog.Error("unable to clone repository: " + err.Error())
		return "", err
	}

	if repo != nil {
		return dir, nil
	}

	return "", fmt.Errorf("unable to clone the repository")
}

/*
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
*/
