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

type CommandRoll struct {
	*Common

	Dice        Entry
	Dices       Entry
	Face        Entry
	Faces       Entry
	Title       Entry
	Description Entry
}

type CommandUnFollow struct {
	*Common

	NotFollowingAny   Entry
	UnFollowedAll     Entry
	MatchNotFound     Entry
	NotFollowingMatch Entry
	UnfollowedMatch   Entry
}

type CommandFollow struct {
	*Common

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

	Title       Entry
	WithMask    Entry
	WithoutMask Entry
}

type CommandCNPJ struct {
	*Common

	Title       Entry
	WithMask    Entry
	WithoutMask Entry
}

type CommandPassword struct {
	*Common

	Title       Entry
	Description Entry
}

type CommandNews struct {
	*Common

	Title Entry
}

type CommandPing struct {
	*Common

	Title            Entry
	Footer           Entry
	APILatency       Entry
	StillCalculating Entry
}

type CommandDonate struct {
	*Common

	Title Entry
}

type CommandUptime struct {
	*Common

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

	Description Entry
}

type CommandResume struct {
	*Common

	NotPaused Entry
	Resumed   Entry
}

type CommandShuffle struct {
	*Common

	Shuffled Entry
}

type CommandSkip struct {
	*Common

	Skipped Entry
}

type CommandStop struct {
	*Common

	Stopped Entry
}

type CommandPause struct {
	*Common

	CannotPause   Entry
	AlreadyPaused Entry
	Paused        Entry
}

type CommandRadio struct {
	*Common

	ListTitle         Entry
	CannotConnect     Entry
	NotInRightChannel Entry
	AddedToQueue      Entry
	ListFooter        Entry
}

type CommandPlaying struct {
	*Common

	Title      Entry
	Never      Entry
	Entry      Entry
	ComingNext Entry
	AndMore    Entry
}

type CommandLyric struct {
	*Common

	NothingPlaying Entry
	NoResults      Entry
	NotConnected   Entry
}
