package model

import "github.com/Pauloo27/aryzona/internal/i18n"

type Guild struct {
	ID              string            `xorm:"varchar(255) not null pk"`
	PreferredLocale i18n.LanguageName `xorm:"varchar(5) not null"`
}
