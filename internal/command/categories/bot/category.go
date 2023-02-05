package bot

import "github.com/Pauloo27/aryzona/internal/command"

var Bot = command.CommandCategory{
	Name:  "bot",
	Emoji: "🤖",
	Commands: []*command.Command{
		&HelpCommand, &PingCommand, &SourceCommand, &UptimeCommand,
		&DonateCommand,
	},
}

func init() {
	command.RegisterCategory(Bot)
}
