package repos

import (
	"github.com/pauloo27/aryzona/internal/data/db"
	"github.com/pauloo27/aryzona/internal/data/db/entity"
)

type (
	UserRepo  = db.Repository[entity.User, string]
	GuildRepo = db.Repository[entity.Guild, string]
)

var (
	User  = UserRepo{}
	Guild = GuildRepo{}
)
