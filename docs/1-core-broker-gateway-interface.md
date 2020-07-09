# Core — Message Broker — Gateway Interfaces

## Versioning

The API will be versioned using three integers separated by decimal points. Using the example version `1.2.3`:

1. The first integer represeents an API version the guarantees only backwards compatable changes will be made. I.e. a _core_ implementing API version `1.3` will be able to communicate with Gateways implementing API version `1.0`, `1.1`, `1.2`, or `1.3`.  Before the first stable release, this first digit will be `0`.

2. The second integer is incremented when backwards-compatable features are added.

3. The third integer is incremented when bug fixes are made.

## Initial Connection and Registration Flow

Gateways need a way to gain authorization to connect to the message broker. 

```
User                                 Core                                Gateway
  |                                    |                                     |
1.|<------------------- displays the unique device ID (e.g. serial number) --|
  |                                    |                                     |
2.|-- enters uniq. ID of gateway ----->|                                     |
  |                                    |                                     |
3.|<----- displays one time password --|                                     |
  |                                    |                                     |
4.|-- enters one time password --------------------------------------------->|
  |                                    |                                     |
5.|                                    |<--- makes HTTP registration reqst --|
  |                                    |                                     |
6.|                                    |-- returns credentials, URI's, &c -->|
  |                                    |                                     |
7.|                                    |<--------- connects to msg. borker --|
  |                                    |                                     |
  V                                    V                                     V  
  ```

1. The _user_ reads a unique device ID from the physical _gateway_ device ot it's UI.

2. The _user_ enter the unique ID from the _gateway_ into the _core_. 

3. The _core_ generates and displays a one time password.

4. The _user_ enters the one time password into the _gateway_.

5. The _gateway_ makes a regsitration request via HTTP to the _core_'s registration API.

6. The _core_ authenticates the one time password; returns ccredenentials, URLs, etc.; and deltes the OTP. 

7. The _gateway_ connects to the message broker.

## Routine Date Exchange

_NOTE: This information is out of data and needs to be updated!_

Similar to the connection flow described above, the routine data exchange between the _core_ and _gateways_ is designed such that _gateways_ can be implemented as simply as possible. 

_Gateways_ are responsible for initiating and mantaining a WebSocket connection with the _core_. All data echange occurs over that connection. This allows the _gateways_ to only have to implement an WebSocket client and prevents them from having to advertize themselves using mDNS/DNS-SD.

All communications between the _core_ and a _gateway_ fit the following template:

```json
{
    "gateways": [],
    "devices": [],
    "errors": []
}
```

Only relevent, non-empty elements ("porperties" in JSON lingo) should to be included in each request. A `description` element with an empty value (`""`) would imply that it should be erased. _Gateways_ are not responsible for maintaing any information about the state of the core across connection sessions, other than their login credentials. Every time a `gateway` connects to the core, it should send complete information about itself and its _devices_ to the core as if it had never connected before. The core will match `gateways`, `devices`, and features across sessions by their `id` values. 

It may seem odd that `gateways` holds an array. This is done monstly for consistancy with the other root-level elements. In some cases, however, it may make sense for multiple software _gateways_ to share a single WebSocket connection and thus a single set of authentication credentials. For example, a hardware device that is able to connect to ZigBee and Z-Wave _devices_ may want to define those protocols as separate _gateways_ but share a single registration flow, share a single set of credentials, and use a single WebSocket connection. 

### `gateways`

The `gateways` array contains objects in the following format:

```json
{
    "hardwareId": "a-serial-number-or-similar",
    "id": "a-universally-unique-id-self-assigned-by-the-gateway",
    "name": "A Friendly Name",
    "description": "A frindly decription of what the gateway does.",
    "connected": true,
    "features": []
}
```

Before a _gateway_ disconnects from the core, if possible, it is polite for it to set its `connected` status to `false`.

### `devices`

The `devices` array contains objects in the following format:

```json
{
    "hardwareId": "a-serial-number-or-similar",
    "gatewayId": "the-id-of-the-gateway-this-device-is-connected-to",
    "id": "a-unique-id-within-the-gateway",
    "name": "A Friendly Name",
    "description": "A friendly decription of what this device does.",
    "features": []
}
```

### `features`

`gateways` and `devices` have `features`.

```json
{
    "id": "a-unique-id-within-the-device",
    "name": "A Friendly Name",
    "description": "A friendly decription of what this feature does.",
    "definition": {
        "standard": "the-name-of-the-standard-feature-type",
        "unit": "",
    },
    "setting": {
        "value": ,
    },
    "readings": [
        {
            "value": ,
            "time": "2019-04-22T21:03:49+00:00"
        }
    ],
}
```

The above template represents a standardized _feature_, if the _feature_ were a custom feature, the `definition` object would contain:

```json
    "definition": {
        "standard": "custom",
        "name": "A Friendly Name for the Feature Type",
        "description": "A friendly description of the feature type.",
        "units": [""],
        "jsonType": "",
        "min": ,
        "max": ,
        "options": [""],
        "settable": true
    }
```

Standard _features_ will be defined in a similar way and will be a part of the standard for each API version. Using standard features enables actions such as "turn off all the lights". `Readings` is an array so that _gateways_ may provide historic readings if they or the device become disconnected or if the device is designed to only report intermittently. 

A few quick notes:

- `jsonType` may be: "boolean", "number," or "string."
- `units` are something like "fahrenheit," "percent," or "volt". I am not sure yet whether I want the whole name spelled out or whether a standard abbreviations will be used. In either case these will need to be defined somewhere -- probably along with the the standared _features_.
- `min` and `max` are used to set bounds for `value` when `jsonType` is "number."
- `options` is a list of valid `value` strings when `jsonType` is "string."
- Some _features_ may not be settable. A device providing a weather report, for example, won't have settable _features_.
- Some standard _features_ may be manditory. One such feature may be "connected".

Below are some example hypothetical `device` objects:

```json
{
    "devices": [
        {
            "hardwareId": "a-serial-number-or-similar",
            "gatewayId": "the-id-of-the-gateway-this-device-is-connected-to",
            "id": "a-unique-id-within-the-gateway",
            "name": "Porch Light",
            "description": "",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Switch",
                    "description": "",
                    "definition": {
                        "standard": "switch-binary"
                    },
                    "readings": [
                        {
                            "value": "on",
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": "on"
                    }
                }
            ]
        },
        {
            "hardwareId": "a-serial-number-or-similar",
            "gatewayId": "the-id-of-the-gateway-this-device-is-connected-to",
            "id": "a-unique-id-within-the-gateway",
            "description": "",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Ambient Temperature",
                    "description": "",
                    "definiton": {
                        "standard": "thermometer",
                        "unit": "fahrenheit"
                    },
                    "readings": [
                        {
                            "value": 65,
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ]
                },
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Temperature Setting",
                    "description": "",
                    "definiton": {
                        "standard": "thermostat-temperature-setting",
                        "unit": "fahrenheit"
                    },
                    "readings": [
                        {
                            "value": 70,
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": 70
                    }
                },
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Mode",
                    "description": "",
                    "definiton": {
                        "standard": "thermostat-mode"
                    },
                    "readings": [
                        {
                            "value": "auto",
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": "auto"
                    }
                }
            ]
        },
        {
            "hardwareId": "a-serial-number-or-similar",
            "gatewayId": "the-id-of-the-gateway-this-device-is-connected-to",
            "id": "a-unique-id-within-the-gateway",
            "name": "Party Lights",
            "description": "Get the party started with this custom device that controls party lights (whatever those are)",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "mode",
                    "description": "",
                    "definition": {
                        "standard": "custom",
                        "jsonType": "string",
                        "options": ["flash", "fade", "wave"],
                        "settable": true
                    },
                    "readings": [
                        {
                            "value": "flash",
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": "flash"
                    }
                },
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "intensity",
                    "description": "",
                    "standard": "custom",
                    "definition": {
                        "jsonType": "number",
                        "min": 0,
                        "max": 9,
                        "settable": true
                    },
                    "readings": [
                        {
                            "value": 5,
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": 5
                    }
                },
            ]
        },
    ]
}
```

## Errors

TODO ...
