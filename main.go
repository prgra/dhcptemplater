package main

import "github.com/prgra/dhcptemplater/dhcp"

func main() {
	app := dhcp.NewApp(dhcp.Cfg{})
	app.GetTemplates()
}
