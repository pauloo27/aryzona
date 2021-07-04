package utils

import "os/exec"

func GetStreamMetadata(url string) ([]byte, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		url,
	)
	return cmd.Output()
}
