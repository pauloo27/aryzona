package bootstrap

import (
	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/providers/git"
)

func loadGitInfo(commitHash, commitMessage string) {
	git.CommitHash = commitHash
	git.CommitMessage = commitMessage
	git.RemoteRepo = config.Config.GitRepoURL
}
