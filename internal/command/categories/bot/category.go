package bot

import "github.com/Pauloo27/aryzona/internal/command"

var Bot = command.CommandCategory{
	Name:  "Bot",
	Emoji: "🤖",
	Commands: []*command.Command{
		&HelpCommand, &PingCommand, &SourceCommand, &UptimeCommand,
	},
}

func init() {
	command.RegisterCategory(Bot)
}
