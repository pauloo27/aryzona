package sysmon

import "github.com/Pauloo27/aryzona/internal/command"

var SysMon = command.CommandCategory{
	Name:     "System Monitor",
	Emoji:    "💻",
	Commands: []*command.Command{&Sys, &Bash},
}

func init() {
	command.RegisterCategory(SysMon)
}
