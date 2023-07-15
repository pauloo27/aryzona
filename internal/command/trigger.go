package command

import (
	"time"

	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
)

type TriggerEvent struct {
	EventTime                    time.Time
	PreferedLanguage             i18n.LanguageName
	Type                         TriggerType
	Channel                      model.TextChannel
	MessageID, GuildID, AuthorID string
	DeferResponse                func() error
	Reply                        func(*Context, *model.ComplexMessage) error
	Edit                         func(*Context, *model.ComplexMessage) error
}
