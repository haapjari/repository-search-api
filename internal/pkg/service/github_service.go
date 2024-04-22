package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v61/github"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/model"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/util"
	"log/slog"
	"os"
	"strings"
	"time"
)

type GitHubService struct {
	QueryParameters *model.QueryParameters

	*github.Client
}

func NewGitHubService(token string, params *model.QueryParameters) *GitHubService {
	return &GitHubService{
		QueryParameters: params,
		Client:          github.NewClient(nil).WithAuthToken(strings.Split(token, " ")[1]),
	}
}

func (g *GitHubService) Query() ([]*model.Repository, error) {
	repos, err := g.repositoriesSearch()
	if err != nil {
		return nil, util.Error(err)
	}

	var result []*model.Repository

	for _, repo := range repos {
		result = append(result, &model.Repository{
			Name:           repo.GetName(),
			FullName:       repo.GetFullName(),
			CreatedAt:      repo.GetCreatedAt().Format("2006-01-02"),
			StargazerCount: repo.GetStargazersCount(),
			Language:       repo.GetLanguage(),
			OpenIssues:     repo.GetOpenIssuesCount(),
			// ClosedIssues:           repo.GetClosedIssuesCount(),  // TODO
			// OpenPullRequestCount:   repo.GetOpenPullRequests(),   // TODO
			// ClosedPullRequestCount: repo.GetClosedPullRequests(), // TODO
			Forks:              repo.GetForksCount(),
			SubscriberCount:    repo.GetSubscribersCount(),
			WatcherCount:       repo.GetWatchersCount(),
			CommitCount:        0,  // TODO
			NetworkCount:       0,  // TODO
			LatestRelease:      "", // TODO
			TotalReleasesCount: 0,  // TODO
			ContributorCount:   0,  // TODO
			ThirdPartyLOC:      0,  // TODO
			SelfWrittenLOC:     0,  // TODO
		})

	}
	// Rate Limits (?)

	// err := g.repositoriesSearch()
	// if err != nil {
	// 	return nil, util.Error(err)
	// }

	return result, nil
}

func (g *GitHubService) repositoriesSearch() ([]*github.Repository, error) {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		Order:       g.QueryParameters.Order,
		Sort:        "stars",
	}

	query := fmt.Sprintf("language:%s stars:%s..%s created:%s..%s", g.QueryParameters.Language, g.QueryParameters.MinStars, g.QueryParameters.MaxStars, g.QueryParameters.FirstCreationDate, g.QueryParameters.LastCreationDate)

	var (
		result     []*github.Repository
		retryCount int
	)

	for {
		r, resp, err := g.Client.Search.Repositories(context.Background(), query, opt)
		if err != nil {
			var rateLimitError *github.RateLimitError
			if errors.As(err, &rateLimitError) {
				retryCount++
				time.Sleep(10 * time.Second)
				if retryCount < 5 {
					continue
				}
			}
			return nil, util.Error(fmt.Errorf(": %v", err))
		}

		slog.Debug("Response: " + resp.Status)
		slog.Debug(fmt.Sprintf("Rate Limit Left: %v", resp.Rate.Remaining))

		result = append(result, r.Repositories...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return result, nil
}

func (g *GitHubService) clone(url string) (string, error) {
	dir, err := os.MkdirTemp("", "clone-")
	if err != nil {
		return "", util.Error(fmt.Errorf("unable to create a temporary directory: %v", err))
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     nil,
	})
	if err != nil {
		return "", util.Error(fmt.Errorf("unable to clone the repository: %v", err))
	}

	if repo != nil {
		return dir, nil
	}

	return "", util.Error(fmt.Errorf("unable to clone the repository"))
}

/*


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
*/
