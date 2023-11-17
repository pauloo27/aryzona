package routine

import (
	"log/slog"
)

func GoAndRecover(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Go routine error recovered", "err", err)
			}
		}()
		f()
	}()
}
