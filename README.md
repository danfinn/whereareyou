# whereareyou
Pinpoints the location of your http request on a map.

## Usage
When run this app listens on port 8080.  A request to / will check your ```Remote_Addr``` header and map your location from that.  
You can also specify a hostname or IP as a query parameter:
```http://host:8080/&map=8.8.8.8```
```http://host:8080/&map=www.google.com"
and it will display the Geolocation information for that IP or hostname.

## Built With

* [IPStack](http://www.ipstack.com) - Geolocation API
* [HERE](https://www.here.com/) - Mapping

