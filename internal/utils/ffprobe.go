package utils

import "os/exec"

/* #nosec G204 */
func GetStreamMetadata(url string) ([]byte, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		// oh shit... some bad shit can happen if that came from the user...
		url,
	)
	return cmd.Output()
}
