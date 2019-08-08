# Gateway API

The Gateway API is a WebSocket/JSON API. The Gateway API server is implemented in the _core_ and connected to by clients implemented in _gateways_. 

## Initial Connection and Registration

A procedure will be established by which new _gateways_ will be able to securely connect to the _core_ and register themselves. All further communication will also be secured and encrypted. 

Z-Wave Gateway registration object:
```
{
    "gateways": [
        {
            "name": "Z-Wave Gateway",
            "description": "A Z-Wave Gateway",
            "version": "0",
            "hardwareId": "a-serial-number-or-similar",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Core Connection",
                    "description": "",
                    "standard": {
                        "type": "connection",
                        "version": "0"
                    },
                    "readings": [
                        {
                            "value": "on",
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ]
                },
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Restart",
                    "description": "Restart the gateway",
                    "standard": {
                        "type": "command",
                        "version": "0"
                    },
                    "readings": [
                        {
                            "value": true,
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": true
                    }
                }
            ]
        }
    ]
}
```

_Gateways_ will only ever be registered individually, but an array is used for consistency. _Gateways_ are not aware of the unique ID or friendly name assigned to them in the _core_ which is why none are present here. 

## Routine Connection

A _gateway_ should maintain an open WebSocket connection to the core at all times. In the event that the _gateway_ must disconnect from the _core_ such as to reboot, it should update its state with the core to indicate such before doing so. 

## Device Registration

At any time after initial connection and registration, a _gateway_ may register its _devices_ with the _core_.

Switch, thermostat, and custom registration objects:
```
{
    "devices": [
        {
            "id": "a-unique-id-within-the-gateway",
            "name": "Porch Light",
            "description": "",
            "version": "0",
            "hardwareId": "a-serial-number-or-similar",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Switch",
                    "description": "",
                    "standard": {
                        "type": "switch-binary",
                        "version": "0"
                    },
                    "readings": [
                        {
                            "value": "on",
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ]
                    "setting": {
                        "value": "on"
                    }
                }
            ]
        },
        {
            "id": "a-unique-id-within-the-gateway",
            "name": "Living Room Thermostat",
            "description": "",
            "version": "0",
            "hardwareId": "a-serial-number-or-similar",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Ambient Temperature",
                    "description": "",
                    "standard": {
                        "type": "thermometer",
                        "version": "0"
                    },
                    "readings": [
                        {
                            "value": {
                                "value": 65,
                                "unit": "fahrenheit"
                            }
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ]
                },
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Temperature Setting",
                    "description": "",
                    "standard": {
                        "type": "thermostat-temperature-setting",
                        "version": "0"
                    },
                    "readings": [
                        {
                            "value": {
                                "value": 65,
                                "unit": "fahrenheit"
                            }
                            "time": "2019-04-22T21:03:49+00:00"
                        }
                    ],
                    "setting": {
                        "value": {
                            "value": 70,
                            "unit": "fahrenheit"
                        }
                        "unit": "fahrenheit"
                    }
                },
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "Mode",
                    "description": "",
                    "standard": {
                        "type": "thermostat-mode",
                        "version": "0"
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
            "id": "a-unique-id-within-the-gateway",
            "name": "Party Lights",
            "description": "Get the party started with this custom device that controls party lights (whatever those are)",
            "version": "0",
            "hardwareId": "a-serial-number-or-similar",
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "name": "mode",
                    "description": "",
                    "standard": {
                        "type": "custom",
                        "version": "0"
                    },
                    "definition": {
                        "type": "string",
                        "options": ["flash", "fade", "wave"],
                        "settable": true
                    }
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
                    "standard": {
                        "type": "custom",
                        "version": "0"
                    },
                    "definition": {
                        "type": "integer",
                        "min": 0,
                        "max": 9,
                        "settable": true
                    }
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

TODO: Figure out how to report gateway-device connection status. 
TODO: Units as part of feature?

Features may be custom or standard. Using standard features enables actions such as "turn off all the lights". The definition of standard features is done separately. Likely in separate, liberally licensed, Git repo. `readings` is an array so that _gateways_ may provide historic readings if they or the device become disconnected or if the device is designed to only report intermittently. 

## State Updates

After a _device_ is registered, its _gateway_ should update its state (reading or setting) with the _core_ immediately upon a state change. This is done by sending messages similar to those created during registration, but with only the updated fields. For example, if an already registered light switch is turned off, the following message might be sent to the core:

{
    "devices": [
        {
            "id": "a-unique-id-within-the-gateway"
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "readings": [
                        {
                            "value": "off",
                            "time": "2019-04-23T11:00:00+00:00"
                        }
                    ],
                    "setting": {
                        "value": "off"
                    }
                }
            ]
        }
    ]
}

Note that only the fields which have changed and relevent `id` fields need to be re-sent. 

## Desired State and Configuration Requests

The _core_ may update the desired state and configuration of _gateways_ and _devices_ by sending requests to the _gateway_. To turn on the light above, the core would send the following message to the appropriate gateway:

{
    "devices": [
        {
            "id": "a-unique-id-within-the-gateway"
            "features": [
                {
                    "id": "a-unique-id-within-the-device",
                    "setting": {
                        "value": "on"
                    }
                }
            ]
        }
    ]
}
