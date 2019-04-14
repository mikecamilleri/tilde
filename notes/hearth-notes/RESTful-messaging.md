# Core & RESTful API Design

## Terms

- device: A single physical device in the home such as a lightswitch or thermostat. Devices may communicate directly with the core using HTTP/REST or via a gateway.
- gateway (hub/controller/interface): A software and hardware device that communicates with the core using HTTP/REST and devices using their specific networking technology. E.g. Z-Wave.
- core: The central piece of software that implements a REST API usedused to communicate with devices and gatways

## Version

During development, version will be 0, first release or stable API will be version 1. Bugfix releases will be given decimal version numbers

## Authentication

Uses HTTP basic auth. Herein "gateway" refers both to true gateways and devices that act as their own gateways.

- for communication initiated by gateways to the core, the username is the gateway's unique name and the password is a random string. 
- for communication initiated by the core to gateways, the username and password are arbitrary. These may be defined by the gateway during gateway registration. The gateway may change this password.
- for communication initiated by humans to the core, the username and password are arbitrary and set at account setup. The user may change this password.

## Authorization

Gateways, devices that communicate directly with the core, and humans should only be given the minimal access necessary to do their job. 

## General Security Considerations

- In order to avid having to implement HTTPS, which would be problematic on low-power hardware such as Arduino, the network that the core uses to communicate with gateways and devices ("device network") should be separate from the general internet connected home network. Ethernet jacks should be secured and WiFi should be encrypted using appropriate encryption. WiFi SSID should not be broadcast. 

- Communication by humans to the core that occurs over the interent should be encrypted with HTTPS, VPN, or similar secure technology.

- Gateways are responsible for secure communication with devices connected to them. The IoT and Home Automation space has a lot of insecure technologies. Products should be selected that allow for reasonable security given the type of devices connected. 

## HTTP methods and status codes

- HTTP methods GET, POST, PUT, PATCH, DELETE, and OPTIONS are supported where appropriate. https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#Request_methods
- Standard HTTP status codes are used. https://en.wikipedia.org/wiki/List_of_HTTP_status_codes

## URL structure

`api-root/version-number/endpoint`

## Gateway and Device Registration

Gateways must be registered by a human. 

Request:
POST `api-root/0/gateway` HTTP/1.1
Content-Type: application/json
Authorization: Basic the-users-basic-auth-string
{
    "gateways": [
        {
            "name": "unique-gateway-name",
            "url": "http://193.168.1.100:80/rest/",
            "username": "arbitrary",
            "password": "a-random-string"
        }
    ]
}

Response:
200 OK
{
    "gateways": [
        {
            "name": "unique-gateway-name",
            "username": "unique-gateway-name",
            "password": "a-random-string",
        }
    ]
}
