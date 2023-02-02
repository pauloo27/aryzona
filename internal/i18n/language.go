package i18n

type Entry string

// TODO: add "params" for formatting...
func (e Entry) Str() string {
	return string(e)
}

type Language struct {
	Meta     *Meta
	Common   *Common
	Commands *Commands

	commands map[string]any
}

type Meta struct {
	Name LanguageName

	DisplayName Entry
	Authors     Entry
}

type Common struct {
	HelloWorld Entry
}

type Commands struct {
	Even *CommandEven
}

func GetCommand(l *Language, name string) any {
	return l.commands[name]
}
