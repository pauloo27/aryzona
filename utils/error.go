package utils

func HandleFatal(err error) {
	if err != nil {
		panic(err)
	}
}
