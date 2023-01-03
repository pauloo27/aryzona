package model

type MessageComponent any

type ButtonComponent struct {
	Label    string
	ID       string
	Emoji    string
	Style    ButtonStyle
	Disabled bool
}

type ButtonStyle int

const (
	_ ButtonStyle = iota
	PrimaryButtonStyle
	SecondaryButtonStyle
	SuccessButtonStyle
	DangerButtonStyle
)
