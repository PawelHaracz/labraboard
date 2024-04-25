package models

type CommitType string

var (
	SHA    CommitType = "SHA"
	BRANCH CommitType = "BRANCH"
	TAG    CommitType = "TAG"
)
