package i18n

import (
	"fmt"
	"strconv"
	"strings"
)

type Entry string

func (e Entry) Str(params ...any) string {
	if len(params) == 0 {
		return string(e)
	}

	formattedStr := string(e)
	for i, param := range params {
		formattedStr = strings.ReplaceAll(formattedStr, "{"+strconv.Itoa(i)+"}", fmt.Sprint(param))
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
