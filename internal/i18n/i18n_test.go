package i18n_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/stretchr/testify/assert"
)

var (
	entryType = reflect.TypeOf(i18n.Entry(""))
)

func TestEntryFormat(t *testing.T) {
	var entry i18n.Entry

	entry = i18n.Entry("Hello")
	assert.Equal(t, "Hello", entry.Str())

	entry = i18n.Entry("Hello {0:name}")
	assert.Equal(t, "Hello {0:name}", entry.Str())

	entry = i18n.Entry("Hello {0:name}")
	assert.Equal(t, "Hello World", entry.Str("World"))

	entry = i18n.Entry("Hello {0:name}")
	assert.Equal(t, "Hello World", entry.Str("World", "Invalid"))

	entry = i18n.Entry("Hello {0:name}, {1:greet}")
	assert.Equal(t, "Hello World, welcome", entry.Str("World", "welcome"))

	entry = i18n.Entry("Hello {0:name}, {1:greet}")
	assert.Equal(t, "Hello World, welcome", entry.Str("World", "welcome"))

	entry = i18n.Entry("Hello {0:name}, {1:greet}")
	assert.Equal(t, "Hello 123, true", entry.Str(123, true))
}

func TestDefaultLangs(t *testing.T) {
	i18n.I18nRootDir = "../../assets/i18n"

	l, err := i18n.GetLanguage(i18n.DefaultLanguageName)
	assert.Nil(t, err)
	assert.NotNil(t, l)

	assert.Equal(t, "en_US", string(l.Name))

	lType := reflect.TypeOf(l).Elem()
	lValue := reflect.ValueOf(l).Elem()

	missingTranslations := checkForMissingTranslations(lType, lValue, "")

	assert.Empty(t, missingTranslations)
}

func TestOtherLangs(t *testing.T) {
	i18n.I18nRootDir = "../../assets/i18n"

	for _, lang := range i18n.LanguagesName {
		if lang == i18n.DefaultLanguageName {
			continue
		}
		t.Run(string(lang), func(t *testing.T) {
			l, err := i18n.GetLanguage(lang)
			assert.Nil(t, err)
			assert.NotNil(t, l)

			lType := reflect.TypeOf(l).Elem()
			lValue := reflect.ValueOf(l).Elem()

			missingTranslations := checkForMissingTranslations(lType, lValue, "")

			assert.Empty(t, missingTranslations)
		})
	}
}

func checkForMissingTranslations(t reflect.Type, value reflect.Value, parentPath string) (missing []string) {
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		fieldValue := value.Field(i)

		if structField.Type == entryType {
			if fieldValue.Interface().(i18n.Entry).Str() == "" {
				fmt.Printf("Missing translation for %s%s\n", parentPath, structField.Name)
				missing = append(missing, structField.Name)
			}
			continue
		}

		path := fmt.Sprintf("%s%s.", parentPath, structField.Name)

		if structField.Type.Kind() == reflect.Struct {
			missing = append(missing, checkForMissingTranslations(structField.Type, fieldValue, path)...)
			continue
		}

		if structField.Type.Kind() == reflect.Ptr {
			if structField.Type.Elem().Kind() == reflect.Struct {
				missing = append(missing, checkForMissingTranslations(structField.Type.Elem(), fieldValue.Elem(), path)...)
			}
			continue
		}
	}
	return missing
}
