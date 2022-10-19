package models

type Repository struct {
	Id   int    `json: "id" gorm: "primary_key"`
	Name string `json: "name"`
	Uri  string `json: "uri"`
}

type CreateRepositoryInput struct {
	Name string `json: "name" binding:"required"`
	Uri  string `json: "uri" binding:"required"`
}

type UpdateRepositoryInput struct {
	Name string `json: "name" binding:"required"`
	Uri  string `json: "uri" binding:"required"`
}

type Library struct {
	Id     int    `json: "id" gorm: "primary_key"`
	Name   string `json: "name"`
	Uri    string `json: "uri"`
	Liblin int    `json: "liblin"`
}

type CreateLibraryInput struct {
	Name   string `json: "name" binding:"required"`
	Uri    string `json: "uri" binding:"required"`
	Liblin int    `json: "liblin"`
}

type UpdateLibraryInput struct {
	Name   string `json: "name" binding:"required"`
	Uri    string `json: "uri" binding:"required"`
	Liblin int    `json: "liblin"`
}

type Entity struct {
	Id     int    `json: "id" gorm: "primary_key"`
	Name   string `json: "name"`
	Uri    string `json: "uri"`
	Lin    int    `json: "lin"`
	Issct  int    `json: "issct"`
	Liblin int    `json: "liblin"`
}

type CreateEntityInput struct {
	Name   string `json: "name"`
	Uri    string `json: "uri"`
	Lin    int    `json: "lin"`
	Issct  int    `json: "issct"`
	Liblin int    `json: "liblin"`
}

type UpdateEntityInput struct {
	Name   string `json: "name"`
	Uri    string `json: "uri"`
	Lin    int    `json: "lin"`
	Issct  int    `json: "issct"`
	Liblin int    `json: "liblin"`
}

type Owner struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarUrl       string `json:"gravatar_id"`
	Url               string `json:"url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	SiteAdmin         bool   `json:"site_admin"`
}

type License struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	SpdxId  string `json:"spdx_id"`
	NodeId  string `json:"node_id"`
	HtmlUrl string `json:"html_url"`
}

type Item struct {
	Id               int     `json:"id"`
	NodeId           string  `json:"node_id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	Owner            Owner   `json:"owner"`
	Private          bool    `json:"private"`
	HtmlUrl          string  `json:"html_url"`
	Description      string  `json:"description"`
	Fork             bool    `json:"fork"`
	Url              string  `json:"url"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	PushedAt         string  `json:"pushed_at"`
	Homepage         string  `json:"homepage"`
	Size             int     `json:"size"`
	StargazersCount  int     `json:"stargazers_count"`
	WatchersCount    int     `json:"watchers_count"`
	Language         string  `json:"language"`
	ForksCount       int     `json:"forks_count"`
	OpenIssuesCount  int     `json:"open_issues_count"`
	MasterBranch     string  `json:"master_branch"`
	DefaultBranch    string  `json:"default_branch"`
	Score            int     `json:"score"`
	ArchiveUrl       string  `json:"archive_url"`
	AssigneesUrl     string  `json:"assignees_url"`
	BlobsUrl         string  `json:"blobgs_url"`
	BranchesUrl      string  `json:"branches_url"`
	CollaboratorsUrl string  `json:"collaborators_url"`
	CommentsUrl      string  `json:"comments_url"`
	CommitsUrl       string  `json:"commits_url"`
	CompareUrl       string  `json:"compare_url"`
	ContentsUrl      string  `json:"contents_url"`
	ContributorsUrl  string  `json:"contributors_url"`
	DeploymentsUrl   string  `json:"deployments_url"`
	DownloadsUrl     string  `json:"downloads_url"`
	EventsUrl        string  `json:"events_url"`
	ForksUrl         string  `json:"forks_url"`
	GitCommitsUrl    string  `json:"git_commits_url"`
	GitRefsUrl       string  `json:"git_refs_url"`
	GitTagsUrl       string  `json:"git_tags_url"`
	GitUrl           string  `json:"git_url"`
	IssueCommentUrl  string  `json:"issue_comment_url"`
	IssueEventsUrl   string  `json:"issue_events_url"`
	IssuesUrl        string  `json:"issues_url"`
	KeysUrl          string  `json:"keys_url"`
	LabelsUrl        string  `json:"labels_url"`
	LanguagesUrl     string  `json:"languages_url"`
	MergesUrl        string  `json:"merges_url"`
	MilestonesUrl    string  `json:"milestones_url"`
	NotificationsUrl string  `json:"notifications_urls"`
	PullsUrl         string  `json:"pulls_url"`
	ReleasesUrl      string  `json:"releases_url"`
	SshUrl           string  `json:"ssh_url"`
	StargazersUrl    string  `json:"stargazers_url"`
	StatusesUrl      string  `json:"statuses_url"`
	SubscribersUrl   string  `json:"subscripers_url"`
	SubscriptionUrl  string  `json:"subscription_url"`
	TagsUrl          string  `json:"tags_url"`
	TeamsUrl         string  `json:"teams_url"`
	TreesUrl         string  `json:"trees_url"`
	CloneUrl         string  `json:"clone_url"`
	MirrorUrl        string  `json:"mirror_url"`
	HooksUrl         string  `json:"hooks_url"`
	SvnUrl           string  `json:"svn_url"`
	Forks            int     `json:"forks"`
	OpenIssues       int     `json:"open_issues"`
	Watchers         int     `json:"watchers"`
	HasIssues        bool    `json:"has_issues"`
	HasProjects      bool    `json:"has_projects"`
	HasPages         bool    `json:"has_pages"`
	HasWiki          bool    `json:"has_wiki"`
	HasDownloads     bool    `json:"has_downloads"`
	Archived         bool    `json:"archived"`
	Disabled         bool    `json:"disabled"`
	Visibility       bool    `json:"visibility"`
	License          License `json:"license"`
}

type RepositoryResponse struct {
	TotalCount        int    `json:"total_count"`
	IncompleteResults bool   `json:"incomplete_results"`
	Items             []Item `json:"items"`
}

type CodeAnalysisStruct struct {
	Languages []Language `json:"languages"`
	Total     Files      `json:"total"`
}

type Language struct {
	Name    string `json:"name"`
	Files   int    `json:"files"`
	Code    int    `json:"code"`
	Comment int    `json:"comment"`
	Blank   int    `json:"blank"`
}

type Files struct {
	Files   int `json:"files"`
	Code    int `json:"code"`
	Comment int `json:"comment"`
	Blank   int `json:"blank"`
}

type Issue struct {
	Url                   string      `json:"url"`
	RepositoryUrl         string      `json:"repository_url"`
	LabelsUrl             string      `json:"labels_url"`
	CommentsUrl           string      `json:"comments_url"`
	EventsUrl             string      `json:"events_url"`
	HtmlUrl               string      `json:"html_url"`
	Id                    string      `json:"id"`
	NodeId                string      `json:"node_id"`
	Number                int         `json:"number"`
	Title                 string      `json:"title"`
	User                  User        `json:"user"`
	Labels                []string    `json:"labels"`
	State                 string      `json:"state"`
	Locked                bool        `json:"locked"`
	Assignee              string      `json:"assignee"`
	Assignees             []string    `json:"assignees"`
	Milestone             string      `json:"milestone"`
	Comments              int         `json:"comments"`
	CreatedAt             string      `json:"created_at"`
	UpdatedAt             string      `json:"updated_at"`
	ClostedAt             string      `json:"closed_at"`
	AuthorAssociation     string      `json:"author_assocation"`
	ActiveLockReason      string      `json:"active_lock_reason"`
	Draft                 bool        `json:"draft"`
	PullRequest           PullRequest `json:"pull_request"`
	Body                  string      `json:"body"`
	Reactions             Reactions   `json:"reactions"`
	TimelineUrl           string      `json:"timeline_url"`
	PerformedViaGithubApp string      `json:"performed_via_github_app"`
	StateReason           string      `json:"state_reason"`
}

type Reactions struct {
	Url        string `json:"url"`
	TotalCount int    `json:"total_count"`
	PlusOne    int    `json:"+1"`
	MinusOne   int    `json:"-1"`
	Laugh      int    `json:"laugh"`
	Hooray     int    `json:"hooray"`
	Confused   int    `json:"confused"`
	Heart      int    `json:"heart"`
	Rocket     int    `json:"rocket"`
	Eyes       int    `json:"eyes"`
}

type PullRequest struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	MergedAt string `json:"merged_at"`
}

type User struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         string `json:"site_admin"`
}
