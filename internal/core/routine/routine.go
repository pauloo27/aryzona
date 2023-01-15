package routine

import "github.com/Pauloo27/logger"

func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()
		f()
	}()
}
