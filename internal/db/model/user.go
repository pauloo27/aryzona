package model

import "github.com/Pauloo27/aryzona/internal/i18n"

type User struct {
	ID                     string            `xorm:"varchar(255) not null pk"`
	PreferredLocale        i18n.LanguageName `xorm:"varchar(5) null"`
	LastSlashCommandLocale i18n.LanguageName `xorm:"varchar(5) null"`
}
