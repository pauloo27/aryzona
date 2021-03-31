package sysmon

import "github.com/Pauloo27/aryzona/command"

var SysMon = command.Category{
	Name:     "System Monitor",
	Commands: []*command.Command{&Sys, &Eval},
}
