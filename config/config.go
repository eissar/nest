package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type NestConfig struct {
	Port           int    `json:"port"`
	Host           string `json:"host"`
	AddFromPath    string `json:"addFromPath"`
	GetLibraryInfo string `json:"getLibraryInfo"`
	GetThumbnail   string `json:"getThumbnail"`
	ApiKey         string `json:"apiKey"`
}

func (n NestConfig) FmtURL() string {
	return "http://" + n.Host + ":" + strconv.Itoa(n.Port)
}

// gets config path,
// creating paths if they don't exist
// or panics.
func GetConfigPath() string {
	prf, ok := os.LookupEnv("userprofile")
	if !ok {
		log.Fatalf("config.getPath: `userprofile` env variable is nil ?")
	}

	configPath := filepath.Join(prf, ".config", "nest")

	createFlag := false
	_, err := os.Stat(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			createFlag = true
		} else {
			log.Fatalf("config.get: %v", err.Error())
		}
	}
	if createFlag {
		err := os.MkdirAll(configPath, 0666) //perm bits do nothing windows
		if err != nil {
			log.Fatalf("config.getPath: %v", err.Error())
		}
	}
	return configPath
}
func GetPath() string {
	prf, ok := os.LookupEnv("userprofile")
	if !ok {
		log.Fatalf("config.getPath: `userprofile` env variable is nil ?")
	}

	configPath := filepath.Join(prf, ".config", "nest")

	createFlag := false
	_, err := os.Stat(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			createFlag = true
		} else {
			log.Fatalf("config.get: %v", err.Error())
		}
	}
	if createFlag {
		err := os.MkdirAll(configPath, 0666) //perm bits do nothing windows
		if err != nil {
			log.Fatalf("config.getPath: %v", err.Error())
		}
	}
	return configPath
}

func initialConfig() NestConfig {
	return NestConfig{
		Port:           41595,
		Host:           "127.0.0.1",
		AddFromPath:    "/api/item/addFromPath",
		GetLibraryInfo: "/api/library/info",
		GetThumbnail:   "/api/item/thumbnail",
	}
}

// load once during startup.
// [ ] - validate config
func GetConfig() NestConfig {
	//fmt.Printf("path", GetConfigPath())
	a := filepath.Join(GetConfigPath(), "config.json")
	cfg, err := os.ReadFile(a)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("[INFO] getconfig: creating new config file at", a)
			MustNewConfig()
			return initialConfig()
		} else {
			log.Fatalf("getconfig: error reading file err=%s", err)
		}
	}
	var v NestConfig
	err = json.Unmarshal(cfg, &v)
	if err != nil {
		log.Fatalf("getconfig: error marshalling config err=%s", err)
	}

	v.ApiKey = os.Getenv("eagle_api_key")
	return v
}

// Ptr is a type constraint for pointers to any type.
type Ptr[T any] interface{ *T }

// populates v with contents of <configdir>/filename.json
func PopulateJson[T any, P Ptr[T]](p string, v P) (int, error) {
	cfg := filepath.Join(GetPath(), p)
	bytes, err := os.ReadFile(cfg)
	if err != nil {
		return -1, fmt.Errorf("populateJson: err=%w", err)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return len(bytes), fmt.Errorf("populateJson: unmarshalling json err=%w", err)
	}
	return len(bytes), nil
}

//func PopulateLibraries(e Eagle)

// creates config path at
// `~/.config/nest/config.json
// and creates necessary elements
// or panics
func MustNewConfig() string {
	cfgDir := GetConfigPath()
	jsonPath := filepath.Join(cfgDir, "config.json")

	// validation check :
	//
	// CreateJson
	func() {
		mustCreateFile := false
		_, err := os.Stat(jsonPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				mustCreateFile = true
			}
		}
		if !mustCreateFile {
			return
		}
		handle, err := os.Create(jsonPath)
		if err != nil {
			log.Fatalf("config.new: error creating config.json: %s", err.Error())
		}
		defer handle.Close()

		bytes, err := json.Marshal(initialConfig())
		if err != nil {
			log.Fatalf("config.new: error marshalling defualt config err=%s", err.Error())
		}

		_, err = handle.Write(bytes)
		if err != nil {
			log.Fatalf("config.new: error writing to config.json err=%s", err.Error())
		}
	}()

	return cfgDir
}
