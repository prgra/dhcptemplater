package dhcp

import (
	"fmt"
	"net"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Cfg struct {
	DBURL   string `toml:"mysql"`
	NetName string `toml:"netname"`
}

type App struct {
	DB  *sqlx.DB
	Cfg Cfg
}

func NewApp(c Cfg) (a App, err error) {
	a.Cfg = c
	a.DB, err = sqlx.Open("mysql", c.DBURL)
	return a, err
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

type dbent struct {
	ID  int    `db:"id"`
	IP  int64  `db:"ip"`
	Mac string `db:"mac"`
}

func (a *App) GetDHCP() (dta []byte, err error) {
	t, err := template.ParseFiles("templates/dhcpd.tmpl")
	if err != nil {
		return []byte(""), err
	}
	var hosts []dbent
	err = a.DB.Select(&hosts, "select id, ip, mac from ip_groups where mac != ''")
	if err != nil {
		return []byte(""), err
	}
	hs := DHCP{
		NetName: a.Cfg.NetName,
	}
	for i := range hosts {
		// mysql inet_ntoa don't work
		sip := inet_ntoa(uint32(hosts[i].IP))
		ip := net.ParseIP(sip)
		if !ip.IsGlobalUnicast() {
			continue
		}
		mac, err := net.ParseMAC(hosts[i].Mac)
		if err != nil {
			continue
		}
		hs.Hosts = append(hs.Hosts,
			Host{
				IP:   ip.String(),
				Mac:  mac.String(),
				Name: fmt.Sprintf("id_%d", hosts[i].ID),
			})
	}
	t.Execute(os.Stdout, hs)
	return
}

func inet_ntoa(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
