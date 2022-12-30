package goplg

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"github.com/hhatto/gocloc"
)

func (g *GoPlugin) writeSourceGraphResponseToDatabase(length int, repositories []SourceGraphRepositoriesStruct) {
	var wg sync.WaitGroup

	// Semaphore is a safeguard to goroutines, to allow only "MaxThreads" run at the same time.
	semaphore := make(chan int, g.MaxThreads)

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

	// close(semaphore) // TODO
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
func runGoCloc(path string) int {
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
