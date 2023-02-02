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

type CommandNews struct {
	*Common

	Title Entry
}
