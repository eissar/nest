package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type LibrariesConfig struct {
	AutoLoad bool     `json:"autoLoad"`
	Paths    []string `json:"paths"`
}

type NestConfig struct {
	Libraries   LibrariesConfig `json:"libraries"`
	Directories []string        `json:"directories"`
	Port        int             `json:"port"`
	Host        string          `json:"host"`
	ApiKey      string          `json:"apiKey,omitempty"`
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

// param i initialLibraries
func initialConfig(libraryPaths []string) NestConfig {
	return NestConfig{
		Port: 41595,
		Host: "127.0.0.1",
		Libraries: LibrariesConfig{
			AutoLoad: true,
			Paths:    libraryPaths,
		},
		Directories: []string{},
	}
}

// load once during startup.
// [ ] - validate config
func GetConfigV0() NestConfig {
	//fmt.Printf("path", GetConfigPath())
	a := filepath.Join(GetConfigPath(), "config.json")
	cfg, err := os.ReadFile(a)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("[INFO] getconfig: creating new config file at", a)
			MustNewConfig()
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

func tryReadConfig(a string) (NestConfig, error) {
	var v NestConfig
	bytes, err := os.ReadFile(a)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return v, fmt.Errorf("[WARN] readconfig: config path=%s does not exist. error=%w", a, err)
		} else {
			return v, err
		}
	}
	err = json.Unmarshal(bytes, &v)
	return v, nil
}

func GetConfig() NestConfig {
	cfg_pth := filepath.Join(GetConfigPath(), "config.json")
	if cfg, err := tryReadConfig(cfg_pth); err == nil {
		cfg.ApiKey = os.Getenv("eagle_api_key")
		return cfg
	} else if errors.Is(err, os.ErrNotExist) {
		cfg_pth, cfg := MustNewConfig()
		fmt.Printf("[INFO] getconfig: new config file created at %s\n", cfg_pth)
		cfg.ApiKey = os.Getenv("eagle_api_key")
		return cfg
	} else {
		fmt.Printf("unknown error:")
		panic(err)
	}
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
func MustNewConfig() (string, NestConfig) {
	var out NestConfig
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

		out = initialConfig(getRecentLibraries())
		bytes, err := json.Marshal(out)
		if err != nil {
			log.Fatalf("config.new: error marshalling defualt config err=%s", err.Error())
		}

		_, err = handle.Write(bytes)
		if err != nil {
			log.Fatalf("config.new: error writing to config.json err=%s", err.Error())
		}
	}()

	return cfgDir, out
}

func filterLibraries(libraries []string) []string {
	filteredibraries := make([]string, 0, len(libraries))

	for _, path := range libraries {
		if _, err := os.Stat(path); err == nil {
			filteredibraries = append(filteredibraries, path)
		}
	}

	return filteredibraries
}

func getRecentLibraries() []string {
	resp, err := http.Get("http://localhost:41595/api/library/history")
	if err != nil {
		log.Fatalf("newconfig: could not retrieve recent libraries from eagle app. is the app open?")
	}

	if resp.StatusCode != 200 {
		log.Fatalf("newconfig: could not retrieve recent libraries from eagle app. is the app open?")
	}

	var bytes []byte
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("newconfig: could not retrieve recent libraries from eagle app. err=%s", err.Error())
	}

	var responseData struct {
		Data   []string `json:"data"`
		Status string   `json:"status"`
	}
	err = json.Unmarshal(bytes, &responseData)
	if err != nil {
		log.Fatalf("newconfig: could not retrieve recent libraries from eagle app. err=%s", err.Error())
	}
	if responseData.Status != "success" {
		log.Fatalf("newconfig: could not retrieve recent libraries from eagle app (response object from eagle was not success.)")
	}
	if len(responseData.Data) < 1 {
		log.Fatalf("newconfig: could not retrieve recent libraries from eagle app (response object had missing or malformed key `data`.)")
	}

	return filterLibraries(responseData.Data)
}

// also populates libraries with defaults from eagle.
// creates config path at
// `~/.config/nest/config.json
// and creates necessary elements
// or panics
func MustNewConfigAndPopulateLibs() (string, NestConfig) {
	var out NestConfig
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

		out = initialConfig(getRecentLibraries())
		bytes, err := json.Marshal(out)
		if err != nil {
			log.Fatalf("config.new: error marshalling defualt config err=%s", err.Error())
		}

		_, err = handle.Write(bytes)
		if err != nil {
			log.Fatalf("config.new: error writing to config.json err=%s", err.Error())
		}
	}()

	return cfgDir, out
}
