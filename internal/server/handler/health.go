package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/data/db"
	"github.com/pauloo27/aryzona/internal/discord"
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
		slog.Error("Cannot encode status", tint.Err(err))
	}
}
