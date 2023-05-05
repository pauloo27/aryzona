package i18n

import (
	"fmt"
	"reflect"
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
	*Locale  `json:"locale"`
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
	Took               Entry

	TranslationInProgress Entry

	MatchInfo   *MatchInfo
	PlayingInfo *PlayingInfo

	DurationNow             Entry
	DurationLessThanAMinute Entry
	DurationDay             Entry
	DurationDays            Entry
	DurationHour            Entry
	DurationHours           Entry
	DurationMinute          Entry
	DurationMinutes         Entry

	DurationSecond  Entry
	DurationSeconds Entry

	Validations *Validations
	Categories  map[string]Entry
}

type MatchInfo struct {
	Match       Entry
	Time        Entry
	TimePenalty Entry
}

type PlayingInfo struct {
	SongTitle         Entry
	Artist            Entry
	Source            Entry
	ETAKey            Entry
	ETANever          Entry
	ETAValue          Entry
	DurationKey       Entry
	DurationLive      Entry
	Position          Entry
	RequestedBy       Entry
	Warning           Entry
	SongPausedWarning Entry
}

type Validations struct {
	MustHaveVoicerOnGuild      *ValidationMustHaveVoicerOnGuild
	MustBePlaying              *ValidationMustBePlaying
	MustBeOnVoiceChannel       *ValidationMustBeOnVoiceChannel
	MustBeOnAValidVoiceChannel *ValidationMustBeOnAValidVoiceChannel
	MustBeOnSameVoiceChannel   *ValidationMustBeOnSameVoiceChannel
	PreCommandValidation       *PreCommandValidation
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
	Dog      *CommandDog
	Cat      *CommandCat
	Fox      *CommandFox
	UUID     *CommandUUID
	Joke     *CommandJoke
	Xkcd     *CommandXkcd
	Language *CommandLanguage
	Server   *CommandServer
}

func GetCommand(l *Language, name string) any {
	return l.commands[name]
}

func MustGetCommandDefinition(l *Language, name string) *CommandDefinition {
	cmd := l.commands[name]
	if cmd == nil {
		return nil
	}

	return reflect.ValueOf(cmd).Elem().FieldByName("Definition").Interface().(*CommandDefinition)
}
