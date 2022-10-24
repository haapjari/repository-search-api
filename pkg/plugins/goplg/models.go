package goplg

type GitHubResponse struct {
	Data GitHubDataStruct `json:"data"`
}

type GitHubDataStruct struct {
	Repository GitHubRepositoryStruct `json:"repository"`
}

type GitHubRepositoryStruct struct {
	DefaultBranchRef DefaultBranchRefStruct   `json:"defaultBranchRef"`
	OpenIssues       GitHubOpenIssuesStruct   `json:"openIssues"`
	ClosedIssues     GitHubClosedIssuesStruct `json:"closedIssues"`
	Languages        GitHubLanguagesStruct    `json:"languages"`
}

type DefaultBranchRefStruct struct {
	Target GitHubTargetStruct `json:"target"`
}

type GitHubTargetStruct struct {
	History GitHubHistoryStruct `json:"history"`
}

type GitHubHistoryStruct struct {
	TotalCount int `json:"totalCount"`
}

type GitHubOpenIssuesStruct struct {
	TotalCount int `json:"totalCount"`
}

type GitHubClosedIssuesStruct struct {
	TotalCount int `json:"totalCount"`
}

type GitHubLanguagesStruct struct {
	TotalSize int `json:"totalSize"`
}
