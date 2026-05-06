//go:build !darwin

package main

func initTray(app *App, visible bool) {}

func showTrayIcon() {}

func hideTrayIcon() {}
