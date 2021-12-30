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

type dbnet struct {
	IPID    int    `db:"ip_zone_id"`
	Name    string `db:"name"`
	Net     int64  `db:"net"`
	Mask    int64  `db:"mask"`
	GateWay int64  `db:"gateway"`
}

func (a *App) GetDHCP() (dta []byte, err error) {
	t, err := template.ParseFiles("templates/dhcpd.tmpl")
	if err != nil {
		return []byte(""), err
	}
	hs := DHCP{
		NetName: a.Cfg.NetName,
	}
	var nets []dbnet
	err = a.DB.Select(&nets, "select ip_zone_id, z.name, net, mask, gateway from ip_zones z join ip_zones_detail zd on zd.id = z.id")
	if err != nil {
		return []byte(""), err
	}
	// var gwmap map[string]string
	for i := range nets {
		nt := net.ParseIP(inet_ntoa(uint32(nets[i].Net)))
		if !nt.IsGlobalUnicast() {
			continue
		}
		gw := net.ParseIP(inet_ntoa(uint32(nets[i].GateWay)))
		if !gw.IsGlobalUnicast() {
			continue
		}
		hs.Nets = append(hs.Nets, Net{
			Net:     nt.String(),
			Mask:    inet_ntoa(uint32(nets[i].Mask)),
			GateWay: gw.String(),
			DNSes:   []string{"1.1.1.1", "8.8.8.8"},
		})
	}

	var hosts []dbent
	err = a.DB.Select(&hosts, "select id, ip, mac from ip_groups where mac != ''")
	if err != nil {
		return []byte(""), err
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
