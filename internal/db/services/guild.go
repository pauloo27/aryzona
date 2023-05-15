package services

import (
	"github.com/pauloo27/aryzona/internal/db/entity"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var Guild = &GuildService{}

type GuildService struct {
}

func (s *GuildService) SetGuildOptions(guildID string, language i18n.LanguageName) error {
	return upsert(&entity.Guild{
		ID:              guildID,
		PreferredLocale: language,
	})
}
