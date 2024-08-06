package caching

import (
	"fmt"
	"os"
	"zehd/pkg"

	"github.com/APoniatowski/boillog"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"golang.org/x/sys/unix"
)

func gitCloner() error {
	err := makeDir()
	if err != nil {
		return err
	}

	cloneOptions := &git.CloneOptions{
		URL:             pkg.GitLink,
		InsecureSkipTLS: true, // This bypasses SSL certificate verification
	}

	if len(pkg.GitUsername) > 0 && len(pkg.GitToken) > 0 {
		cloneOptions.Auth = &http.BasicAuth{
			Username: pkg.GitUsername,
			Password: pkg.GitToken,
		}
	}

	if !isGitRepo(pkg.TemplatesDir) {
		_, err = git.PlainClone(pkg.TemplatesDir, false, cloneOptions)
		if err != nil {
			boillog.LogIt("gitCloner", "WARNING", err.Error())
			return err
		}
	}

	return nil
}

func gitPull() error {
	r, err := git.PlainOpen(pkg.TemplatesDir)
	if err != nil {
		boillog.LogIt("gitPull", "ERROR", "Failed to open repository: "+err.Error())
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		boillog.LogIt("gitPull", "ERROR", "Failed to get worktree: "+err.Error())
		return err
	}

	pullOptions := &git.PullOptions{
		RemoteName:      "origin",
		InsecureSkipTLS: true,
	}

	if len(pkg.GitUsername) > 0 && len(pkg.GitToken) > 0 {
		pullOptions.Auth = &http.BasicAuth{
			Username: pkg.GitUsername,
			Password: pkg.GitToken,
		}
	}

	err = w.Pull(pullOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			ref, err := r.Head()
			if err != nil {
				boillog.LogIt("gitPull", "ERROR", "Failed to get HEAD reference: "+err.Error())
				return err
			}
			commitHash := ref.Hash().String()
			boillog.LogIt("gitPull", "INFO", fmt.Sprintf("Repository is already up-to-date at commit %s", commitHash))
			return nil
		}
		boillog.LogIt("gitPull", "WARNING", "Failed to pull updates: "+err.Error())

		return err
	}
	ref, err := r.Head()
	if err != nil {
		boillog.LogIt("gitPull", "ERROR", "Failed to get HEAD reference: "+err.Error())
		return err
	}
	commitHash := ref.Hash().String()

	boillog.LogIt("gitPull", "INFO", fmt.Sprintf("Repository updated successfully to commit %s", commitHash))

	return nil
}

func makeDir() error {
	err := os.Mkdir(pkg.TemplatesDir, 0744)
	if err != nil {
		if os.IsExist(err) {
			boillog.LogIt("gitCloner", "INFO", error.Error(err))

			return nil
		} else if os.IsPermission(err) {
			boillog.LogIt("gitCloner", "ERROR", error.Error(err))
		} else {
			if pathErr, ok := err.(*os.PathError); ok {
				if errno, ok := pathErr.Err.(unix.Errno); ok {
					switch errno {
					case unix.ENOSPC:
						boillog.LogIt("gitCloner", "ERROR", "No space left on device.")

					case unix.EROFS:
						boillog.LogIt("gitCloner", "ERROR", "Read-only file system.")

					default:
						boillog.LogIt("gitCloner", "ERROR", "Unknown error: "+error.Error(err))
					}
				} else {
					boillog.LogIt("gitCloner", "ERROR", "Error creating directory: "+error.Error(err))
				}
			} else {
				boillog.LogIt("gitCloner", "ERROR", "Unexpected error type: "+error.Error(err))
			}
		}

		return err
	}

	return nil
}

func isGitRepo(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	gitDir := path + ".git"
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return false
	}

	return true
}
