package i18n

type CommandEven struct {
	*Common

	Even Entry
	Odd  Entry
}

type CommandPick struct {
	*Common

	Title       Entry
	Description Entry
}

type CommandNews struct {
	*Common

	Title Entry
}
