package bot

import "github.com/pauloo27/aryzona/internal/command"

var Bot = command.Category{
	Name:  "bot",
	Emoji: "🤖",
	Commands: []*command.Command{
		&HelpCommand, &PingCommand, &SourceCommand, &UptimeCommand,
		&DonateCommand, &LanguageCommand, &ServerCommand,
	},
}

func init() {
	command.RegisterCategory(Bot)
}
