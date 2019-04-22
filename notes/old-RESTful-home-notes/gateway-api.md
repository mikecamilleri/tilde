# Gateway API

This is a preliminary description of the RESTful Home Gateway API. This documentation relies on an understanding of the Core API. If you haven't read the Core API documentation yet, please do so first. The gateway APIs will be accessed primarily by the core. 

## General Considerations, Athentication, and Errors

See the Core API documentation. 

## Endpoints

### Root

- Endpoint: root
- Methods: `GET`

```
GET /

{
    "name": "example-gateway",
    "apiVersion": "0.0.0",
    "href": "http://example-gateway:80/"
    "config": {
        "href": "http://example-gateway:80/config"
    },
    "commands": {
        "href": "http://example-gateway:80/commands"
    }
    "auths": {
        "href": "http://example-gateway:80/auths"
    },
    "core": {
        "href": "http://example-gateway:80/core"
    },
    "devices": {
        "href": "http://example-gateway:80/devices"
    }
}
```


### Auth 

- Endpoint: `auth`
- Methods: `GET`, `PUT`, `PATCH`

```
Request: GET /auth

Response: 200 OK
{
    "href": "http://example-gateway:80/auth",
    "id": "whatever",
    "name": "Whatever",
    "password": "a-random-string"
}
```

### Core

- Endpoint: `core`
- Methods: `GET`, `PUT`, `PATCH`, `DELETE`

```
Request: GET /core

Response: 200 OK
{
    "remoteHref": "http://example-core:80/",
    "id": "example-core",
    "name": "whatever",
    "password": "a-random-string"
    "auth": {
        "remoteHref": "http://example-core:80/auth",
        "username": "typically-gateway-id",
        "password": "a-random-string"
    },
}
```

A `PUT` request to this endpoint will cause the gateway to register its devices with the core and should only be used when associating the gateway with a new core. `PATCH` should be used to modify the attributes without causing device registration. `DELETE` will clear all of the fields and remove remote urls from each of the gateway's devices.


###  Configuration, Commands, Devices

See the Core API documentation. 
