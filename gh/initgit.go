package gh

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/pkg/errors"
)

func InitGitRepo(backupDir, remote, branch string, auth *transport.AuthMethod) (*git.Repository, error) {
	if branch == "" {
		branch = "master"
	}

	if remote == "" {
		return nil, errors.New("remote is required")
	}

	repo, err := git.PlainClone(backupDir, false, &git.CloneOptions{
		URL:           remote,
		Auth:          *auth,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return repo, nil
}
