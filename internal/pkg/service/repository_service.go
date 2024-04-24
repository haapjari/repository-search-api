package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v61/github"
	"github.com/haapjari/repository-search-api/internal/pkg/model"
	"github.com/haapjari/repository-search-api/internal/pkg/util"
	"log/slog"
	"os"
	"strings"
	"time"
)

type RepositoryService struct {
	QueryParameters *model.QueryParameters

	token      string
	stop       chan struct{}
	errorCh    chan error
	completed  chan *model.Repository
	retryCount int
	*github.Client
}

func NewRepositoryService(token string, params *model.QueryParameters) *RepositoryService {
	g := &RepositoryService{
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

// Query is a method of the RepositoryService struct. It queries GitHub repositories based on the provided query parameters.
// It retrieves detailed information about the repositories and returns the result as a slice of model.Repository structs.
func (rs *RepositoryService) Query() ([]*model.Repository, error) {
	repos, err := rs.multiRepoSearch()
	if err != nil {
		return nil, util.Error(err)
	}

	result := make([]*model.Repository, len(repos))
	rs.completed = make(chan *model.Repository, len(repos))

	for _, r := range repos {
		rs.worker(r)
	}

	for range len(repos) {
		result = append(result, <-rs.completed)
	}

	return result, nil
}

// Stop is a method of the RepositoryService struct. It stops the service by closing the stop channel.
func (rs *RepositoryService) Stop() {
	if rs.stop != nil {
		close(rs.stop)
	}
}

// errorHandler is a method of the RepositoryService struct. It listens for errors that occur during the processing of
// GitHub repositories and logs them using the slog package.
func (rs *RepositoryService) errorHandler() {
	for {
		select {
		case <-rs.stop:
			return
		case err := <-rs.errorCh:
			slog.Error(err.Error())
		}
	}
}

// worker is a method of the RepositoryService struct. It processes a GitHub repository based on the provided repository
// information. It retrieves detailed information about the repository and sends the result to the completed channel.
func (rs *RepositoryService) worker(r *github.Repository) {
	select {
	case <-rs.stop:
		return
	default:
		slog.Debug("Processing: " + r.GetFullName())

		startTime := time.Now()

		pullRequests, err := rs.repoPulls(r.GetFullName())
		if err != nil {
			rs.errorCh <- err
		}

		issues, err := rs.repoIssues(r.GetFullName())
		if err != nil {
			rs.errorCh <- err
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

		commits, err := rs.repoCommits(r.GetFullName())

		if err != nil {
			rs.errorCh <- err
		}

		releases, err := rs.repoLatestRelease(r.GetFullName())
		if err != nil {
			rs.errorCh <- err
		}

		latestRelease, err := rs.repositoryLatestRelease(r.GetFullName())
		if err != nil {
			rs.errorCh <- err

		}

		contributors, err := rs.repoContributors(r.GetFullName())
		if err != nil {
			rs.errorCh <- err
		}

		path, err := util.Clone(rs.token, r.GetCloneURL())
		if err != nil {
			rs.errorCh <- err
		}

		selfWrittenLOC, err := util.CalcLOC(path, r.GetLanguage())
		if err != nil {
			rs.errorCh <- err
		}

		libs, err := util.ParseModFile(path)
		if err != nil {
			rs.errorCh <- err
		}

		err = os.RemoveAll(path)
		if err != nil {
			rs.errorCh <- err
		}

		thirdPartyLOC := 0

		for _, lib := range libs {
			slog.Debug(fmt.Sprintf("Processing %v | Library: %v", r.GetFullName(), lib))

			p, fetchErr := util.FetchLibrary(lib)
			if fetchErr != nil {
				rs.errorCh <- fetchErr
			}

			l, calcErr := util.CalcLOC(p, r.GetLanguage())
			if calcErr != nil {
				rs.errorCh <- calcErr
			}

			thirdPartyLOC += l
		}

		rs.completed <- &model.Repository{
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
			ThirdPartyLOC:          thirdPartyLOC,
			SelfWrittenLOC:         selfWrittenLOC,
		}

		slog.Debug(fmt.Sprintf("Completed Processing: %v | Processing Time: %.2f sec", r.GetFullName(), time.Since(startTime).Seconds()))

		return
	}
}

// repoContributors is a method of the RepositoryService struct. It retrieves detailed information about the contributors
// of a GitHub repository based on the provided full name of the repository. It makes use of the List contributors API
// endpoint (https://docs.github.com/en/rest/reference/repos#list-repository-contributors).
func (rs *RepositoryService) repoContributors(name string) ([]*github.Contributor, error) {
	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.Contributor

	for {
		contributors, resp, err := rs.Client.Repositories.ListContributors(context.Background(), owner, repo, opt)
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

// repositoryLatestRelease is a method of the RepositoryService struct. It retrieves detailed information about the latest
// release of a GitHub repository based on the provided full name of the repository. It makes use of the Get latest
// release API endpoint (https://docs.github.com/en/rest/reference/repos#get-the-latest-release).
func (rs *RepositoryService) repositoryLatestRelease(name string) (*github.RepositoryRelease, error) {
	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	for {
		latestRelease, resp, err := rs.Client.Repositories.GetLatestRelease(context.Background(), owner, repo)
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

// repoLatestRelease is a method of the RepositoryService struct. It retrieves detailed information about the releases
// of a GitHub repository based on the provided full name of the repository. It makes use of the List releases API
// endpoint (https://docs.github.com/en/rest/reference/repos#list-releases).
func (rs *RepositoryService) repoLatestRelease(name string) ([]*github.RepositoryRelease, error) {
	opt := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.RepositoryRelease

	for {
		releases, resp, err := rs.Client.Repositories.ListReleases(context.Background(), owner, repo, opt)
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

// repoIssues is a method of the RepositoryService struct. It retrieves detailed information about the issues
// of a GitHub repository based on the provided full name of the repository. It makes use of the List issues API
// endpoint (https://docs.github.com/en/rest/reference/issues#list-repository-issues).
func (rs *RepositoryService) repoIssues(name string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		State:       "all",
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.Issue

	for {
		r, resp, err := rs.Client.Issues.ListByRepo(context.Background(), owner, repo, opt)
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

// repoPulls is a method of the RepositoryService struct. It retrieves detailed information about the pull requests
// of a GitHub repository based on the provided full name of the repository. It makes use of the List pull requests API
// endpoint (https://docs.github.com/en/rest/reference/pulls#list-pull-requests).
func (rs *RepositoryService) repoPulls(name string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		State:       "all",
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var all []*github.PullRequest

	for {
		r, resp, err := rs.Client.PullRequests.List(context.Background(), owner, repo, opt)
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

// repoCommits is a method of the RepositoryService struct. It retrieves detailed information about the commits
// of a GitHub repository based on the provided full name of the repository. It makes use of the List commits API
// endpoint (https://docs.github.com/en/rest/reference/repos#list-commits).
func (rs *RepositoryService) repoCommits(name string) ([]*github.RepositoryCommit, error) {
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var result []*github.RepositoryCommit

	for {
		r, resp, err := rs.Client.Repositories.ListCommits(context.Background(), owner, repo, opt)
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

// singleRepoSearch is a method of the RepositoryService struct. It retrieves detailed information about a GitHub repository
// based on the provided full name of the repository. It makes use of the Get a Repository API endpoint
// (https://docs.github.com/en/rest/reference/repos#get-a-repository).
//
// The method takes a string argument, fullName, which is the full name of the repository in the format "owner/singleRepoSearch".
// It returns a pointer to a github.Repository struct representing the queried repository and an error if any occurs during the process.
//
// The method handles GitHub's rate limit by retrying the API call after a wait time if a rate limit error is encountered.
// The number of retries and the wait time are defined in the RepositoryService struct.
func (rs *RepositoryService) singleRepoSearch(name string) (*github.Repository, error) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		Order:       rs.QueryParameters.Order,
	}

	owner, repo := strings.Split(name, "/")[0], strings.Split(name, "/")[1]

	var result *github.Repository

	for {
		r, resp, err := rs.Client.Repositories.Get(context.Background(), owner, repo)
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

// multiRepoSearch is a method of the RepositoryService struct. It searches for GitHub repositories based on the provided
// query parameters. It makes use of the Search repositories API endpoint
// (https://docs.github.com/en/rest/reference/search#search-repositories).
func (rs *RepositoryService) multiRepoSearch() ([]*github.Repository, error) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		Order:       rs.QueryParameters.Order,
		Sort:        "stars",
	}

	query := fmt.Sprintf("language:%s stars:%s..%s created:%s..%s", rs.QueryParameters.Language, rs.QueryParameters.MinStars, rs.QueryParameters.MaxStars, rs.QueryParameters.FirstCreationDate, rs.QueryParameters.LastCreationDate)

	all := make([]*github.Repository, 0)

	for {
		r, resp, err := rs.Client.Search.Repositories(context.Background(), query, opt)
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
