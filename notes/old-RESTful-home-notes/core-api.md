# Core API

The Core API is used by both gateways and users (directly or via other user interfaces) to interact with the core.

## General Considerations

Although not shown in the examples here, all requests and responses will have appropriate [HTTP headers](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields). When necessary, HTTP `Accept` and `Content-type` headers should be `application/json`. 

`PATCH` is used to make partial changes to an entity. `PUT` is used in a very few instances to completely replace an entity. `POST` is used exlusively to create subordinate resources. 

Property names are `camelCase`.

## Authentication

Athentication is handled using HTTP Basic authentication because it is easy to use, widely understood, and includes a separate username and password.

## Errors

Endpoints will return appropriate [HTTP status codes](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes). 

Endpoints may provide additional information about an error in an error response such as below. The error in this example is for illustration purposes only. A list of error codes will be written later.

```
{
    "errors": [
        {
            "code": "10",
            "name": "UNMODIFIABLE_FIELD",
            "message": "The field 'features.thermostat.value' may not be modified."
        },
        {
            "code": "10",
            "name": "UNMODIFIABLE_FIELD",
            "message": "The field 'features.temperature.value' may not be modified."
        }
    ]
}
```

In order to avoid confusion with HTTP error codes (a major HTTP API pet peeve of mine), application specific error codes will be two digits and range from 10-99. Each block of 10 will represent a distinct group of errors. Structuring the error codes in this way facilitates clients easily handling groups of similar errors.

## Endpoints

### Root

- Endpoint: root

```
GET /

{
    "apiVersion": "0.0.0",
    "href": "http://example-core:80/",
    "id": "example-core",
    "name": "Example Core",
    "config": {
        "href": "http://example-core:80/configuration",
        "fields": [
            {
                "id": "latitude",
                "name": "Latitude",
                "type": "float",
                "mutable": true,
                "range": [-90.0, 90.0]
            }
            {
                "id": "latitude",
                "name": "Latitude",
                "type": "float",
                "mutable": true,
                "range": [-180.0, 180.0]
            }
        ],
        "desired": {
            "latitude": 45.512794,
            "longitude": -122.679565
        },
        "current": {
            "latitude": 45.512794,
            "longitude": -122.679565
        },
        "history": {
            "href": "http://example-core:80/configuration/history"
        }
    },
    "commands": {
        "href": "http://example-core:80/commands"
        "fields": [
            {
                "id": "discover",
                "name": "Discover",
                "type": "boolean",
                "mutable": true
            }
        ]
        "desired": {
            "discover": true,
        }
        "history": {
            "href": "http://example-core:80/commands/history"
        }
    },
    "auths": {
        "href": "http://example-core:80/auths"
    },
    "gateways": {
        "href": "http://example-core:80/gateways"
    },
    "devices": {
        "href": "http://example-core:80/devices"
    }
}
```

The root URL may contain a path such as `api` or `rest` so that other HTTP services may operate on the same port. The endpoint paths above should be considered standard.

### Auths

- Endpoint: `auth`

```
Request: GET /auths

Response: 200 OK
{
    "href": "http://example-core:80/auths",
    "auths": [
        {
            "href": "http://example-core:80/auths/mike",
            "id": "mike",
            "name": "Michael Camilleri",
            "roles": ["admin"]
        },
        {
            "href": "http://example-core:80/auths/example-gateway",
            "id": "example-gateway",
            "name": "Example Gateway",
            "roles": ["gateway"]
        }
    ]
}
```

```
Request: POST /auths
{
    "id": "new-auth",
    "name": "New Auth",
    "password": "a-password-string",
    "roles": ["user"]
}

Response: 201 Created
{
    "href": "http://example-core:80/auths/new-auth",
    "id": "new-auth",
    "name": "New Auth",
    "password": "a-password-string",
    "roles": ["user"]
}
```

Auths may have the following roles:
- `admin`: A superuser. May use all available methods on all endpoints.
- `user`: Privilages are appropriate for someone occupying the residence but not needing administrative access.
- `gateway`: Has access to its own status and devices.

### Gateways

- Endpoint: `gateways`

```
Request: GET /gateways

Response: 200 OK
{
    "href": "http://example-core:80/gateways"
    "gateways": [
        {
            "href": "http://example-core:80/gateways/example-gateway",
            "remoteHref": "http://example-gateway:80/",
            "id": "example-gateway",
            "name": "Example Gateway",
            "auth": {
                "remoteHref": "http://example-gateway:80/auth",
                "username": "something",
                "password": "a-random-string"
            },
            "config": {
                "href": "http://example-core:80/gateways/example-gateway/config",
            },
            "commands": {
                "href": "http://example-core:80/gateways/example-gateway/commands"
            },
            "devices": {
                "href": "http://example-core:80/gateways/example-gateway/devices"
            },
        }
    ]
}
```

The `devices` object may not be necessary. 

### Devices

- Endpoint: `devices`

```
Request: GET /devices

Response: 200 OK
{
    "href": "http://example-core:80/devices"
    "devices": [
        {
            "href": "http://example-core:80/devices/kitchen-ceiling-light",
            "id": "kitchen-ceiling-light",
            "name": "Kitchen Ceiling Light",
            "type": "switch",
            "gatewayId": "gateway-id",
            "state": {
                "href": "http://example-core:80/devices/kitchen-ceiling-light/state",
                "fields": [
                    {
                        "id": "connected",
                        "name": "Connected",
                        "type": "boolean",
                        "mutable": true
                    },
                    {
                        "id": "powered",
                        "name": "Powered",
                        "type": "boolean",
                        "mutable": true
                    }
                ],
                "desired": {
                    "connected": true,
                    "powered": true,
                    "updatedTime": "2018-07-29T13:55:27-0700"
                },
                "current": {
                    "connected": true,
                    "powered": true,
                    "updatedTime": "2018-07-29T13:55:27-0700",
                    "collectedTime": "2018-07-29T13:55:25-0700"
                },
                "history": {
                    "href": "http://example-core:80/devices/kitchen-ceiling-light/state/history"
                }

            }
        },
        {
            "href": "http://example-core:80/devices/first-floor-thermostat",
            "id": "first-floor-thermostat",
            "name": "First Floor Thermostat",
            "type": "thermostat",
            "gatewayId": "gateway-id",
            "state": {
                "href": "http://example-core:80/devices/first-floor-thermostat/state",
                "fields": [
                    {
                        "id": "connected",
                        "name": "Connected",
                        "type": "boolean",
                        "mutable": true
                    },
                    {
                        "id": "mode",
                        "name": "Mode",
                        "type": "string",
                        "mutable": true,
                        "options": ["auto", "heat", "cool", "off"]
                    },
                    {
                        "id": "heating",
                        "name": "Heating",
                        "type": "bool",
                        "mutable": false
                    },
                    {
                        "id": "cooling",
                        "name": "Cooling",
                        "type": "bool",
                        "mutable": false
                    },
                    {
                        "id": "temperatureSetting",
                        "name": "Temperature Setting",
                        "type": "float",
                        "mutable": true,
                        "range": [50.0, 100.0]
                    },
                    {
                        "id": "temperatureReading",
                        "name": "Temperature Reading",
                        "type": "float",
                        "mutable": false
                    }
                ],
                "desired": {
                    "connected": true,
                    "mode": "auto",
                    "temperatureSetting": "70",
                    "updatedTime": "2018-07-29T13:55:27-0700"
                },
                "current": {
                    "connected": true,
                    "mode": "auto",
                    "temperatureSetting": "70.0",
                    "temperatureReading": "65.5",
                    "heating": true,
                    "cooling": false,
                    "updatedTime": "2018-07-29T13:55:40-0700",
                    "collectedTime": "2018-07-29T13:55:32-0700"
                }
                "history": {
                    "href": "http://example-core:80/devices/first-floor-thermostat/state/history"
                }
            }
        }
    ]
}
```
