package goplg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"github.com/hhatto/gocloc"
	JSONParser "github.com/tidwall/gjson"
)

func (g *GoPlugin) updateCoreSize(name string, lines int) {
	// Copy the repository struct to a new variable.
	var repositoryStruct models.Repository

	// Find matching repository from the database.
	if err := g.DatabaseClient.Where("repository_name = ?", name).First(&repositoryStruct).Error; err != nil {
		utils.CheckErr(err)
	}

	// Update the OriginalCodebaseSize variable, with calculated value.
	repositoryStruct.OriginalCodebaseSize = strconv.Itoa(lines)

	// Update the database.
	g.DatabaseClient.Model(&repositoryStruct).Updates(repositoryStruct)
}

func (g *GoPlugin) restoreModuleFiles() {
	utils.RemoveFiles("go.mod", "go.sum")
	utils.CopyFile("go.mod.bak", "go.mod")
	utils.CopyFile("go.sum.bak", "go.sum")
	utils.RemoveFiles("go.mod.bak", "go.sum.bak")
}

func (g *GoPlugin) resetModuleFiles() {
	utils.RemoveFiles("go.mod", "go.sum")
	utils.CopyFile("go.mod.bak", "go.mod")
	utils.CopyFile("go.sum.bak", "go.sum")
}

func (g *GoPlugin) saveModuleFiles() {
	utils.CopyFile("go.mod", "go.mod.bak")
	utils.CopyFile("go.sum", "go.sum.bak")
}

func (g *GoPlugin) processSourceGraphResponse(length int, repositories []models.SourceGraphRepositories) {
	var wg sync.WaitGroup

	// Semaphore is a safeguard to goroutines, to allow only "MaxThreads" run at the same time.
	semaphore := make(chan int, g.MaxRoutines)

	for i := 0; i < length; i++ {
		semaphore <- 1
		wg.Add(1)

		go func(i int) {
			r := models.Repository{RepositoryName: repositories[i].Name, RepositoryUrl: repositories[i].Name, OpenIssueCount: "", ClosedIssueCount: "", OriginalCodebaseSize: "", LibraryCodebaseSize: "", RepositoryType: "", PrimaryLanguage: ""}
			g.DatabaseClient.Create(&r)

			defer func() { <-semaphore }()
		}(i)
		wg.Done()
	}
	wg.Wait()

	// When the Channel Length is not 0, there is still running goroutines.
	for !(len(semaphore) == 0) {
		continue
	}
}

// TODO: Refactor, to not to use HTTP requests (?)
func (g *GoPlugin) getAllRepositories() []models.Repository {
	var repositories []models.Repository

	g.DatabaseClient.Find(&repositories)

	return repositories
}

// Calculates the lines of code using https://github.com/hhatto/gocloc
// in the path provided and return the value.
func (g *GoPlugin) calculateCodeLines(path string) int {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	paths := []string{
		path,
	}

	processor := gocloc.NewProcessor(languages, options)

	result, err := processor.Analyze(paths)
	utils.CheckErr(err)

	return int(result.Total.Code)

}

// Delete the contents of tmp -folder.
func (g *GoPlugin) pruneTemporaryFolder() {
	utils.Command("chmod", "-R", "777", utils.GetProcessDirPath())
	os.RemoveAll(utils.GetProcessDirPath())
	os.MkdirAll(utils.GetProcessDirPath(), os.ModePerm)
}

// Parse repository name from url.
func parseName(s string) (string, string, error) {
	parts := strings.Split(s, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid repository string: %s", s)
	}
	return parts[1], parts[2], nil
}

// Creates a slice of repositories, which are duplicates in an original list.
func findDuplicates(repositories []models.Repository) []models.Repository {
	// Create a map to store the names of the repositories that we've seen so far
	seenRepositories := make(map[string]bool)

	// Create a slice to store the duplicate entries
	duplicateEntries := []models.Repository{}

	// Iterate through the slice of repositories
	for _, repository := range repositories {
		// If we've already seen this repository, add it to the slice of duplicate entries
		if seenRepositories[repository.RepositoryName] {
			duplicateEntries = append(duplicateEntries, repository)
		} else {
			// Otherwise, mark the repository as seen
			seenRepositories[repository.RepositoryName] = true
		}
	}

	return duplicateEntries
}

// Parse "github.com/mholt/archiver/v3 v3.5.1" into the format "github.com/mholt/archiver/v3@v3.5.1"
func convertToDownloadableFormat(input string) string {
	// Split the input string on the first space character
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	// Return the first part (the package name) followed by an @ symbol and the second part (the version)
	return parts[0] + "@" + parts[1]
}

// Check if the string has uppercase.
func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func parseLibraryUrl(input string) string {
	// Split the input string on the first space character
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	// Split the package name on the '/' character
	packageNameParts := strings.Split(parts[0], "/")

	// Add the '!' prefix and lowercase each part of the package name
	for i, part := range packageNameParts {
		// Split the part on each uppercase character
		subParts := strings.Split(part, "")
		for j, subPart := range subParts {
			if hasUppercase(subPart) {
				subParts[j] = "!" + strings.ToLower(subPart)
			}
		}
		packageNameParts[i] = strings.Join(subParts, "")
	}

	// Join the modified package name parts with '/' characters
	packageName := strings.Join(packageNameParts, "/")

	// Return the modified package name followed by an '@' symbol and the version
	return packageName + "@" + parts[1]
}

// Export the JSON Parser to separate function.
func extractDefaultBranchCommitBlobContent(sourceGraphResponseBody []byte) string {
	outerModFile := JSONParser.Get(
		string(sourceGraphResponseBody),
		"data.repository.defaultBranch.target.commit.blob.content",
	)

	return outerModFile.String()
}
