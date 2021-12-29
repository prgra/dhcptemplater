package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/prgra/dhcptemplater/dhcp"
)

func main() {
	var conf dhcp.Cfg
	_, err := toml.DecodeFile("dhcp.toml", &conf)
	if err != nil {
		log.Println("can't read config", err)
		os.Exit(1)
	}
	fmt.Println(conf)
	app := dhcp.NewApp(dhcp.Cfg{})
	app.GetTemplates()
}
