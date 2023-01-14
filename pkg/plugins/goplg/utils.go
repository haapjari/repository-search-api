package goplg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"github.com/hhatto/gocloc"
	"github.com/pingcap/errors"
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

func (g *GoPlugin) checkIfRepositoryExists(r models.Repository) bool {
	existingRepos := g.getAllRepositories()

	for _, existingRepo := range existingRepos {
		if existingRepo.RepositoryUrl == r.RepositoryUrl {
			return true
		}
	}

	return false
}

func (g *GoPlugin) hasBeenEnriched(r models.Repository) bool {
	existingRepos := g.getAllRepositories()

	for _, existingRepo := range existingRepos {
		if existingRepo.RepositoryUrl == r.RepositoryUrl {
			if existingRepo.OpenIssueCount == "" && existingRepo.ClosedIssueCount == "" && existingRepo.CommitCount == "" && existingRepo.RepositoryType == "" && existingRepo.PrimaryLanguage == "" && existingRepo.CreationDate == "" && existingRepo.StargazerCount == "" && existingRepo.LicenseInfo == "" && existingRepo.LatestRelease == "" {
				return false
			}
		}
	}

	return true
}

func (g *GoPlugin) getAllRepositories() []models.Repository {
	var repositories []models.Repository

	g.DatabaseClient.Find(&repositories)

	return repositories
}

// Calculates the lines of code using https://github.com/hhatto/gocloc
// in the path provided and return the value.
func (g *GoPlugin) calculateCodeLines(path string) (int, error) {
	if !folderExists(path) {
		return 0, errors.New("Error - " + path + " doesn't exist.")
	}

	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	paths := []string{
		path,
	}

	processor := gocloc.NewProcessor(languages, options)

	result, err := processor.Analyze(paths)
	utils.CheckErr(err)

	return int(result.Total.Code), nil
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

// Check if a folder exists in the file system.
func folderExists(folderPath string) bool {
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
func convertToDownloadableFormat(input string) string {
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

// Export the JSON Parser to separate function.
func extractDefaultBranchCommitBlobContent(sourceGraphResponseBody []byte) string {
	outerModFile := JSONParser.Get(
		string(sourceGraphResponseBody),
		"data.repository.defaultBranch.target.commit.blob.content",
	)

	return outerModFile.String()
}
