package main

import (
	"fmt"
	"net"
	"os"
	"text/template"
)

type Data struct {
	Nic      string
	Priority string
        Router   string
	Vip      string
}

var kaconfig = `! Configuration File for keepalived
global_defs {
   vrrp_version 2
}
vrrp_instance VI_{{ .Router }} {
    state BACKUP
    interface {{ .Nic }}
    virtual_router_id {{ .Router }}
    priority {{ .Priority }}
    advert_int 1
    nopreempt
    virtual_ipaddress {
        {{ .Vip }}/32
    }
}`

func findNicForVip(vip string) (string, string, string, error) {
	netvip := net.ParseIP(vip)
	if netvip == nil {
		return "", "", "", fmt.Errorf("%s is not a valid ip", vip)
	}
	if netvip.To4() == nil {
		return "", "", "", fmt.Errorf("%s is not an ipv4", vip)
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing interfaces")
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", "", "", fmt.Errorf("error parsing interface %s address", i.Name)
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil && !v.IP.IsLoopback() && v.Contains(netvip) {
					p := (int(v.IP[len(v.IP)-2]) + 1) * (int(v.IP[len(v.IP)-1]) + 1) % 255
                                        r := int(netvip[len(netvip)-1])
					return i.Name, fmt.Sprintf("%v", p), fmt.Sprintf("%v", r), nil
				}
			}
		}
	}
	return "", "", "", fmt.Errorf("No valid interface found for %s", vip)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <vip>\n", os.Args[0])
		os.Exit(1)
	}
	nic, prio, router, err := findNicForVip(os.Args[1])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	data := Data{nic, prio, router, os.Args[1]}
	tmpl, _ := template.New("test").Parse(kaconfig)
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
