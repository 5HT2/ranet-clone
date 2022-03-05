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
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	OcrData string `json:"ocr_data"`
}

type configOperation func(*Config)
type Config struct {
	Mutex      sync.Mutex  `json:"-"` // not saved in db
	NumDone    int64       `json:"-"` // not saved in db
	DlQueue    []ImageInfo `json:"-"` // not saved in db
	OcrQueue   []ImageInfo `json:"-"` // not saved in db
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
				//SaveConfig()
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
		log.Printf("failed to marshal config: %v\n", err)
		return
	}

	err = os.WriteFile(configPath, bytes, fileMode)
	if err != nil {
		log.Printf("failed to write config: %v\n", err)
	} else {
		log.Printf("successfully saved config\n")
	}
}

func LoadConfig() {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error loading config: %v\n", err)
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("error unmarshalling config: %v\n", err)
	}
}
