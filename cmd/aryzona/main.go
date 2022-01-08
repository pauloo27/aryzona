package main

import "github.com/Pauloo27/aryzona/internal/bootstrap"

var (
	commitHash, commitMessage string
)

func main() {
	bootstrap.Start(commitHash, commitMessage)
}
