package commons

import (
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var syncOnce sync.Once

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {
	// load .env file
	syncOnce.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			//log.Fatalf("Error loading .env file")
		}
	})

	return os.Getenv(key)
}
