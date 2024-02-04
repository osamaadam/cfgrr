package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	RunE:    pushRun,
	Short:   "Push the configuration files to the remote git repository",
	Long: `
This command automatically replicates the files in the backup directory so that they are browsable. and then pushes the changes to the remote git repository.
For this to work properly, the user must have already set up the global git configuration, and the remote repository must exist.`,
	Example: strings.Join([]string{
		"cfgrr push",
		"cfgrr push origin",
		"cfgrr push origin/master",
		"cfgrr push origin master",
	}, "\n"),
}

func pushRun(cmd *cobra.Command, args []string) (err error) {
	config := vconfig.GetConfig()
	remote, branch := "", ""
	if len(args) > 0 {
		remote = args[0]
		if len(args) > 1 {
			branch = args[1]
		}
	}
	if remote == "" {
		remote = config.GitRemote
	}
	if branch == "" {
		branch = config.GitBranch
	}
	// TODO: Initialize the remote if it doesn't exist.
	repo, err := git.PlainInit(config.BackupDir, false)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			// Repository already exists.
			// Open the repository.
			repo, err = git.PlainOpen(config.BackupDir)
			if err != nil {
				// Failed to open the repository.
				return err
			}
		} else {
			// Failed to initialize the repository.
			return err
		}
	}

	// Stage the changes.
	w, err := repo.Worktree()

	if branch != "" {
		if err := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
			Create: true,
			Keep:   true,
		}); err != nil {
			if err != git.ErrBranchExists {
				if err := w.Checkout(&git.CheckoutOptions{
					Branch: plumbing.NewBranchReferenceName(branch),
					Keep:   true,
				}); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Replicate the files to make them browsable.
	all, clean = true, true
	if err := runReplicate(cmd, nil); err != nil {
		return err
	}

	if _, err := w.Add("."); err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		fmt.Println("No changes to push")
		return nil
	}

	commitMsg := fmt.Sprintf("cfgrr push (%s)", time.Now().Format(time.RFC1123))

	if _, err := w.Commit(commitMsg, &git.CommitOptions{}); err != nil {
		return err
	}

	fmt.Println(commitMsg)

	fmt.Println("Pushing to", remote)

	if err := repo.Push(&git.PushOptions{
		RemoteName: remote,
		Progress:   os.Stdout,
	}); err != nil {
		return err
	}

	fmt.Println("Pushed to", remote)

	return nil
}
