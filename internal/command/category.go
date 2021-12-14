package command

type CommandCategory struct {
	Name, Emoji string
	Commands    []*Command
	OnLoad      func()
}
