package utils

func ConditionalString(b bool, tr, fal string) string {
	if b {
		return tr
	}
	return fal
}
