package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pauloo27/logger"
)

type LanguageName string

const (
	EnglishLang    LanguageName = "en_US"
	PortugueseLang LanguageName = "pt_BR"

	DefaultLanguageName = EnglishLang
)

var (
	I18nRootDir = "./assets/i18n"

	LanguagesName = []LanguageName{EnglishLang, PortugueseLang}

	loadedLanguages = make(map[LanguageName]*Language)
)

func (l LanguageName) DiscordName() string {
	return strings.Replace(string(l), "_", "-", 1)
}

func FindLanguageName(name string) LanguageName {
	langName := strings.Split(name, "_")[0]

	switch langName {
	case "en":
		return EnglishLang
	case "pt":
		return PortugueseLang
	default:
		return DefaultLanguageName
	}
}

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

	var rawMap RawJSONMap
	err = json.Unmarshal(data, &rawMap)
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
	lang.langName = name
	lang.RawMap = rawMap

	commonValue := reflect.ValueOf(lang.Common)
	metaValue := reflect.ValueOf(lang.Meta)
	localeValue := reflect.ValueOf(lang.Locale)
	rawMapValue := reflect.ValueOf(rawMap)

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

		localeField := fieldValue.Elem().FieldByName("Locale")
		if localeField.IsValid() {
			localeField.Set(localeValue)
		}

		rawMapField := fieldValue.Elem().FieldByName("RawMap")
		if rawMapField.IsValid() {
			rawMapField.Set(rawMapValue)
		}

		lang.commands[strings.ToLower(structField.Name)] = fieldValue.Interface()
	}

	return &lang, err
}

type RawJSONMap map[string]any

var (
	ErrInvalidKeyType   = errors.New("invalid key type")
	ErrInvalidMapAccess = errors.New("invalid map access")
	ErrInvalidArrayAccess = errors.New("invalid array access")
)

func (r RawJSONMap) Get(keys ...any) (any, error) {
	if len(keys) == 0 {
		return r, nil
	}

	value := r[keys[0].(string)]

	for _, key := range keys[1:] {
		switch k := key.(type) {
		case string:
			typedV, ok := value.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("%w: %v (%T)", ErrInvalidMapAccess, key, key)
			}
			value = typedV[k]
		case int:
			typedV, ok := value.([]any)
			if !ok {
				return nil, fmt.Errorf("%w: %v (%T)", ErrInvalidArrayAccess, key, key)
			}
			value = typedV[k]
		default:
			return nil, fmt.Errorf("%w: %v (%T)", ErrInvalidKeyType, key, key)
		}
	}

	return value, nil
}
