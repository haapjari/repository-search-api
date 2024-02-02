package models

type RepositoryResponse struct {
	RepositoryData []Repository `json:"data"`
}

type Repository struct {
	Id                   int    `json:"id" gorm:"primary_key"`
	RepositoryName       string `json:"repository_name"`
	RepositoryUrl        string `json:"repository_url"`
	OpenIssueCount       string `json:"open_issue_count"`
	ClosedIssueCount     string `json:"closed_issue_count"`
	CommitCount          string `json:"commit_count"`
	OriginalCodebaseSize string `json:"original_codebase_size"`
	LibraryCodebaseSize  string `json:"library_codebase_size"`
	RepositoryType       string `json:"repository_type"`
	PrimaryLanguage      string `json:"primary_language"`
	CreationDate         string `json:"creation_date"`
	StargazerCount       string `json:"stargazer_count"`
	LicenseInfo          string `json:"license_info"`
	LatestRelease        string `json:"latest_release"`
}

type CreateRepositoryInput struct {
	RepositoryName       string `json:"repository_name"`
	RepositoryUrl        string `json:"repository_url"`
	OpenIssueCount       string `json:"open_issue_count"`
	ClosedIssueCount     string `json:"closed_issue_count"`
	CommitCount          string `json:"commit_count"`
	OriginalCodebaseSize string `json:"original_codebase_size"`
	LibraryCodebaseSize  string `json:"library_codebase_size"`
	RepositoryType       string `json:"repository_type"`
	PrimaryLanguage      string `json:"primary_language"`
	CreationDate         string `json:"creation_date"`
	StargazerCount       string `json:"stargazer_count"`
	LicenseInfo          string `json:"license_info"`
	LatestRelease        string `json:"latest_release"`
}

type UpdateRepositoryInput struct {
	RepositoryName       string `json:"repository_name"`
	RepositoryUrl        string `json:"repository_url"`
	OpenIssueCount       string `json:"open_issue_count"`
	ClosedIssueCount     string `json:"closed_issue_count"`
	CommitCount          string `json:"commit_count"`
	OriginalCodebaseSize string `json:"original_codebase_size"`
	LibraryCodebaseSize  string `json:"library_codebase_size"`
	RepositoryType       string `json:"repository_type"`
	PrimaryLanguage      string `json:"primary_language"`
	CreationDate         string `json:"creation_date"`
	StargazerCount       string `json:"stargazer_count"`
	LicenseInfo          string `json:"license_info"`
	LatestRelease        string `json:"latest_release"`
}

type Commit struct {
	Id             int    `json:"id" gorm:"primary_key"`
	RepositoryName string `json:"repository_name"`
	CommitDate     string `json:"commit_date"`
	CommitUser     string `json:"commit_user"`
}

type CreateCommitInput struct {
	RepositoryName string `json:"repository_name"`
	CommitDate     string `json:"commit_date"`
	CommitUser     string `json:"commit_user"`
}

type UpdateCommitInput struct {
	RepositoryName string `json:"repository_name"`
	CommitDate     string `json:"commit_date"`
	CommitUser     string `json:"commit_user"`
}

type GitHubResponseStruct struct {
	Data GitHubResponseDataStruct `json:"data"`
}

type GitHubResponseDataStruct struct {
	Repository GitHubRepositoryStruct `json:"repository"`
}

type GitHubRepositoryStruct struct {
	DefaultBranchRef GitHubDefaultBranch   `json:"defaultBranchRef"`
	OpenIssues       GitHubOpenIssues      `json:"openIssues"`
	ClosedIssues     GitHubClosedIssues    `json:"closedIssues"`
	Languages        GitHubLanguages       `json:"languages"`
	StargazerCount   int                   `json:"stargazerCount"`
	CreatedAt        string                `json:"createdAt"`
	PrimaryLanguage  GitHubPrimaryLanguage `json:"primaryLanguage"`
	LicenseInfo      GitHubLicenseInfo     `json:"licenseInfo"`
	LatestRelease    GitHubLatestRelease   `json:"latestRelease"`
}

type GitHubLatestRelease struct {
	PublishedAt string `json:"publishedAt"`
}

type GitHubLicenseInfo struct {
	Key string `json:"key"`
}

type GitHubPrimaryLanguage struct {
	Name string `json:"name"`
}

type GitHubDefaultBranch struct {
	Target GitHubTarget `json:"target"`
}

type GitHubTarget struct {
	History GitHubHistory `json:"history"`
}

type GitHubHistory struct {
	TotalCount int `json:"totalCount"`
}

type GitHubOpenIssues struct {
	TotalCount int `json:"totalCount"`
}

type GitHubClosedIssues struct {
	TotalCount int `json:"totalCount"`
}

type GitHubLanguages struct {
	TotalSize int `json:"totalSize"`
}

type SourceGraphResponseStruct struct {
	Data SourceGraphDataStruct `json:"data"`
}

type SourceGraphDataStruct struct {
	Search SourceGraphSearch `json:"search"`
}

type SourceGraphSearch struct {
	Results SourceGraphResults `json:"results"`
}

type SourceGraphResults struct {
	Repositories []SourceGraphRepositories `json:"repositories"`
}

type SourceGraphRepositories struct {
	Name string `json:"name"`
}
