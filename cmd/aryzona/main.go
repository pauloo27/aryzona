package main

import "github.com/pauloo27/aryzona/internal/bootstrap"

var (
	commitHash, commitMessage string
)

func main() {
	bootstrap.Start(commitHash, commitMessage)
}
