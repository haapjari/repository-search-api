package goplg

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/haapjari/glsgen/pkg/models"
	"github.com/haapjari/glsgen/pkg/utils"
	"github.com/hhatto/gocloc"
	"github.com/pingcap/errors"
)

func (g *GoPlugin) getAllRepositories() []models.Repository {
	var repositories []models.Repository

	g.DatabaseClient.Find(&repositories)

	return repositories
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func isLocalPath(path string) bool {
	return strings.HasPrefix(path, "./")
}

func (g *GoPlugin) repositoryExists(repository models.Repository, repositories []models.Repository) bool {
	repositoryUrl := repository.RepositoryUrl

	for _, r := range repositories {
		if r.RepositoryUrl == repositoryUrl {
			return true
		}
	}

	return false
}

func (g *GoPlugin) hasBeenEnriched(repository models.Repository, repositories []models.Repository) bool {
	var queriedRepository models.Repository
	repositoryUrl := repository.RepositoryUrl

	for _, r := range repositories {
		if r.RepositoryUrl == repositoryUrl {
			queriedRepository = r
			break
		}
	}
	v := reflect.ValueOf(queriedRepository)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			return false
		}
	}
	return true
}

// Calculates the lines of code using https://github.com/hhatto/gocloc
// in the path provided and return the value.
func (g *GoPlugin) calculateCodeLines(path string) (int, error) {
	if !pathExists(path) {
		return 0, errors.New("Error - " + path + " doesn't exist.")
	}

	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	paths := []string{
		path,
	}

	processor := gocloc.NewProcessor(languages, options)

	result, err := processor.Analyze(paths)
	if err != nil {
		return 0, errors.New("Error - Failed to calculate lines of code.")
	}

	return int(result.Total.Code), nil
}

// Delete the contents of tmp -folder.
func (g *GoPlugin) pruneTemporaryFolder() {
	utils.Command("chmod", "-R", "777", utils.GetProcessDirPath())
	os.RemoveAll(utils.GetProcessDirPath())
	os.MkdirAll(utils.GetProcessDirPath(), os.ModePerm)
}

// Check if a folder exists in the file system.
func pathExists(folderPath string) bool {
	// Use os.Stat to get the file information for the folder
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			// The folder does not exist
			return false
		} else {
			// Some other error occurred
			fmt.Printf("Error checking if folder exists: %v", err)
			return false
		}
	}

	// The folder exists
	return true
}

// Parse "github.com/mholt/archiver/v3 v3.5.1" into the format "github.com/mholt/archiver/v3@v3.5.1"
func downloadableFormat(input string) string {
	// Split the input string on the first space character
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	// Return the first part (the package name) followed by an @ symbol and the second part (the version)
	return parts[0] + "@" + parts[1]
}

// Parses the input string to correct format.
// Input: "github.com/example/component v1.0.0"
// Output: "github.com/example/component@v1.0.0"
// Input: "github.com/example/Component v1.0.0"
// Output: "github.com/example/!component@v1.0.0"
func parseLibraryUrl(input string) string {
	// Split the input string on the first space character
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	// Split the package name on the '/' character
	packageNameParts := strings.Split(parts[0], "/")

	// Add the '!' prefix and lowercase each upper character, of the package name
	for i, part := range packageNameParts {
		modifiedPart := ""
		for _, subPart := range part {
			if unicode.IsUpper(subPart) {
				modifiedPart += "!" + strings.ToLower(string(subPart))
			} else {
				modifiedPart += string(subPart)
			}
		}
		packageNameParts[i] = modifiedPart
	}

	// Join the modified package name parts with '/' characters
	packageName := strings.Join(packageNameParts, "/")

	// Return the modified package name followed by an '@' symbol and the version
	return packageName + "@" + parts[1]
}
