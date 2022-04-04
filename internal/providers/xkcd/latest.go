package xkcd

func GetLatest() (*Comic, error) {
	return GetByNum(0)
}
