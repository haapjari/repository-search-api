package util

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hhatto/gocloc"
	"golang.org/x/mod/modfile"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func FindFile(path, fileName string) (string, error) {
	var foundPath string
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == fileName {
			foundPath = currentPath
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", os.ErrNotExist
	}

	return foundPath, nil
}

func OnlyLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func Clone(token string, url string) (string, error) {
	dir, err := os.MkdirTemp("", "clone-")
	if err != nil {
		return "", Error(fmt.Errorf("unable to create a temporary directory: %v", err))
	}

	auth := &http.BasicAuth{
		Username: "x-access-token",
		Password: strings.Split(token, " ")[1],
	}

	repo, err := git.PlainClone(dir, false,
		&git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     auth,
		})
	if err != nil {
		return "", Error(fmt.Errorf("unable to clone the repository: %v", err))
	}

	if repo != nil {
		return dir, nil
	}

	return "", Error(fmt.Errorf("unable to clone the repository"))
}

// ParseModFile reads the go.mod file and returns the libraries that are required by the project.
func ParseModFile(path string) ([]string, error) {
	data, err := os.ReadFile(path + "/" + "go.mod")
	if err != nil {
		return nil, Error(err)
	}

	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}

	libraries := make([]string, len(file.Require))

	for _, r := range file.Require {
		libraries = append(libraries, r.Mod.Path+"@"+r.Mod.Version)
	}

	return libraries, nil
}

// FetchLibrary creates a temporary module to fetch a library without affecting the main go.mod.
func FetchLibrary(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "temp-mod")
	if err != nil {
		return "", Error(fmt.Errorf("unable to create a temporary directory: %v", err))
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	if err = os.Chdir(tempDir); err != nil {
		return "", Error(fmt.Errorf("unable to change the working directory: %v", err))
	}

	if out, execErr := exec.Command("go", "mod", "init", "temp").CombinedOutput(); execErr != nil {
		return "", Error(fmt.Errorf("unable to initialize the temporary module: %v", string(out)))
	}

	if out, execErr := exec.Command("go", "get", url).CombinedOutput(); execErr != nil {
		return "", Error(fmt.Errorf("unable to fetch the library: %v", string(out)))
	}

	p, err := exec.Command("go", "env", "GOMODCACHE").Output()
	if err != nil {
		return "", nil
	}

	return filepath.Join(strings.TrimSpace(string(p)), url), nil
}

// CalcLOC is a method of the RepositoryService struct. It calculates the lines of code
// of a directory based on the provided language.
func CalcLOC(dir string, lang string) (int, error) {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	paths := []string{
		dir,
	}

	processor := gocloc.NewProcessor(languages, options)

	result, err := processor.Analyze(paths)
	if err != nil {
		return -1, Error(fmt.Errorf("unable to analyze the repository: %v", err))
	}

	report, ok := result.Languages[lang]
	if ok {
		return int(report.Code), nil
	}

	if report, ok = result.Languages[strings.ToLower(lang)]; ok {
		return int(report.Code), nil
	}

	return 0, fmt.Errorf("language %s not found in analysis results", lang)
}
