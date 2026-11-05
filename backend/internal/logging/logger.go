// Package logging configure le logger structuré (slog).
package logging

import (
	"log/slog"
	"os"
)

// Init configure et retourne le logger par défaut :
// JSON en production, texte lisible sinon.
func Init(appEnv string) *slog.Logger {
	var h slog.Handler
	if appEnv == "production" {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	l := slog.New(h)
	slog.SetDefault(l)
	return l
}
