package dhcp

import (
	"os"
	"text/template"

	"github.com/jmoiron/sqlx"
)

type Cfg struct {
	DBURL        string `toml:"mysql"`
	TemplatePath string `toml:"templates"`
}

type App struct {
	DB  *sqlx.DB
	Cfg Cfg
}

func NewApp(c Cfg) *App {
	return &App{
		Cfg: c,
	}
}

type DHCP struct {
	NetName string
	Nets    []Net
	Hosts   []Host
}

type Net struct {
	Net     string
	Mask    string
	DNSes   []string
	GateWay string
}

type Host struct {
	Name string
	IP   string
	Mac  string
}

func (a *App) GetTemplates() {
	t, err := template.ParseFiles("templates/dhcpd.tmpl")
	if err != nil {
		panic(err)
	}
	hs := DHCP{
		NetName: "HelloWorld",
		Nets: []Net{
			{
				Net:     "127.0.0.1",
				Mask:    "255.255.255.0",
				DNSes:   []string{"8.8.8.8", "5.5.5.5"},
				GateWay: "127.0.0.2",
			},
			{
				Net:     "127.0.0.2",
				Mask:    "255.255.255.0",
				DNSes:   []string{"8.8.8.8", "5.5.5.5"},
				GateWay: "127.0.0.2",
			},
		},
		Hosts: []Host{
			{
				Name: "asda",
				IP:   "127.0.0.1",
				Mac:  "hhbb",
			},
			{
				Name: "asdaaa",
				IP:   "127.0.0.2",
				Mac:  "hhcc",
			},
		},
	}
	t.Execute(os.Stdout, hs)
}
