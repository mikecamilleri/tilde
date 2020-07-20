# Core — Message Broker — Gateway Interfaces

## Versioning

The API will be versioned using three integers separated by decimal points. Using the example version `1.2.3`:

1. The first integer represents an API version the guarantees only backwards compatible changes will be made. I.e. a _core_ implementing API version `1.3` will be able to communicate with Gateways implementing API version `1.0`, `1.1`, `1.2`, or `1.3`.  Before the first stable release, this first digit will be `0`.

2. The second integer is incremented when backwards-compatible features are added.

3. The third integer is incremented when bug fixes are made.

## Initial Connection and Registration Flow

Gateways need a way to gain authorization to connect to the message broker. 

```
User                             Core/RabbitMQ                           Gateway
  |                                    |                                     |
1.|<------------------ displays the unique gateway ID (e.g. serial number) --|
  |                                    |                                     |
2.|-- enters uniq. ID of gateway ----->|                                     |
  |                                    |                                     |
3.|<----- displays one time password --|                                     |
  |                                    |                                     |
4.|-- enters one time password and url of RabbitMQ ------------------------->|
  |                                    |                                     |
5.|                                    |<-------- makes registration rqst. --|
  |                                    |                                     |
6.|                                    |-- returns credentials etc. -------->|
  |                                    |                                     |
7.|                                    |<------ connects using credentials --|
  |                                    |                                     |
  V                                    V                                     V  
  ```

1. The _user_ reads a unique device ID from the physical _gateway_ device or its UI.

2. The _user_ enter the unique device ID from the _gateway_ into the _core_. 

3. The _core_ generates and displays a one time password.

4. The _user_ enters the one time password and RabbitMQ URL into the _gateway_.

5. The _gateway_ makes a registration request via RabbitMQ.

6. The _core_ authenticates the one time password; returns credentials, etc.; and deletes the OTP. 

7. The _gateway_ connects to the RabbitMQ using its new credentials.

## Routine Date Exchange

Similar to the connection flow described above, the routine data exchange between the _core_ and _gateways_ is designed such that _gateways_ can be implemented as simply as possible. 

The credentials sent to the _gateway_ during registration, include both a channel for listening and a channel for reporting. 

### Gateway to Core

Upon first connection, a _gateway_ should send information about it itself to the core via its reporting channel in the following format. An `id` must be unique within the parent level. I.e. two _devices_ may have _features_ with an `id` of `1`, but no device amy have two features with an `id` of `1`

```json
{
    "gateway": {
        "manufacturer": "",
        "model": "",
        "serialNumber": "",
        "softwareVersion": "",
        "features": {},
    },
    "devices": {
        "<id>": {
            "manufacturer": "",
            "model": "",
            "softwareVersion": "",
            "features": {
                "<id>": {
                    "name": "",
                    "description": "",
                    "standard": "",
                    "valueType": "",
                    "settable": true,
                    "settingValueRange": {
                        "unitPrefix": "",
                        "unit": "",
                        "min": 0,
                        "max": 0
                    },
                    "options": [
                        "",
                        ""
                    ],
                    "setting": {
                        "dateTime": "",
                        "unitPrefix": "",
                        "unit": "",
                        "value": ""
                    },
                    "reading": {
                        "dateTime": "",
                        "unitPrefix": "",
                        "unit": "",
                        "value": ""
                    }
                }
            }
        }
    }
}
```

See the User API documentation for more on features. 

All subsequent communications will be in the same format and a subset of the above. For example, to update a single reading, a _gateway_ may send:

```json
{
    "devices": {
        "some-unique-id": {
            "features": {
                "some-other-unique-id": {
                    "reading": {
                        "dateTime": "2020-07-19T15:07:33-04:00",
                        "unit": "DEGREES_FAHRENHEIT",
                        "value": 78
                    }
                }
            }
        }
    }
}
```

or, unformatted:

```json
{"devices":{"some-unique-id":{"features":{"some-other-unique-id":{"reading":{"dateTime":"2020-07-19T15:07:33-04:00","unit":"DEGREES_FAHRENHEIT","value":78}}}}}}
```

The core will be flexable with its parsing of the JSON, so the `value` could just as well ve sent as a string.

### Core to Gateway

The _core_ will send setting updates to a _gateway_  on the `gateway`'s listening channel and in the following format:

```json
{
    "devices": {
        "some-unique-id": {
            "features": {
                "some-other-unique-id": {
                    "setting": {
                        "unit": "DEGREES_FAHRENHEIT",
                        "value": 72
                    }
                }
            }
        }
    }
}
```
