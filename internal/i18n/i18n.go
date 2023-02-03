package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Pauloo27/logger"
)

type LanguageName string

const (
	EnglishLang    LanguageName = "en_US"
	PortugueseLang LanguageName = "pt_BR"

	DefaultLanguageName = EnglishLang
)

var (
	I18nRootDir = "./assets/i18n"

	loadedLanguages = make(map[LanguageName]*Language)
)

func GetLanguage(name LanguageName) (*Language, error) {
	lang, ok := loadedLanguages[name]
	if !ok {
		var err error
		lang, err = loadLanguage(name)
		if err != nil {
			return nil, err
		}
		loadedLanguages[name] = lang
	}
	return lang, nil
}

func MustGetLanguage(name LanguageName) *Language {
	lang, err := GetLanguage(name)
	if err != nil {
		logger.Fatal(err)
	}
	return lang
}

/* #nosec G304 the name does not come from user input */
func loadLanguage(name LanguageName) (*Language, error) {
	fileName := fmt.Sprintf("%s/%s.json", I18nRootDir, name)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var lang Language
	err = json.Unmarshal(data, &lang)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(lang.Commands).Elem()

	lang.commands = make(map[string]any)

	commonValue := reflect.ValueOf(lang.Common)
	metaValue := reflect.ValueOf(lang.Meta)

	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		fieldValue := reflect.ValueOf(lang.Commands).Elem().Field(i)

		commonField := fieldValue.Elem().FieldByName("Common")
		if commonField.IsValid() {
			commonField.Set(commonValue)
		}

		metaField := fieldValue.Elem().FieldByName("Meta")
		if metaField.IsValid() {
			metaField.Set(metaValue)
		}

		lang.commands[strings.ToLower(structField.Name)] = fieldValue.Interface()
	}

	return &lang, err
}
