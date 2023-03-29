package services

import (
	"github.com/Pauloo27/aryzona/internal/db/entity"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var User = &UserService{}

type UserService struct {
}

func (s *UserService) SetPreferredLang(userID string, language i18n.LanguageName) error {
	return upsert(&entity.User{
		ID:              userID,
		PreferredLocale: language,
	})
}

func (s *UserService) SetLastSlashCommandLocale(userID string, language i18n.LanguageName) error {
	return upsert(&entity.User{
		ID:                     userID,
		LastSlashCommandLocale: language,
	})
}
