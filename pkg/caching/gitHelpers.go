package caching

import (
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

	_, err = git.PlainClone(pkg.TemplatesDir, false, cloneOptions)
	if err != nil {
		boillog.LogIt("gitCloner", "WARNING", err.Error())
	}
	return err
}

func gitFetcher() error {
	r, err := git.PlainOpen(pkg.TemplatesDir)
	if err != nil {
		boillog.LogIt("gitFetcher", "ERROR", "Failed to open repository: "+err.Error())
		return err
	}

	fetchOptions := &git.FetchOptions{
		RemoteName:      "origin",
		InsecureSkipTLS: true, // This bypasses SSL certificate verification
	}

	if len(pkg.GitUsername) > 0 && len(pkg.GitToken) > 0 {
		fetchOptions.Auth = &http.BasicAuth{
			Username: pkg.GitUsername,
			Password: pkg.GitToken,
		}
	}

	err = r.Fetch(fetchOptions)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		boillog.LogIt("gitFetcher", "WARNING", "Failed to fetch updates: "+err.Error())
		return err
	}
	boillog.LogIt("gitFetcher", "INFO", "Repository updated successfully")
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
