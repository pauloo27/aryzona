package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pauloo27/aryzona/internal/db"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/logger"
)

func Health(w http.ResponseWriter, r *http.Request) {
	dbErr := db.DB.Ping()

	isDBOk := dbErr == nil
	isDiscordOk := discord.Bot.IsLive()

	status := map[string]bool{
		"db":      isDBOk,
		"discord": isDiscordOk,
	}

	isAlright := isDBOk && isDiscordOk

	if isAlright {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		logger.Error(err)
	}
}
