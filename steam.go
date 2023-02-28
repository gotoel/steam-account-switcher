package main

import "os/exec"

const DEFAULT_STEAM_PATH = "C:\\Program Files\\Steam"

func startSteam() error {
	return exec.Command("cmd", "/C", "start", "steam://open/main").Start()
}

func stopSteam() error {
	return exec.Command("taskkill", "/F", "/IM", "steam.exe").Run()
}
