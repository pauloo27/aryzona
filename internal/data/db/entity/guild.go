package entity

import "github.com/pauloo27/aryzona/internal/i18n"

type Guild struct {
	ID              string            `xorm:"id varchar(255) not null pk"`
	PreferredLocale i18n.LanguageName `xorm:"preferred_locale varchar(5) not null"`
}
