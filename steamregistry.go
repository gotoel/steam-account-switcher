package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
)

func switchAccountRegistry(account *Account) error {
	_ = stopSteam()

	reg, err := getAccountRegistry()
	if err != nil {
		return fmt.Errorf("could not get registry: %w", err)
	}

	if err = reg.SetStringValue("AutoLoginUser", account.Username); err != nil {
		return fmt.Errorf("failed to set value in registry: %w", err)
	}
	if err = reg.SetDWordValue("RememberPassword", 1); err != nil {
		return fmt.Errorf("failed to set value in registry: %w", err)
	}

	return startSteam()
}

func getActiveUsername() (string, error) {
	reg, err := getAccountRegistry()
	if err != nil {
		return "", fmt.Errorf("could not get registry: %w", err)
	}

	username, _, err := reg.GetStringValue("AutoLoginUser")
	if err != nil {
		return "", fmt.Errorf("could not get value in registry: %w", err)
	}

	return username, nil
}

func getSteamInstallPath() string {
	installPath := DEFAULT_STEAM_PATH

	reg, err := getSteamRegistry()
	if err != nil {
		return installPath
	}

	installPath, _, err = reg.GetStringValue("InstallPath")
	if err != nil {
		// Try with 64bit registry
		reg, err = getSteamRegistry64bit()
		if err != nil {
			return installPath
		}

		installPath, _, err = reg.GetStringValue("InstallPath")
	}

	return filepath.ToSlash(installPath)
}

func getAccountRegistry() (registry.Key, error) {
	return registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.READ|registry.WRITE)
}

func getSteamRegistry() (registry.Key, error) {
	return registry.OpenKey(registry.LOCAL_MACHINE, `Software\Valve\Steam`, registry.READ)
}

func getSteamRegistry64bit() (registry.Key, error) {
	return registry.OpenKey(registry.LOCAL_MACHINE, `Software\WOW6432Node\Valve\Steam`, registry.READ)
}
