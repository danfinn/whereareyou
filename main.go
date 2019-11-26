package main

import (
	"flag"
	"fmt"
	"github.com/qioalice/ipstack"
	"log"
	"net"
	"net/http"
)

var ipstackKey = "da9aa3691dc70204247e74e750447d34"
var hereMapID = "JGc5IdagIJa1autqH5Ns"
var hereMapCode = "Wk9lyx1sVuBDB-cjU6MuJQ"

type point struct {
	lat, long, city, country string
}

// Takes an IP or a hostname and outputs a map jpg.  That IP/hostname could come in the mapip query param
// or if not specified it will pull it from the request header
func mapHost(w http.ResponseWriter, r *http.Request) {
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "userip: %q is not IP:port format", r.RemoteAddr)
	}
	// Check if host query param was set, if it is we will use that
	query := r.URL.Query()
	host, present := query["map"]
	if present && len(host) > 0 {
		value := host[0]
		if net.ParseIP(value) == nil {
			// Must be a hostname, lets resolve it to an IP
			hostIP, err := net.LookupHost(value)
			if err == nil {
				ip_addr := net.ParseIP(hostIP[0])
				ip = ip_addr.String()
			} else {
				fmt.Fprint(w, "<p>Unable to parse IP or lookup hostname</p>")
			}
		} else {
			ip_addr := net.ParseIP(value)
			ip = ip_addr.String()
		}
	}

	forward := r.Header.Get("X-Forwarded-For")
	if forward == "" {
		forward = "unset"
	}

	fmt.Fprint(w, "<html>")
	fmt.Fprint(w, "<head>")
	fmt.Fprint(w, "<link rel='stylesheet' type='text/css' href='static/styles.css' />")
	fmt.Fprint(w, "</head>")
	fmt.Fprint(w, "<title>WhereAmI</title>")
	fmt.Fprint(w, "<h1>Where you at?</h1>")
	fmt.Fprint(w, "<div>")
	fmt.Fprintf(w, "<p>IP: %s</p>", ip)
	fmt.Fprintf(w, "<p>port:  %s</p>", port)
	fmt.Fprintf(w, "<p>X-Forwarded-For: %s</p>", forward)
	fmt.Fprint(w, "</div")

	if ip == "::1" || privateIP(ip) {
		fmt.Fprint(w, "<p>Unable to map private IP</p>")
	} else {
		ip_geo := geoIP(ip)
		fmt.Fprint(w, "<p>")
		fmt.Fprintf(w, "<img src=https://image.maps.api.here.com/mia/1.6/mapview?app_id=%s&app_code=%s&i&lat=%s&lon=%s&h=512&w=512&vt=0&z=14</img>",
			hereMapID, hereMapCode, ip_geo.lat, ip_geo.long)
		fmt.Fprint(w, "</p>")
	}
}

// Takes an IP and gets it's Lattitude and Longitude
func geoIP(ip string) point {
	cli, err := ipstack.New(
		ipstack.ParamToken(ipstackKey),
		ipstack.ParamDisableFirstMeCall())
	res, err := cli.IP(ip)
	if err != nil {
		fmt.Println("error getting IP info")
	}
	place := point{
		lat:     fmt.Sprintf("%f", res.Latitide),
		long:    fmt.Sprintf("%f", res.Longitude),
		city:    res.City,
		country: res.CountryName,
	}
	return place
}

// Checks if an IP is a RFC 1918 private IP
func privateIP(ip string) bool {
	private := false
	IP := net.ParseIP(ip)
	if IP == nil {
		fmt.Println("func privateIP unable to parse IP")
	} else {
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
		private = private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
	}
	return private
}

func main() {

	// Get user flags
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", mapHost)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
