package iac

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"labraboard/internal/models"
	"os"
)

type Git struct {
	folderPath string
	repoUrl    string
	repository *git.Repository
	gitSha     string
	commitType models.CommitType
}

func GitClone(repoUrl string, folderPath string, commitName string, commitType models.CommitType) (*Git, error) {
	gitRepo, err := git.PlainClone(folderPath, false, &git.CloneOptions{
		URL:      repoUrl,
		Tags:     git.AllTags,
		Progress: nil, //os.Stdout,
	})
	if err != nil {
		return nil, err
	}
	w, err := gitRepo.Worktree()
	if err != nil {
		return nil, err
	}

	var commitSha = ""
	switch commitType {
	case models.TAG:
		tag, err := gitRepo.Tag(commitName)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannot checkin tag %s", commitName))
		}
		checkoutOptions := git.CheckoutOptions{
			Hash:  tag.Hash(),
			Force: true,
		}
		if err = w.Checkout(&checkoutOptions); err != nil {
			mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", commitName, commitName)
			if err = fetchOrigin(gitRepo, mirrorRemoteBranchRefSpec); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Cannot check in branch %s", commitName))
			}
		}
		head, err := gitRepo.Head()
		if err != nil {
			return nil, err
		}
		commitSha = head.Hash().String()
	case models.SHA:
		object, err := gitRepo.CommitObject(plumbing.NewHash(commitName))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannot checkin commit %s", commitName))
		}
		checkoutOptions := git.CheckoutOptions{
			Hash:  object.Hash,
			Force: true,
		}
		if err = w.Checkout(&checkoutOptions); err != nil {
			mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", commitName, commitName)
			if err = fetchOrigin(gitRepo, mirrorRemoteBranchRefSpec); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Cannot check in branch %s", commitName))
			}
		}
		commitSha = object.Hash.String()
	case models.BRANCH:
		branchRefName := plumbing.NewBranchReferenceName(commitName)
		branchCoOpts := git.CheckoutOptions{
			Branch: branchRefName,
			Force:  true,
		}

		if err = w.Checkout(&branchCoOpts); err != nil {
			mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", commitName, commitName)
			if err = fetchOrigin(gitRepo, mirrorRemoteBranchRefSpec); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Cannot check in branch %s", commitName))
			}
		}
		ref, err := gitRepo.Reference(branchRefName, false)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannot read referances from branch %s", commitName))
		}
		commitSha = ref.Hash().String()
	}

	g := &Git{
		folderPath: folderPath,
		repoUrl:    repoUrl,
		repository: gitRepo,
		gitSha:     commitSha,
	}

	return g, nil

}

func (g *Git) Clear() error {
	err := os.RemoveAll(g.folderPath)
	if err != nil {
		return err
	}
	return nil
}

func (g *Git) GetCommitSha() string {
	return g.gitSha
}

func fetchOrigin(repo *git.Repository, refSpecStr string) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}
