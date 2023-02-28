package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/getlantern/systray"
	"github.com/pelletier/go-toml"
)

type Settings struct {
	SteamPath string
}

type Account struct {
	Username    string
	Description string

	menuItem *systray.MenuItem
}

func getSettings() (*Settings, []*Account, error) {
	path := filepath.Join(applicationDir, "settings.toml")

	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("failed to load file: %w", err)
		}

		file, err = os.Create(path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create file: %w", err)
		}

		_, err = file.WriteString(fmt.Sprintf(`[settings]
steampath = "%s"
[accounts]
username1 = "description..."
username2 = "another description"
username3 = "one more"`, getSteamInstallPath()))
		if err != nil {
			return nil, nil, fmt.Errorf("failed writing to file: %w", err)
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			return nil, nil, fmt.Errorf("failed seeking back: %w", err)
		}
	}
	defer file.Close()

	tree, err := toml.LoadReader(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing file: %w", err)
	}

	settingsTree, ok := tree.Get("settings").(*toml.Tree)
	if !ok {
		return nil, nil, fmt.Errorf("failed parsing file, missing settings tree")
	}

	accountsTree, ok := tree.Get("accounts").(*toml.Tree)
	if !ok {
		return nil, nil, fmt.Errorf("failed parsing file, missing accounts tree")
	}

	settings := &Settings{
		SteamPath: settingsTree.Get("steampath").(string),
	}

	var accounts []*Account
	for key, value := range accountsTree.ToMap() {
		account := &Account{Username: key}

		switch v := value.(type) {
		case string:
			account.Description = v
		case int64:
			account.Description = strconv.FormatInt(v, 10)
		}

		accounts = append(accounts, account)
	}

	return settings, accounts, nil
}
