package info

import (
	"os"

	"github.com/naeimc/health/api"
)

func Complete(form map[string]any) (completed map[string]any) {
	completed = make(map[string]any)
	for key := range form {
		value := completeKey(key)
		if value != nil {
			completed[key] = value
		}
	}
	return
}

func completeKey(key string) any {
	switch key {
	case api.Key_Hostname:
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "error: " + err.Error()
		}
		return hostname
	case api.Key_DirUsrHome:
		dir, err := os.UserHomeDir()
		if err != nil {
			dir = "error: " + err.Error()
		}
		return dir
	case api.Key_DirUsrConfig:
		dir, err := os.UserConfigDir()
		if err != nil {
			dir = "error: " + err.Error()
		}
		return dir
	case api.Key_DirUsrCache:
		dir, err := os.UserCacheDir()
		if err != nil {
			dir = "error: " + err.Error()
		}
		return dir
	default:
		return nil
	}
}
