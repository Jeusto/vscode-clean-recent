package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

type DbRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DbRecentlyOpened struct {
	Entries []DbEntry `json:"entries"`
}

type DbEntry struct {
	Workspace       *DbWorkspace `json:"workspace,omitempty"`
	FolderUri       string       `json:"folderUri,omitempty"`
	FileUri         string       `json:"fileUri,omitempty"`
	Label           string       `json:"label,omitempty"`
	RemoteAuthority string       `json:"remoteAuthority,omitempty"`
}

type DbWorkspace struct {
	Id         string `json:"id"`
	ConfigPath string `json:"configPath"`
}

func main() {
	db, err := sql.Open("sqlite3", getDBPath())
	panic_if_error(err)
	defer db.Close()

	recentlyOpened := readHistory()
	filteredRecentlyOpened := removeNonExistantEntries(recentlyOpened)

	fmt.Println(len(recentlyOpened.Entries), "entries initially")
	fmt.Println(len(filteredRecentlyOpened.Entries), "entries after cleaning")

	filteredRecentlyOpenedJson, err := json.Marshal(recentlyOpened)
	panic_if_error(err)

	_, err = db.Exec("UPDATE ItemTable SET value = ? WHERE key = 'history.recentlyOpenedPathsList'", string(filteredRecentlyOpenedJson))
	panic_if_error(err)
}

func getDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home + "/.config/Code/User/globalStorage/state.vscdb"
}

func readHistory() DbRecentlyOpened {
	path := getDBPath()
	var recentlyOpened DbRecentlyOpened

	db, err := sql.Open("sqlite3", path)
	panic_if_error(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ItemTable WHERE key = 'history.recentlyOpenedPathsList'")
	panic_if_error(err)

	for rows.Next() {
		var key string
		var value string

		err = rows.Scan(&key, &value)
		panic_if_error(err)

		err = json.Unmarshal([]byte(value), &recentlyOpened)
		panic_if_error(err)
	}

	return recentlyOpened
}

func removeNonExistantEntries(recentlyOpened DbRecentlyOpened) DbRecentlyOpened {
	var filteredRecentlyOpened DbRecentlyOpened

	for i, entry := range recentlyOpened.Entries {
		var uri string

		if entry.FolderUri != "" {
			uri = entry.FolderUri
		} else if entry.FileUri != "" {
			uri = entry.FileUri
		} else {
			uri = entry.Workspace.ConfigPath
		}

		uri = strings.TrimPrefix(uri, "file://")

		_, err := os.Stat(uri)
		if !os.IsNotExist(err) && uri != "" {
			filteredRecentlyOpened.Entries = append(filteredRecentlyOpened.Entries, recentlyOpened.Entries[i])
		}
	}

	return filteredRecentlyOpened
}

func panic_if_error(err error) {
	if err != nil {
		panic(err)
	}
}
