package utils

func AsMention(userID string) string {
	return Fmt("<@%s>", userID)
}
