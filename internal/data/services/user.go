package services

import (
	"errors"

	"github.com/pauloo27/aryzona/internal/data/db"
	"github.com/pauloo27/aryzona/internal/data/db/entity"
	"github.com/pauloo27/aryzona/internal/data/db/repos"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
)

var User = &UserService{}

type UserService struct {
}

func (s *UserService) SetPreferredLang(userID string, language i18n.LanguageName) error {
	return repos.User.Upsert(userID, &entity.User{
		ID:              userID,
		PreferredLocale: language,
	})
}

func (s *UserService) SetLastSlashCommandLocale(userID string, language i18n.LanguageName) error {
	return repos.User.Upsert(userID, &entity.User{
		ID:                     userID,
		LastSlashCommandLocale: language,
	})
}

func (s *UserService) GetLanguage(userID, guildID string) i18n.LanguageName {
	user, err := repos.User.FindOneByID(userID)

	userNotFound := errors.Is(err, db.ErrNotFound)

	if err != nil && !userNotFound {
		logger.Error(err)
	}

	if !userNotFound {
		if user.PreferredLocale != "" {
			return user.PreferredLocale
		}

		if user.LastSlashCommandLocale != "" {
			return user.LastSlashCommandLocale
		}
	}

	if guildID == "" {
		return i18n.DefaultLanguageName
	}

	guild, err := repos.Guild.FindOneByID(guildID)
	guildNotFound := errors.Is(err, db.ErrNotFound)

	if err != nil && !guildNotFound {
		logger.Error(err)
	}

	if !guildNotFound {
		return guild.PreferredLocale
	}

	return i18n.DefaultLanguageName
}
