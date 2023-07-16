package services

import (
	"github.com/pauloo27/aryzona/internal/data/db/entity"
	"github.com/pauloo27/aryzona/internal/data/db/repos"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var Guild = &GuildService{}

type GuildService struct {
}

func (s *GuildService) SetGuildOptions(guildID string, language i18n.LanguageName) error {
	return repos.Guild.Upsert(guildID, &entity.Guild{
		ID:              guildID,
		PreferredLocale: language,
	})
}
