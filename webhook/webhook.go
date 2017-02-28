package webhook

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

const (
	defaultGitBranch = "master"
	localRepoPath    = "/hans/ansible"
)

type pullPayload struct {
	Ref string `json:"ref" binding:"required"`
}

func RefreshGitRepo(c *gin.Context) {
	gitBranch := os.Getenv("HANS_GIT_BRANCH")
	if gitBranch == "" {
		gitBranch = defaultGitBranch
	}

	var p pullPayload
	c.BindJSON(&p)

	if isRefreshRequired(gitBranch, p.Ref) {
		err := gitPull(p.Ref)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "git pull failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "done",
		"message": "updated local repo",
	})
}

func isRefreshRequired(gitBranch, webhookRef string) bool {
	if webhookRef != "refs/heads/"+gitBranch {
		return false
	}

	return true
}

func gitPull(ref string) error {
	cmd := exec.Command("git", "-C", localRepoPath, "pull", "--rebase", "origin", ref)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
