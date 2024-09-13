package services

import (
	"errors"
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/data/db"
	"github.com/pauloo27/aryzona/internal/data/db/entity"
	"github.com/pauloo27/aryzona/internal/data/db/repos"
	"github.com/pauloo27/aryzona/internal/i18n"
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
		slog.Warn("Cannot load user from db", tint.Err(err))
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
		slog.Error("Cannot load guild from db", "err", err)
	}

	if !guildNotFound {
		return guild.PreferredLocale
	}

	return i18n.DefaultLanguageName
}
