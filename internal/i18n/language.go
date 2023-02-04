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
	*Meta    `json:"meta"`
	*Common  `json:"common"`
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
	Even     *CommandEven
	Pick     *CommandPick
	News     *CommandNews
	Roll     *CommandRoll
	Follow   *CommandFollow
	UnFollow *CommandUnFollow
	Score    *CommandScore
	CNPJ     *CommandCNPJ
	CPF      *CommandCPF
	Password *CommandPassword
	Source   *CommandSource
	Ping     *CommandPing
	Donate   *CommandDonate
	Uptime   *CommandUptime
	Resume   *CommandResume
	Shuffle  *CommandShuffle
	Skip     *CommandSkip
	Stop     *CommandStop
	Lyric    *CommandLyric
	Pause    *CommandPause
	Playing  *CommandPlaying
	Radio    *CommandRadio
	Help     *CommandHelp
	Play     *CommandPlay
}

func GetCommand(l *Language, name string) any {
	return l.commands[name]
}
