package valueobjects

import "errors"

var (
	ErrRepoUrlEmpty = errors.New("error: url is empty")
)

type IaCRepo struct {
	Url           string
	DefaultBranch string
	Path          string
}

func NewIaCRepo(url string, defaultBranch string, path string) (*IaCRepo, error) {
	if url == "" {
		return nil, ErrRepoUrlEmpty
	}
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	return &IaCRepo{
		Url:           url,
		DefaultBranch: defaultBranch,
		Path:          path,
	}, nil
}
