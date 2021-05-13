package animal

import "github.com/Pauloo27/aryzona/utils"

func GetRandomDog() (string, error) {
	url, err := utils.GetString("https://random.dog/woof")
	if err != nil {
		return "", err
	}
	return utils.Fmt("https://random.dog/%s", url), nil
}
