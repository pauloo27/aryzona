package services

import (
	"github.com/Pauloo27/aryzona/internal/db"
	"github.com/Pauloo27/aryzona/internal/db/entity"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var User = &UserService{}

type UserService struct {
}

func (s *UserService) SetPreferredLang(userID string, language i18n.LanguageName) error {
	return s.upsertUser(&entity.User{
		ID:              userID,
		PreferredLocale: language,
	})
}

func (s *UserService) SetLastSlashCommandLocale(userID string, language i18n.LanguageName) error {
	return s.upsertUser(&entity.User{
		ID:              userID,
		LastSlashCommandLocale: language,
	})
}

func (s *UserService) upsertUser(user *entity.User) error {
	// EAFP: easy to ask for forgiveness than permission,
	// let's try to update the user and if it fails, create it

	aff, err := db.DB.Update(user)

	if aff == 1 || err != nil {
		return err
	}

	_, err = db.DB.Insert(user)

	return err
}
