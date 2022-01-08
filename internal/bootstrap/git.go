package bootstrap

import (
	"os"

	"github.com/Pauloo27/aryzona/internal/providers/git"
)

func loadGitInfo(commitHash, commitMessage string) {
	git.CommitHash = commitHash
	git.CommitMessage = commitMessage
	git.RemoteRepo = os.Getenv("DC_BOT_REMOTE_REPO")
}
