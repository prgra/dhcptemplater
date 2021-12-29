package main

import (
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
	app, err := dhcp.NewApp(conf)
	if err != nil {
		panic(err)
	}
	_, err = app.GetDHCP()
	if err != nil {
		panic(err)
	}
}
