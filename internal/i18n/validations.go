package i18n

type ValidationMustHaveVoicerOnGuild struct {
	Description     Entry
	BotNotConnected Entry
}

type ValidationMustBePlaying struct {
	Description    Entry
	NothingPlaying Entry
}

type ValidationMustBeOnVoiceChannel struct {
	Description             Entry
	YouAreNotInVoiceChannel Entry
}

type ValidationMustBeOnAValidVoiceChannel struct {
	Description            Entry
	CannotConnectToChannel Entry
}

type ValidationMustBeOnSameVoiceChannel struct {
	Description       Entry
	NotInRightChannel Entry
}

type PreCommandValidation struct {
	MustBeExecutedAsSlashCommand Entry
	PermissionRequired           Entry
	MissingSubCommand            Entry
	InvalidSubCommand            Entry
}
