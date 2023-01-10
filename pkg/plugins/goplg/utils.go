package goplg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unicode"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"github.com/hhatto/gocloc"
	JSONParser "github.com/tidwall/gjson"
)

func (g *GoPlugin) writeSourceGraphResponseToDatabase(length int, repositories []SourceGraphRepositoriesStruct) {
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
func (g *GoPlugin) getAllRepositories() models.RepositoryResponse {
	getRepositoriesRequest, err := http.NewRequest("GET", REPOSITORY_API_BASEURL, nil)
	utils.CheckErr(err)

	getRepositoriesRequest.Header.Set("Content-Type", "application/json")

	getRepositoriesResponse, err := g.HttpClient.Do(getRepositoriesRequest)
	utils.CheckErr(err)

	defer getRepositoriesResponse.Body.Close()

	getRepositoriesResponseBody, err := ioutil.ReadAll(getRepositoriesResponse.Body)
	utils.CheckErr(err)

	var repositories models.RepositoryResponse
	json.Unmarshal([]byte(getRepositoriesResponseBody), &repositories)

	return repositories
}

// Filter empty strings from slice.
func filterEmpty(slice []string) []string {
	var result []string
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

// Performs a GET request to the specified URL.
func performGetRequest(url string) string {
	// Make a GET request to the specified URL
	resp, err := http.Get(url)
	utils.CheckErr(err)

	defer resp.Body.Close()

	// Read the response body into a variable
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckErr(err)

	return (string(body))
}

// Calculates the lines of code using https://github.com/hhatto/gocloc
// in the path provided and return the value.
func runGocloc(path string) int {
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

// Delete the contents of tmp/ -folder.
func pruneTempGoPath() {
	// Delete the processed libraries from the disk, in order not to suffocate.
	out, err := runCommand("chmod", "-R", "777", utils.GetTempGoPath())
	if out != "" {
		fmt.Println(out)
	}

	if err != "" {
		fmt.Println(err)
	}

	os.RemoveAll(utils.GetTempGoPath())
	os.MkdirAll(utils.GetTempGoPath(), os.ModePerm)
}

// Wrapper for "exec/os" command execution.
// Copied from blog: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func runCommand(name string, arg ...string) (string, string) {
	cmd := exec.Command(name, arg...)

	var stdout, stderr []byte
	var errStdout, errStderr error

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	err := cmd.Start()
	utils.CheckErr(err)

	// WaitGroup ensures, cmd.Wait() is called, after we finish reading from stdin and stdout.
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	outStr, errStr := string(stdout), string(stderr)

	return outStr, errStr
}

// Helper function for running commands with "os/exec".
// Copied from blog: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// removeDuplicates removes duplicates from a slice of strings
func removeDuplicates(slice []string) []string {
	// Create a map to keep track of the elements that have already been seen
	seen := make(map[string]struct{}, len(slice))

	// Initialize the result slice
	var result []string

	// Iterate over the slice
	for _, elem := range slice {
		// Check if the element has already been seen
		if _, ok := seen[elem]; !ok {
			// If it has not been seen, add it to the result slice and mark it as seen
			result = append(result, elem)
			seen[elem] = struct{}{}
		}
	}

	return result
}

// Parse repository name from url.
func parseRepositoryName(s string) (string, string, error) {
	parts := strings.Split(s, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid repository string: %s", s)
	}
	return parts[1], parts[2], nil
}

// Creates a slice of repositories, which are duplicates in an original list.
func findDuplicateRepositoryEntries(repositories []models.Repository) []models.Repository {
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
func parseUrlToDownloadFormat(input string) string {
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

func parseGoLibraryUrl(input string) string {
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
