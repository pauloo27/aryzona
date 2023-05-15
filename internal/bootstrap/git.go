package bootstrap

import (
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/providers/git"
)

func loadGitInfo(commitHash, commitMessage string) {
	git.CommitHash = commitHash
	git.CommitMessage = commitMessage
	git.RemoteRepo = config.Config.GitRepoURL
}
