package i18n

import (
	"fmt"
	"regexp"
)

type Entry string

func (e Entry) Str(params ...any) string {
	if len(params) == 0 {
		return string(e)
	}

	formattedStr := string(e)
	for i, param := range params {
		re := regexp.MustCompile(fmt.Sprintf(`{%d:[\w]+}`, i))
		formattedStr = re.ReplaceAllString(formattedStr, fmt.Sprint(param))
	}
	return formattedStr
}

type Language struct {
	*Meta
	*Common
	Commands *Commands

	commands map[string]any
}

type Meta struct {
	Name LanguageName

	DisplayName Entry
	Authors     Entry
}

type Common struct {
	HelloWorld         Entry
	SomethingWentWrong Entry
}

type Commands struct {
	Even *CommandEven
	Pick *CommandPick
	News *CommandNews
}

func GetCommand(l *Language, name string) any {
	return l.commands[name]
}
