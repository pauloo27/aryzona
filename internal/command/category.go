package command

type CommandCategory struct {
	Commands    []*Command
	Name, Emoji string
	OnLoad      func()
}
