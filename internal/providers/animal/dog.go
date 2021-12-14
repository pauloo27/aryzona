package animal

import "github.com/Pauloo27/aryzona/internal/utils"

func GetRandomDog() (string, error) {
	url, err := utils.GetString("https://random.dog/woof")
	if err != nil {
		return "", err
	}
	return utils.Fmt("https://random.dog/%s", url), nil
}

func GetRandomDogImage() (string, error) {
	url, err := utils.GetString("https://random.dog/woof?include=jpg")
	if err != nil {
		return "", err
	}
	return utils.Fmt("https://random.dog/%s", url), nil
}
