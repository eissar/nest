package databaseutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

/*
mtime key value

key: {mtime:(...),value:(...)}

e.g., songData:
"songData": {"mtime":1257894000,"value":{"song":"Beyond the Clouds","channel":"You'll Never Get to Heaven","duration":"249"}}

	func updateDatabase(key string,value string){
		mtime:= time.Now
		...
	}
*/
func getDatabaseFilePath() (string, error) {
	dbPath := "database.json"
	pth, err := filepath.Abs(dbPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
		} else {
			return pth, fmt.Errorf("error resolving database path (%s) : %s", dbPath, err)
		}
	}
	return pth, nil
}

type Database struct {
	Data map[string]interface{} `json:"data"`
	Path string                 `json:"path"`
}

func (db Database) Save() error {
	// get json data as bytes.Buffer
	jsonData, err := (func() (bytes.Buffer, error) {
		var prettyData bytes.Buffer

		if db.Data == nil {
			// populate with default data.
			db.Data = map[string]interface{}{
				"key": "value",
			}
		}
		data, err := json.Marshal(db.Data)
		if err != nil {
			return prettyData, fmt.Errorf("error marshalling data as json: %s", err)
		}

		err = json.Indent(&prettyData, data, "", "    ")
		if err != nil {
			return prettyData, fmt.Errorf("error prettifying json data (really?) %s", err)
		}
		return prettyData, nil
	})()
	if err != nil {
		return err
	}

	err = os.WriteFile(db.Path, jsonData.Bytes(), 0644) // 0666 The perm argument is irrelevant on Windows
	if err != nil {
		return fmt.Errorf("error while writing database to file: %s", err)
	}

	return nil
}
func (db Database) Update(key string, value string) error {
	// use https://pkg.go.dev/github.com/gofrs/flock#Flock.Lock
	return nil
}

func GetDatabase() (Database, error) {
	var db Database //map[string]interface{}
	// pth, err := getDatabaseFilePath()
	// if err != nil {
	// 	//if errors.Is(err, os.ErrNotExist) {
	// 	return db, fmt.Errorf("error getting db filepath %s:", err)
	// }
	db.Path = "database.json"
	// create the database if it does not exist.
	db.Save()

	bytes, err := os.ReadFile(db.Path)
	if err != nil {
		return db, fmt.Errorf("error reading bytes from database path (%s) : %s", db.Path, err)
	}
	err = json.Unmarshal(bytes, &db.Data)
	if err != nil {
		return db, fmt.Errorf("error unmarshalling db bytes from database : %s", err)
	}
	db.Data["test"] = "Hey"
	db.Save()

	return db, nil
}
