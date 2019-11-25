package main

import ("fmt"
	"github.com/qioalice/ipstack"
	"log"
	"net"
	"net/http"
)

var ipstackKey = "da9aa3691dc70204247e74e750447d34"
var hereMapID = "JGc5IdagIJa1autqH5Ns"
var hereMapCode = "Wk9lyx1sVuBDB-cjU6MuJQ"

type point struct {
	lat string
	long string
	city string
	country string

}

// Takes an IP and outputs a map jpg.  That IP could come in the mapip query param 
// or if not specified it will pull it from the request header 
func mapIP(w http.ResponseWriter, r *http.Request) {
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "userip: %q is not IP:port format", r.RemoteAddr)
	}
	// Check if mapip query param was set, if it is we will use that
	query := r.URL.Query()
	mapip, present := query["mapip"]
	if present && len(mapip) > 0 {
		value := mapip[0]
		if net.ParseIP(value) == nil {
			fmt.Fprintf(w, "IP %s not in correct format", value)
			return
		} else {
			ip_addr := net.ParseIP(value)
			ip = ip_addr.String()
		}	
	}

  	forward := r.Header.Get("X-Forwarded-For")
       	if forward == "" { forward = "unset" }

       	fmt.Fprintf(w, "<p>IP: %s</p>", ip)
       	fmt.Fprintf(w, "<p>port:  %s</p>", port)
       	fmt.Fprintf(w, "<p>X-Forwarded-For: %s</p>", forward)

	if ip == "::1" || privateIP(ip) {
		fmt.Fprintf(w, "<p>Unable to map private IP</p>")
	} else {
		ip_geo := geoIP(ip)
		fmt.Fprintf(w, "<img src=https://image.maps.api.here.com/mia/1.6/mapview?app_id=%s&app_code=%s&i&lat=%s&lon=%s&h=512&w=512&vt=0&z=14</img>",
			hereMapID, hereMapCode, ip_geo.lat, ip_geo.long)
	}
}

// Takes an IP and gets it's Lattitude and Longitude
func geoIP(ip string) point {
	cli, err := ipstack.New(
		ipstack.ParamToken(ipstackKey),
		ipstack.ParamDisableFirstMeCall())
	res,err := cli.IP(ip)
	if err != nil {
		fmt.Println("error getting IP info")
	}
	place := point {
		lat: fmt.Sprintf("%f", res.Latitide),
		long: fmt.Sprintf("%f", res.Longitude),
		city: res.City,
		country: res.CountryName,
	}
	return place
}

// Checks if an IP is a RFC 1918 private IP
func privateIP(ip string) (bool) {
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
	http.HandleFunc("/", mapIP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
