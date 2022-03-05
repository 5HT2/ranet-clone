package cfg

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var (
	config     Config
	configPath = "config.json"
	fileMode   = os.FileMode(0700)
)

type ImageInfo struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type configOperation func(*Config)
type Config struct {
	Mutex      sync.Mutex  `json:"-"` // not saved in db
	NumDone    int64       `json:"-"` // not saved in db
	Downloaded []ImageInfo `json:"downloaded,omitempty"`
}

// Config.run will modify a Config non-concurrently.
func (c *Config) run(co configOperation) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	co(c)
}

// SetupConfigSaving will run SaveConfig every minute with a ticker
func SetupConfigSaving() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				SaveConfig()
			}
		}
	}()
}

func SaveConfig() {
	var bytes []byte
	var err error = nil

	config.run(func(c *Config) {
		bytes, err = json.MarshalIndent(c, "", "    ")
	})

	if err != nil {
		log.Printf("Failed to marshal config: %v\n", err)
		return
	}

	err = os.WriteFile(configPath, bytes, fileMode)
	if err != nil {
		log.Printf("Failed to write config: %v\n", err)
	} else {
		log.Printf("Successfully saved config\n")
	}
}

func LoadConfig() {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("Error unmarshalling config: %v\n", err)
	}
}
