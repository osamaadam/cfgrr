package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:     "clone <remote>",
	Aliases: []string{"c"},
	RunE:    cloneRun,
	Args:    cobra.ExactArgs(1),
	Short:   "Pull the configuration files from the remote git repository",
	Long: `
This command pulls the configuration files from the remote git repository and then replicates them to the backup directory.`,
	Example: strings.Join([]string{
		"cfgrr clone git@github.com:osamaadam/dotfiles.git",
		"cfgrr clone git@github.com:osamaadam/dotfiles.git --branch main",
		"cfgrr clone git@github.com:osamaadam/dotfiles.git -b main",
	}, "\n"),
}

func cloneRun(cmd *cobra.Command, args []string) (err error) {
	config := vconfig.GetConfig()
	url := args[0]

	branchRef := plumbing.NewBranchReferenceName(branch)

	fmt.Println("Cloning the configurations..")
	fmt.Println("Remote:", url)
	fmt.Println("Branch:", branch)
	if _, err := git.PlainClone(config.BackupDir, false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: branchRef,
	}); err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			fmt.Println("Repository already exists, pulling the latest changes..")
			r, err := git.PlainOpen(config.BackupDir)
			if err != nil {
				return err
			}
			w, err := r.Worktree()
			if err != nil {
				return err
			}

			if err := w.Pull(&git.PullOptions{
				RemoteName:    config.GitRemote,
				RemoteURL:     url,
				Progress:      os.Stdout,
				ReferenceName: branchRef,
			}); err != nil {
				if err == git.NoErrAlreadyUpToDate {
					fmt.Println("No changes to pull")
					return nil
				}
				return err
			}

			fmt.Println("Pulled the latest changes")
			return nil
		} else {
			return err
		}
	}

	fmt.Printf("Cloned configurations from %s to %s\n", url, config.BackupDir)
	return nil
}

func init() {
	config := vconfig.GetConfig()
	cloneCmd.Flags().StringVarP(&branch, "branch", "b", config.GitBranch, "branch to clone")
}
