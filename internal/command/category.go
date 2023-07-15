package command

type Category struct {
	Commands    []*Command
	Name, Emoji string
	OnLoad      func()
}
