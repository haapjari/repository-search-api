package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v61/github"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/model"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/util"
	"github.com/hhatto/gocloc"
	"log/slog"
	"os"
	"strings"
	"time"
)

type GitHubService struct {
	QueryParameters *model.QueryParameters

	token      string
	stop       chan struct{}
	errorCh    chan error
	completed  chan *model.Repository
	retryCount int
	*github.Client
}

func NewGitHubService(token string, params *model.QueryParameters) *GitHubService {
	g := &GitHubService{
		QueryParameters: params,
		token:           token,
		errorCh:         make(chan error),
		stop:            make(chan struct{}),
		retryCount:      5,
		Client:          github.NewClient(nil).WithAuthToken(strings.Split(token, " ")[1]),
	}

	go g.errorHandler()

	return g
}

// Query is a method of the GitHubService struct. It queries GitHub repositories based on the provided query parameters.
// It retrieves detailed information about the repositories and returns the result as a slice of model.Repository structs.
func (g *GitHubService) Query() ([]*model.Repository, error) {
	// TODO
	// repos, err := g.multiRepoSearch()
	// if err != nil {
	// 	return nil, util.errorHandler(err)
	// }

	// Temp Code
	repo, err := g.singleRepoSearch("charmbracelet/bubbletea")
	if err != nil {
		return nil, util.Error(err)
	}

	repos := []*github.Repository{repo}

	// Actual Code
	result := make([]*model.Repository, len(repos))
	g.completed = make(chan *model.Repository, len(repos))

	for _, r := range repos {
		g.worker(r)
	}

	for range len(repos) {
		result = append(result, <-g.completed)
	}

	return result, nil
}

// Stop is a method of the GitHubService struct. It stops the service by closing the stop channel.
func (g *GitHubService) Stop() {
	if g.stop != nil {
		close(g.stop)
	}
}

// errorHandler is a method of the GitHubService struct. It listens for errors that occur during the processing of
// GitHub repositories and logs them using the slog package.
func (g *GitHubService) errorHandler() {
	for {
		select {
		case <-g.stop:
			return
		case err := <-g.errorCh:
			slog.Error(err.Error())
		}
	}
}

// worker is a method of the GitHubService struct. It processes a GitHub repository based on the provided repository
// information. It retrieves detailed information about the repository and sends the result to the completed channel.
func (g *GitHubService) worker(r *github.Repository) {
	select {
	case <-g.stop:
		return
	default:
		slog.Debug("Processing: " + r.GetFullName())
		startTime := time.Now()

		pullRequests, err := g.repoPulls(r.GetFullName())
		if err != nil {
			g.errorCh <- err
		}

		issues, err := g.repoIssues(r.GetFullName())
		if err != nil {
			g.errorCh <- err
		}

		openPullRequestCount := 0
		closedPullRequestCount := 0

		for _, pr := range pullRequests {
			if pr.GetState() == "open" {
				openPullRequestCount++
			}

			if pr.GetState() == "closed" {
				closedPullRequestCount++
			}
		}

		commits, err := g.repoCommits(r.GetFullName())

		if err != nil {
			g.errorCh <- err
		}

		releases, err := g.repoLatestRelease(r.GetFullName())
		if err != nil {
			g.errorCh <- err
		}

		latestRelease, err := g.repositoryLatestRelease(r.GetFullName())
		if err != nil {
			g.errorCh <- err

		}

		contributors, err := g.repoContributors(r.GetFullName())
		if err != nil {
			g.errorCh <- err
		}

		path, err := g.clone(r.GetCloneURL())
		if err != nil {
			g.errorCh <- err
		}

		selfWrittenLOC, err := g.calculateLinesOfCode(path, r.GetLanguage())
		if err != nil {
			g.errorCh <- err
		}

		if err = os.RemoveAll(path); err != nil {
			g.errorCh <- err
		}

		g.completed <- &model.Repository{
			Name:                   r.GetName(),
			FullName:               r.GetFullName(),
			CreatedAt:              r.GetCreatedAt().Format("2006-01-02"),
			StargazerCount:         r.GetStargazersCount(),
			Language:               r.GetLanguage(),
			OpenIssues:             r.GetOpenIssues(),
			ClosedIssues:           len(issues),
			OpenPullRequestCount:   openPullRequestCount,
			ClosedPullRequestCount: closedPullRequestCount,
			Forks:                  r.GetForksCount(),
			WatcherCount:           r.GetSubscribersCount(),
			CommitCount:            len(commits),
			NetworkCount:           r.GetNetworkCount(),
			LatestRelease:          latestRelease.GetPublishedAt().Format("2006-01-02"),
			TotalReleasesCount:     len(releases),
			ContributorCount:       len(contributors),
			ThirdPartyLOC:          0, // TODO
			SelfWrittenLOC:         selfWrittenLOC,
		}

		slog.Debug(fmt.Sprintf("Completed Processing: %v | Processing Time: %.2f sec", r.GetFullName(), time.Since(startTime).Seconds()))

		return
	}
}

// repoContributors is a method of the GitHubService struct. It retrieves detailed information about the contributors
// of a GitHub repository based on the provided full name of the repository. It makes use of the List contributors API
// endpoint (https://docs.github.com/en/rest/reference/repos#list-repository-contributors).
func (g *GitHubService) repoContributors(name string) ([]*github.Contributor, error) {
	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.Contributor

	for {
		contributors, resp, err := g.Client.Repositories.ListContributors(context.Background(), owner, repo, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(err)
		}

		all = append(all, contributors...)
		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s/contributors | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return all, nil
}

// repositoryLatestRelease is a method of the GitHubService struct. It retrieves detailed information about the latest
// release of a GitHub repository based on the provided full name of the repository. It makes use of the Get latest
// release API endpoint (https://docs.github.com/en/rest/reference/repos#get-the-latest-release).
func (g *GitHubService) repositoryLatestRelease(name string) (*github.RepositoryRelease, error) {
	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	for {
		latestRelease, resp, err := g.Client.Repositories.GetLatestRelease(context.Background(), owner, repo)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(err)
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s/releases/latest | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		return latestRelease, nil
	}
}

// repoLatestRelease is a method of the GitHubService struct. It retrieves detailed information about the releases
// of a GitHub repository based on the provided full name of the repository. It makes use of the List releases API
// endpoint (https://docs.github.com/en/rest/reference/repos#list-releases).
func (g *GitHubService) repoLatestRelease(name string) ([]*github.RepositoryRelease, error) {
	opt := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.RepositoryRelease

	for {
		releases, resp, err := g.Client.Repositories.ListReleases(context.Background(), owner, repo, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(err)
		}

		all = append(all, releases...)
		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s/releases | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return all, nil
}

// repoIssues is a method of the GitHubService struct. It retrieves detailed information about the issues
// of a GitHub repository based on the provided full name of the repository. It makes use of the List issues API
// endpoint (https://docs.github.com/en/rest/reference/issues#list-repository-issues).
func (g *GitHubService) repoIssues(name string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		State:       "all",
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.Issue

	for {
		r, resp, err := g.Client.Issues.ListByRepo(context.Background(), owner, repo, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(err)
		}

		all = append(all, r...)

		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s/issues | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return all, nil
}

// repoPulls is a method of the GitHubService struct. It retrieves detailed information about the pull requests
// of a GitHub repository based on the provided full name of the repository. It makes use of the List pull requests API
// endpoint (https://docs.github.com/en/rest/reference/pulls#list-pull-requests).
func (g *GitHubService) repoPulls(name string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		State:       "all",
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.PullRequest

	for {
		r, resp, err := g.Client.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(err)
		}

		all = append(all, r...)

		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s/pulls | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return all, nil
}

// repoCommits is a method of the GitHubService struct. It retrieves detailed information about the commits
// of a GitHub repository based on the provided full name of the repository. It makes use of the List commits API
// endpoint (https://docs.github.com/en/rest/reference/repos#list-commits).
func (g *GitHubService) repoCommits(name string) ([]*github.RepositoryCommit, error) {
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var result []*github.RepositoryCommit

	for {
		r, resp, err := g.Client.Repositories.ListCommits(context.Background(), owner, repo, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(fmt.Errorf(": %v", err))
		}

		result = r

		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s/commits | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return result, nil
}

// singleRepoSearch is a method of the GitHubService struct. It retrieves detailed information about a GitHub repository
// based on the provided full name of the repository. It makes use of the Get a Repository API endpoint
// (https://docs.github.com/en/rest/reference/repos#get-a-repository).
//
// The method takes a string argument, fullName, which is the full name of the repository in the format "owner/singleRepoSearch".
// It returns a pointer to a github.Repository struct representing the queried repository and an error if any occurs during the process.
//
// The method handles GitHub's rate limit by retrying the API call after a wait time if a rate limit error is encountered.
// The number of retries and the wait time are defined in the GitHubService struct.
func (g *GitHubService) singleRepoSearch(name string) (*github.Repository, error) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		Order:       g.QueryParameters.Order,
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var result *github.Repository

	for {
		r, resp, err := g.Client.Repositories.Get(context.Background(), owner, repo)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(fmt.Errorf(": %v", err))
		}

		result = r

		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /repos/%s/%s | Response: %v | Rate Limit Left: %v", owner, repo, resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return result, nil
}

// multiRepoSearch is a method of the GitHubService struct. It searches for GitHub repositories based on the provided
// query parameters. It makes use of the Search repositories API endpoint
// (https://docs.github.com/en/rest/reference/search#search-repositories).
func (g *GitHubService) multiRepoSearch() ([]*github.Repository, error) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		Order:       g.QueryParameters.Order,
		Sort:        "stars",
	}

	query := fmt.Sprintf("language:%s stars:%s..%s created:%s..%s", g.QueryParameters.Language, g.QueryParameters.MinStars, g.QueryParameters.MaxStars, g.QueryParameters.FirstCreationDate, g.QueryParameters.LastCreationDate)

	var all []*github.Repository

	for {
		r, resp, err := g.Client.Search.Repositories(context.Background(), query, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				waitDuration := time.Until(rateLimitError.Rate.Reset.Time)
				time.Sleep(waitDuration)
				continue
			}
			return nil, util.Error(fmt.Errorf(": %v", err))
		}

		all = append(all, r.Repositories...)

		if resp.NextPage == 0 {
			break
		}

		slog.Debug(fmt.Sprintf("GET /search/repositories | Response: %v | Rate Limit Left: %v", resp.Status, resp.Rate.Remaining))

		opt.Page = resp.NextPage
	}

	return all, nil
}

// clone is a method of the GitHubService struct. It clones a GitHub repository based on the provided URL.
// It makes use of the go-git library (https://pkg.go.dev/github.com/go-git/go-git/v5).
func (g *GitHubService) clone(url string) (string, error) {
	dir, err := os.MkdirTemp("", "clone-")
	if err != nil {
		return "", util.Error(fmt.Errorf("unable to create a temporary directory: %v", err))
	}

	// fmt.Println("")
	// fmt.Println(strings.Split(g.token, " ")[1])
	// fmt.Println("")

	auth := &http.BasicAuth{
		Username: "x-access-token",
		Password: strings.Split(g.token, " ")[1],
	}

	repo, err := git.PlainClone(dir, false,
		&git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     auth,
		})
	if err != nil {
		return "", util.Error(fmt.Errorf("unable to clone the repository: %v", err))
	}

	if repo != nil {
		return dir, nil
	}

	return "", util.Error(fmt.Errorf("unable to clone the repository"))
}

// calculateLinesOfCode is a method of the GitHubService struct. It calculates the lines of code
// of a directory based on the provided language.
func (g *GitHubService) calculateLinesOfCode(dir string, lang string) (int, error) {
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

	selfWrittenGoCode, ok := result.Languages[lang]
	if ok {
		return int(selfWrittenGoCode.Code), nil
	}

	return 0, fmt.Errorf("calculating lines of code failed")
}

/*

func (r *Repo) CalculateLibraryLOC() (int, error) {
	// TODO
	modFilePath, err := util.FindFile(r.opt.Directory, "go.mod")
	if err != nil {
		return -1, err
	}

	modFile := util.NewModFile(modFilePath)
	err = modFile.Parse()
	if err != nil {
		return -1, err
	}

	// libraries are in format of golang.org/x/mod v0.15.0

	// libraryCount := 0

	goPath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		r.opt.Logger.Errorf("unable to exec command: %s", err.errorHandler())
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
*/

/*
func (r *RepositorySearchService) GetLinesOfCode(name, remote string) (int, int, error) {
	remote = strings.Replace(remote, "git://", "https://", 1)
	baseDir := os.TempDir()
	dir, err := os.MkdirTemp(baseDir, name+"-*")
	if err != nil {
		return -1, -1, err
	}

	singleRepoSearch := NewRepo(&RepoOptions{
		Directory: dir,
		URL:       remote,
		Logger:    r.log,
		Name:      name,
		Token:     r.personalAccessToken,
	})

	r.log.Debugf("Cloning Repository '%v' into %v.", remote, dir)

	err = singleRepoSearch.clone()
	if err != nil {
		return -1, -1, err
	}

	selfWrittenLOC, err := singleRepoSearch.SelfWrittenLOC()
	if err != nil {
		r.log.Errorf("error, while querying for self written lines of code: %s", err.errorHandler())
		return -1, -1, err
	}

	libraryLOC, err := singleRepoSearch.CalculateLibraryLOC()
	if err != nil {
		return -1, -1, err
	}

	r.log.Debugf("Deleting Repository: %v | Directory: %v", remote, dir)

	err = singleRepoSearch.Delete()
	if err != nil {
		return -1, -1, err
	}

	return selfWrittenLOC, libraryLOC, nil
}

func (r *RepositorySearchService) GetCommitsCount(owner, singleRepoSearch string) (int, error) {
	page := 1
	count := 0

	for {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=100&page=%d", owner, singleRepoSearch, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return 0, err
		}

		resp, err := r.httpClient.Do(req)
		if err != nil {
			return count, err
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode == 200 {
			var commits []interface{}
			if err = json.NewDecoder(resp.Body).Decode(&commits); err != nil {
				return count, err
			}

			count += len(commits)

			if linkHeader := resp.Header.Get("Link"); !strings.Contains(linkHeader, `rel="next"`) {
				break
			}

			page++
		} else {
			return count, fmt.Errorf("GitHub API error: %s", resp.Status)
		}
	}

	return count, nil
}

func (r *RepositorySearchService) GetOpenPullRequests(owner, singleRepoSearch string) (int, error) {
	page := 1
	count := 0

	for {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=open&per_page=100&page=%d",
				owner, singleRepoSearch, page), nil)
		if err != nil {
			return 0, err
		}

		resp, e := r.httpClient.Do(req)
		if e != nil {
			return count, e
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode == 200 {
			var pullRequests []interface{}
			if err = json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
				return count, err
			}

			count += len(pullRequests)

			if linkHeader := resp.Header.Get("Link"); !strings.Contains(linkHeader, `rel="next"`) {
				break
			} else {
				page = page + 1
			}
		}
	}

	return count, nil
}

func (r *RepositorySearchService) GetClosedPullRequests(owner, singleRepoSearch string) (int, error) {
	page := 1
	count := 0

	for {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=closed&per_page=100&page=%d",
				owner, singleRepoSearch, page), nil)
		if err != nil {
			return 0, err
		}

		resp, e := r.httpClient.Do(req)
		if e != nil {
			return count, e
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode == 200 {
			var pullRequests []interface{}
			if err = json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
				return count, err
			}

			count += len(pullRequests)

			if linkHeader := resp.Header.Get("Link"); !strings.Contains(linkHeader, `rel="next"`) {
				break
			} else {
				page = page + 1
			}
		}
	}

	return count, nil
}

func (r *RepositorySearchService) GetTotalReleases(owner string, name string) (int, error) {
	opt := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	var all []*github.RepositoryRelease

	for {
		releases, resp, err := r.gitHubClient.Repositories.ListReleases(context.Background(), owner, name, opt)
		if err != nil {
			return -1, err
		}

		all = append(all, releases...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return len(all), nil
}

func (r *RepositorySearchService) GetLatestRelease(owner string, name string) (string, error) {
	release, _, err := r.gitHubClient.Repositories.GetLatestRelease(context.Background(), owner, name)
	if err != nil {
		return "", err
	}

	return release.GetCreatedAt().String(), nil
}

func (r *RepositorySearchService) GetContributorCount(owner string, name string) (int, error) {
	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
		Anon: "true",
	}

	var all []*github.Contributor

	for {
		contributors, resp, err := r.gitHubClient.Repositories.ListContributors(context.Background(),
			owner, name, opt)
		if err != nil {
			return -1, err
		}

		all = append(all, contributors...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return len(all), nil
}

func (r *RepositorySearchService) LastCreationDate(language string, stars string) (string, int, error) {
	if language == "" || stars == "" {
		return "", 400, errors.New("language or stars field is empty")
	}

	var (
		date  = time.Now()
		year  = date.Year()
		month = int(date.Month())
		day   = date.Day()
	)

	startDate, e := time.Parse(time.DateOnly, fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if e != nil {
		return "", 500, e
	}

	var totalCount int

	for {
		select {
		case <-time.After(r.queryInterval):
			var (
				queryParameters = fmt.Sprintf("q=language:%s+stars:%s+created:>%s&order=asc",
					language, stars, startDate.Format(time.DateOnly))
				endpoint = fmt.Sprintf("https://api.github.com/search/repositories?%s", queryParameters)
			)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				return "", 500, err
			}

			r.log.Debugf("Request: %v", endpoint)

			if r.personalAccessToken != "" {
				req.Header.Set("Authorization", "token "+r.personalAccessToken)
			}

			resp, err := r.httpClient.Do(req)
			if err != nil {
				return "", 500, err
			}

			r.log.Debugf("Response Status Code: %v", resp.StatusCode)

			if resp.StatusCode == 200 {
				count := Count{}

				body, readError := io.ReadAll(resp.Body)
				if readError != nil {
					return "", 500, err
				}

				if err = json.Unmarshal(body, &count); err != nil {
					return "", 500, err
				}

				totalCount = count.TotalCount

				if totalCount != 0 {
					return startDate.Format(time.DateOnly), 200, nil
				}

				startDate = startDate.AddDate(0, 0, -1)
			}

			if err = resp.Body.Close(); err != nil {
				r.log.Warnf("unable to close response body: %s", err.errorHandler())
			}
		}
	}
}

func (r *RepositorySearchService) FirstCreationDate(language string, stars string) (string, int, error) {
	if language == "" || stars == "" {
		return "", 400, errors.New("language or stars field is empty")
	}

	year, err := r.findFirstYear(language, stars)
	if err != nil {
		return "", 500, err
	}

	month, err := r.findFirstMonth(language, stars, year)
	if err != nil {
		return "", 500, err
	}

	weekDay, err := r.findFirstWeek(language, stars, year, month)
	if err != nil {
		return "", 500, err
	}

	day, err := r.findFirstDay(language, stars, year, month, weekDay)
	if err != nil {
		return "", 500, err
	}

	return fmt.Sprintf("%d-%02d-%02d", year, month, day), 200, nil
}

func (r *RepositorySearchService) findFirstYear(language string, stars string) (int, error) {
	startDate, e := time.Parse("2006-01-02", "2007-01-01")
	if e != nil {
		return -1, e
	}

	var totalCount int

	for {
		select {
		case <-time.After(r.queryInterval):
			var (
				queryParameters = fmt.Sprintf("q=language:%s+stars:%s+created:<%s&order=asc",
					language, stars, startDate.Format("2006-01-02"))
				endpoint = fmt.Sprintf("https://api.github.com/search/repositories?%s", queryParameters)
			)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Request: %v", endpoint)

			if r.personalAccessToken != "" {
				req.Header.Set("Authorization", "token "+r.personalAccessToken)
			}

			resp, err := r.httpClient.Do(req)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Response Status Code: %v", resp.StatusCode)

			if resp.StatusCode == 200 {
				count := Count{}

				body, readError := io.ReadAll(resp.Body)
				if readError != nil {
					return -1, err
				}

				if err = json.Unmarshal(body, &count); err != nil {
					return -1, err
				}

				totalCount = count.TotalCount

				if totalCount != 0 {
					return startDate.Year() - 1, nil
				}

				startDate = startDate.AddDate(1, 0, 0)
			}

			if err = resp.Body.Close(); err != nil {
				r.log.Warnf("unable to close response body: %s", err.errorHandler())
			}
		}
	}
}

func (r *RepositorySearchService) findFirstMonth(language string, stars string, year int) (int, error) {
	startDate, e := time.Parse("2006-01-02", fmt.Sprintf("%v-01-01", year))
	if e != nil {
		return -1, e
	}

	totalCount := -1

	for {
		select {
		case <-time.After(r.queryInterval):
			var (
				queryParameters = fmt.Sprintf("q=language:%s+stars:%s+created:<%s&order=asc",
					language, stars, startDate.Format("2006-01-02"))
				endpoint = fmt.Sprintf("https://api.github.com/search/repositories?%s", queryParameters)
			)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Request: %v", endpoint)

			if r.personalAccessToken != "" {
				req.Header.Set("Authorization", "token "+r.personalAccessToken)
			}

			resp, err := r.httpClient.Do(req)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Response Status Code: %v", resp.StatusCode)

			if resp.StatusCode == 200 {
				count := Count{}

				body, readError := io.ReadAll(resp.Body)
				if readError != nil {
					return -1, err
				}

				if err = json.Unmarshal(body, &count); err != nil {
					return -1, err
				}

				totalCount = count.TotalCount

				if totalCount != 0 {
					return int(startDate.Month() - 1), nil
				}

				startDate = startDate.AddDate(0, 1, 0)
			}

			if err = resp.Body.Close(); err != nil {
				r.log.Warnf("unable to close response body: %s", err.errorHandler())
			}
		}
	}
}

// findFirstWeek
func (r *RepositorySearchService) findFirstWeek(language string, stars string, year int, month int) (int, error) {
	startDate, e := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", year, month, 1))
	if e != nil {
		return -1, e
	}

	totalCount := -1

	for {
		select {
		case <-time.After(r.queryInterval):
			var (
				queryParameters = fmt.Sprintf("q=language:%s+stars:%s+created:<%s&order=asc",
					language, stars, startDate.Format("2006-01-02"))
				endpoint = fmt.Sprintf("https://api.github.com/search/repositories?%s", queryParameters)
			)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Request: %v", endpoint)

			if r.personalAccessToken != "" {
				req.Header.Set("Authorization", "token "+r.personalAccessToken)
			}

			resp, err := r.httpClient.Do(req)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Response Status Code: %v", resp.StatusCode)

			if resp.StatusCode == 200 {
				count := Count{}

				body, readError := io.ReadAll(resp.Body)
				if readError != nil {
					return -1, err
				}

				if err = json.Unmarshal(body, &count); err != nil {
					return -1, err
				}

				totalCount = count.TotalCount

				if totalCount != 0 {
					return startDate.AddDate(0, 0, -7).Day(), nil
				}

				startDate = startDate.AddDate(0, 0, 7)
			}

			if err = resp.Body.Close(); err != nil {
				r.log.Warnf("unable to close response body: %s", err.errorHandler())
			}
		}
	}
}

func (r *RepositorySearchService) findFirstDay(language string, stars string, year int, month int, day int) (int, error) {
	startDate, e := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if e != nil {
		return -1, e
	}

	totalCount := -1

	for {
		select {
		case <-time.After(r.queryInterval):
			var (
				queryParameters = fmt.Sprintf("q=language:%s+stars:%s+created:<%s&order=asc",
					language, stars, startDate.Format("2006-01-02"))
				endpoint = fmt.Sprintf("https://api.github.com/search/repositories?%s", queryParameters)
			)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Request: %v", endpoint)

			if r.personalAccessToken != "" {
				req.Header.Set("Authorization", "token "+r.personalAccessToken)
			}

			resp, err := r.httpClient.Do(req)
			if err != nil {
				return -1, err
			}

			r.log.Debugf("Response Status Code: %v", resp.StatusCode)

			if resp.StatusCode == 200 {
				count := Count{}

				body, readError := io.ReadAll(resp.Body)
				if readError != nil {
					return -1, err
				}

				if err = json.Unmarshal(body, &count); err != nil {
					return -1, err
				}

				totalCount = count.TotalCount

				if totalCount != 0 {
					return startDate.AddDate(0, 0, -1).Day(), nil
				}

				startDate = startDate.AddDate(0, 0, 1)
			}

			if err = resp.Body.Close(); err != nil {
				r.log.Warnf("unable to close response body: %s", err.errorHandler())
			}
		}
	}
}
*/
