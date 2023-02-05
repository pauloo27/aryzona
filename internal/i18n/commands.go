package i18n

type CommandDefinition struct {
	Name        Entry
	Description Entry

	Parameters []ParameterDefinition
}

type ParameterDefinition struct {
	Name        Entry
	Description Entry
}

type CommandEven struct {
	*Common
	Definition *CommandDefinition

	Even Entry
	Odd  Entry
}

type CommandPick struct {
	*Common
	Definition *CommandDefinition

	Title       Entry
	Description Entry
}

type CommandRoll struct {
	*Common
	Definition *CommandDefinition

	Dice        Entry
	Dices       Entry
	Face        Entry
	Faces       Entry
	Title       Entry
	Description Entry
}

type CommandUnFollow struct {
	*Common
	Definition *CommandDefinition

	NotFollowingAny   Entry
	UnFollowedAll     Entry
	MatchNotFound     Entry
	NotFollowingMatch Entry
	UnfollowedMatch   Entry
}

type CommandFollow struct {
	*Common
	Definition *CommandDefinition

	Match              Entry
	Time               Entry
	TimePenalty        Entry
	MatchNotFound      Entry
	MatchFinished      Entry
	AlreadyFollowing   Entry
	FollowLimitReached Entry
}

type CommandScore struct {
	*Common
	Definition *CommandDefinition

	NoMatchesLive Entry
	Title         Entry
	Footer        Entry
	MatchNotFound Entry
	LiveUpdates   Entry
	Time          Entry
	Match         Entry
	TimePenalty   Entry
}

type CommandCPF struct {
	*Common
	Definition *CommandDefinition

	Title       Entry
	WithMask    Entry
	WithoutMask Entry
}

type CommandCNPJ struct {
	*Common
	Definition *CommandDefinition

	Title       Entry
	WithMask    Entry
	WithoutMask Entry
}

type CommandPassword struct {
	*Common
	Definition *CommandDefinition

	Title       Entry
	Description Entry
}

type CommandNews struct {
	*Common
	Definition *CommandDefinition

	Title Entry
}

type CommandPing struct {
	*Common
	Definition *CommandDefinition

	Title            Entry
	Footer           Entry
	APILatency       Entry
	StillCalculating Entry
}

type CommandDonate struct {
	*Common
	Definition *CommandDefinition

	Title Entry
}

type CommandUptime struct {
	*Common
	*Meta
	Definition *CommandDefinition

	Language       Entry
	Title          Entry
	Uptime         Entry
	Implementation Entry
	HostInfoKey    Entry
	HostInfoValue  Entry
	StartedAt      Entry
	LastCommit     Entry
}

type CommandSource struct {
	*Common
	Definition *CommandDefinition

	Description Entry
}

type CommandResume struct {
	*Common
	Definition *CommandDefinition

	NotPaused Entry
	Resumed   Entry
}

type CommandShuffle struct {
	*Common
	Definition *CommandDefinition

	Shuffled Entry
}

type CommandSkip struct {
	*Common
	Definition *CommandDefinition

	Skipped Entry
}

type CommandStop struct {
	*Common
	Definition *CommandDefinition

	Stopped Entry
}

type CommandPause struct {
	*Common
	Definition *CommandDefinition

	CannotPause   Entry
	AlreadyPaused Entry
	Paused        Entry
}

type CommandRadio struct {
	*Common
	Definition *CommandDefinition

	ListTitle         Entry
	CannotConnect     Entry
	NotInRightChannel Entry
	AddedToQueue      Entry
	ListFooter        Entry
}

type CommandPlaying struct {
	*Common
	Definition *CommandDefinition

	Title      Entry
	Never      Entry
	Entry      Entry
	ComingNext Entry
	AndMore    Entry
}

type CommandLyric struct {
	*Common
	Definition *CommandDefinition

	NothingPlaying Entry
	NoResults      Entry
	NotConnected   Entry
}

type CommandHelp struct {
	*Common
	Definition *CommandDefinition

	Title              Entry
	Parameters         Entry
	Validations        Entry
	SubCommands        Entry
	Aliases            Entry
	Permission         Entry
	Category           Entry
	AKA                Entry
	ForMoreInfo        Entry
	CommandNotFound    Entry
	SubCommandNotFound Entry
	RequiresPermission Entry
	Required           Entry
	NotRequired        Entry
}

type CommandPlay struct {
	*Common
	Definition *CommandDefinition

	NotInRightChannel        Entry
	CannotConnect            Entry
	YouTubePlaylist          Entry
	BestResult               Entry
	MultipleResults          Entry
	MultipleResultsSelectOne Entry
	FirstResultWillPlay      Entry
	IfYouFailToSelect        Entry
	Entry                    Entry
	SelectedResult           Entry
	ConfirmBtn               Entry
	PlayOtherBtn             Entry
	Live                     Entry
}

type CommandFox struct {
	Definition *CommandDefinition
}

type CommandDog struct {
	Definition *CommandDefinition
}

type CommandCat struct {
	Definition *CommandDefinition
}

type CommandUUID struct {
	Definition *CommandDefinition
}

type CommandJoke struct {
	Definition *CommandDefinition
}

type CommandXkcd struct {
	Definition *CommandDefinition
}
