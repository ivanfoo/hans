package utils

import (
	"os"
	"os/exec"
)

const (
	repositoryLocalPath = "/hans/ansible"
	remoteName          = "origin"
)

type repository struct {
	localPath string
	remote    string
	branch    string
}

func CloneGitRepository(remote string, branch string) {
	r := &repository{
		localPath: repositoryLocalPath,
		remote:    remote,
		branch:    branch,
	}

	if r.existsLocally() {
		r.gitPull()
		return
	}

	r.gitClone()
}

func (r *repository) existsLocally() bool {
	_, err := os.Stat(r.localPath)

	if os.IsNotExist(err) {
		return false
	}

	return true

}

func (r *repository) gitClone() error {
	cmd := exec.Command("git", "clone", r.remote, r.localPath)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) gitPull() error {
	cmd := exec.Command("git", "-C", r.localPath, "pull", remoteName, r.branch)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
