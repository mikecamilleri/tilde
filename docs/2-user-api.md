# User API

## Architecture

The User API is a simple, REST-like API.

### Endpoints

All endpoint urls follow `/BASE_URL/V0/`

### General Notes

This API accepts and returns `application/json`.

All `dateTime` values are RFC 3339 formatted strings.

## `gateways`

### `GET`

Returns an array of gateways. Any field can be queried, which will limit the results returned to those that match the query. E.g `/BASE_URL/V0/gateways?manufacturer=acme`.

```json
{
    "gateways": []
}
```

### `PATCH`

Multiple gateways may be updated at once via this endpoint by using a `PATCH` request with each gateway `id` as a key. The feature object is as described below.

```json
{
    "<id>": {},
    "<id>": {}
}
```

## `gateways/<id>`

### `GET`

Returns a specific gateway. Features on a gateway will typically be things like configurable settings. See below for the structure of a feature.

```json
{
    "id": "",
    "name": "",
    "description": "",
    "manufacturer": "",
    "model": "",
    "serialNumber": "",
    "softwareVersion": "",
    "features": []
}
```

### `PATCH`

A `PATCH` request may be used to update the following fields. Only specific fields within a feature may be updated as described later. 

```json
{
    "name": "",
    "description": "",
    "features": []
}
```

## `devices`

### `GET`

Returns an array of devices. Any field can be queried, which will limit the results returned to those that match the query. E.g `/BASE_URL/V0/devices?manufacturer=acme`.

```json
{
    "devices": []
}
```

### `PATCH`

Multiple devices may be updated at once via this endpoint by using a `PATCH` request with each device `id` as a key. The feature object is as described below.

```json
{
    "<id>": {},
    "<id>": {}
}
```


## `devices/<id>`

### `GET`

Returns a specific device. See below for the structure of a feature.

```json
{
    "id": "",
    "deviceId": "",
    "name": "",
    "description": "",
    "manufacturer": "",
    "model": "",
    "serialNumber": "",
    "softwareVersion": "",
    "features": []
}
```

### `PATCH`

A `PATCH` request may be used to update the following fields. Only specific fields within a feature may be updated as described later. `id` must be included when not in the path, such as when updating multiple devices.

```json
{
    "name": "",
    "description": "",
    "features": []
}
```

## `features`, `gateways/<id>/features`, and `devices/<id>/features`

### `GET`

Returns an array of features. Any field can be queried, which will limit the results returned to those that match the query. E.g `/BASE_URL/V0/features?standard=SWITCH`.

```json
{
    "features": []
}
```

### `PATCH`

Multiple features may be updated at once via this endpoint by using a `PATCH` request with each feature `id` as a key. The feature object is as described below.

```json
{
    "<id>": {},
    "<id>": {}
}
```


## `features/<id>`

### `GET`

Returns a specific feature.  

```json
{
    "id": "",
    "deviceId": "",
    "gatewayId": "",
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
        "id": "",
        "dateTime": "",
        "unitPrefix": "",
        "unit": "",
        "value": ""
    },
    "reading": {
        "id": "",
        "dateTime": "",
        "unitPrefix": "",
        "unit": "",
        "value": ""
    }
}
```

Features will have either a `deviceId` or a `gatewayId` depending on whether they are on a gateway or device.

`valueType` is the basic type (`STRING`, `INTEGER`, `FLOAT`, `BOOLEAN`) of the setting and reading `value` fields (shown as string here). 

Not all features will return all of these fields. `settingValueRange` is used only to set valid ranges for `INTEGER` and `FLOAT` features. The `options` array is used to define a set of acceptable `STRING` values only. 

`standard` refers to predefined set of restrictions on a feature. As an example, there may be a `TEMPERATURE` standard that requires the `valueType` to be `FLOAT` and that any `unit` be `DEGREE_CELSIUS`, `DEGREE_FAHRENHEIT`, or `KELVIN`.

### `PATCH`

The following fields may be updated using a `PATCH` request. 

`settings` may only be updated if `"settable": true`.

```json
{
    "name": "",
    "description": "",
    "setting": {
        "unitPrefix": "",
        "unit": "",
        "value": ""
    }
}
``` 

## `data`

An endpoint for getting historical setting and reading data for a feature or features. At a minimum this should support querying for a particular time range.

TODO: JSON schema and reconsider endpoint name and location. Better on a feature?

## `gateway-otp`

One time passwords for gateway pairing.

TODO: JSON schema.

## `auths`

Auths for users of the User API.

TODO: JSON schema.
