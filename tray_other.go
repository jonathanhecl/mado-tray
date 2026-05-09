//go:build !darwin

package main

func initTray(app *App) {}

func showTrayIcon() {}

func hideTrayIcon() {}

func setTrayLocale(locale string) {}
