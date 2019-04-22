# General Requirements

## Components/Objects

- **core**: The central service that _gateways_ communicate with.
- **gateway** _Gateways_ serve as an intermediaries between the _core_ and _devices_ using their specific networking technology (E.g. Z-Wave).
- **device**: A single physical device in the home such as a light switch or thermostat, or a virtual "device" such as an integration with a third party weather service. A _device_ is connected to a _gateway_. A _device_ may serve as its own _gateway_. 
- **feature**: A single functionality implemented by a _device_. Examples: A typical household light switch has a single _feature_, "switch" which can be on or off. A thermostat has several _features_ such as "mode" (heat/cool/off), "temperature," and "set-temperature".

## Examples

### Z-Wave

[Z-Wave](https://www.z-wave.com) is a very common protocol used by home automation products such as light switches, thermostats, garge door openers, etc. A Z-Wave network fits into the above model nicely as Z-Wave has concepts similar to _gateways_, _devices_, and _features_. 

### Weather and Air Quality

Weather and air quality reports are avialable from various web APIs. A _geteway_ providing weather and air quality data might map locations on the planet to _devices_ and specific data or reports to _features_.

### Twitter

Twitter may be of interest in home automation because, among other uses, some government agencies use it to communicate emrgency alerts. A Twitter _gateway_ may map #hashtags or @users to _devices_ which may have only one feature, "Tweet."

## Core-Gateway Communication

- Core MUST have an authentication system in place for security.
- Gateways MUST be able to send device status and other updatess to Core.
- Gateways SHOULD be able to accept messages from Gateways.
- Gateways SHOULD allow their associated Core to configure them amd their devices. 

### Gateway & Core HTTP APIs

The HTTP APIs used for core-gateway communication are not RESTful. This has been done to make ease implementation on gateways and to facilitate the use of transports other than HTTP in the future. 

### Pairing

- `Gateways` MUST have a way of communicating a pairing code to the user.
- (Service discovery will hopefully be able to eliminate this, but for now services MUST communicate a URL.)
- For now HTTP Basic Auth is used. Something better will be implemented later

The _core_ sends a `POST` request to the _gateway_ with the pairing code as the basic auth password and a message body containing credentials and the API url:

{
    "pair": {
        "url": <gateway base url>,
        "username": <unique gateway id>,
        "password": <random>
    }
}

The _gateway_ sends a `POST` request to the _core_:

{
    "pair": {
        "url": <core base url>,
        "username": <anything>,
        "password": <random>
    }
}

## Initialization

On successful pairing, the _gateway_ makes one or more `POST` requests that describe itself: 

{
    "gateway": {
        "make": "com.mikecamilleri",
        "model": "Z-Wave",
        "description": "this gateway integrates with Z-Wave devices",
        "configuration": {},
        "status": {},
        "commands": [
            {
                "id": "restart"
                "name": "Restart"
                "description": "restart the gateway"
                "type": "bool"
            }
        ],
    },
    devices: [
        {
            "id": "",
            "name": "",
            "description": "simple light switch",
            "make": "",
            "model": "",
            "serialNumber": "",
            "status": {
                "connected": true,
            }
            "features": [
                {
                    "type": "binary-switch",
                    "id": "switch",
                    "name": "Switch",
                    "description": "switch",
                    "value": {
                        "reading": false,
                        "setting": false,
                        "type": "bool",
                    }
                }
            ],
        },
        {
            "id": "",
            "name": "",
            "description": "thermostat",
            "make": "",
            "model": "",
            "serialNumber": "",
            "status": {
                "connected": true,
            }
            "features": [
                {
                    "type": "thermometer",
                    "id": "thermometer",
                    "name": "Thermometer",
                    "description": "ambient temperature",
                    "value": {
                        "reading": 75.5,
                        "type": "float",
                        "unit": "Fahrenheit"
                    }
                },
                {
                    "type": "temperature-setting",
                    "id": "temperature-setting",
                    "name": "Temperature Setting",
                    "description": "temperature setting",
                    "value": {
                        "reading": 72,
                        "setting": 72,
                        "type": "integer",
                        "range": [50,100],
                        "unit": "Fahrenheit"
                    }
                },
                {
                    "type": "thermostat-mode-setting",
                    "id": "thermostat-mode-setting",
                    "name": "Mode",
                    "description": "thermostat-mode",
                    "value": {
                        "reading": "auto",
                        "setting": "auto",
                        "type": "string",
                        "options": ["auto", "heat", "cool", "off"]
                    }
                },
                {
                    "type": "thermostat-mode-operation",
                    "id": "thermostat-mode-operation",
                    "name": "",
                    "description": "",
                    "value": {
                        "reading": "cooling",
                        "type": "string",
                        "options": ["heating", "cooling", "off"]
                    }
                }
            ],
        },
    ]
}

# Discourse

`Gateways` update the status of their devices as it changes.

{
    devices: [
        {
            "id": "",
            "features": [
                {
                    "id": "thermostat-mode-operation",
                    "reading": "cooling",
                }
            ]
        }
    ]
}

The `core` changes device feature settings where necessary. 

{
    devices: [
        {
            "id": "",
            "features": [
                {
                    "id": "switch",
                    "setting": true,
                }
            ]
        }
    ]
}

`Gateways` may also update `device` "settings" in cases where the setting is changed by measn ohter than the core. For, example, a light swithch is manually switched or a thermostat is manually adjusted. There are cases in which "setting" and "reading" may not be the same. Such cases include a slow device network or devices that, by their nature, take time to change their state. For example, a garage door that takes some time to roll up or down. 




