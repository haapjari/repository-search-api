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
