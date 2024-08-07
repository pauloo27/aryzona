package i18n

type CommandDefinition struct {
	Name        Entry
	Description Entry

	Parameters  []ParameterDefinition
	SubCommands []CommandDefinition
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

type CommandServer struct {
	*Common
	Definition *CommandDefinition

	Title                Entry
	ServerOptionsChanged Entry
}

type CommandLanguage struct {
	*Common
	*Meta
	Definition *CommandDefinition

	Title                Entry
	LanguageNotFound     Entry
	LanguageChanged      Entry
	LanguageList         Entry
	CurrentLanguage      Entry
	UsingDefaultLanguage Entry
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

	MatchNotFound      Entry
	MatchFinished      Entry
	AlreadyFollowing   Entry
	FollowLimitReached Entry
}

type CommandLive struct {
	*Common
	Definition *CommandDefinition

	NoMatchesLive Entry
	Title         Entry
	Page          Entry
	PageNotFound  Entry
	Footer        Entry
}

type CommandScore struct {
	*Common
	Definition *CommandDefinition

	MatchNotFound Entry
	LiveUpdates   Entry
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

	Title   Entry
	Unknown Entry
	SeeMore Entry
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
	*Locale
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

type CommandVolume struct {
	*Common
	Definition *CommandDefinition
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

	ListTitle     Entry
	CannotConnect Entry
	AddedToQueue  Entry
	ListFooter    Entry
}

type CommandPlaying struct {
	*Common
	Definition *CommandDefinition

	Title      Entry
	Entry      Entry
	ComingNext Entry
	Song       Entry
	Songs      Entry
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
	RawMap     RawJSONMap
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
	CancelBtn                Entry
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
