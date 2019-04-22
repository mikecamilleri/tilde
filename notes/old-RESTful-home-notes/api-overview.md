# RESTful Home Overview

## Major Components

- **device**: A single physical device in the home such as a light switch or thermostat, or a virtual "device" such as an integration with a third party weather service. Devices communicate with the core via a gateway. A device may serve as its own gateway. 
- **gateway** (hub/controller/interface): Gateways serve as an intermediaries between the core and devices using their specific networking technology (E.g. Z-Wave).
- **core**: The central service that gateways communicate with.

## Version

During development, the version will be 0.0.0, The first release or stable API will be version 1.0.0. The digit after the first dot represents non-breaking feature releases. The digit after the second dot is for bug fix releases. 

## Authentication

Uses HTTP basic auth.

## Authorization

Gateways and users should only be given the minimal access necessary to do their job. 

## General Security Considerations

- For now, in order to avoid having to implement HTTPS, which would be problematic on low-power hardware such as Arduino, the local IP network that the core uses to communicate with gateways must be secure and ideally should be separate from the general internet connected home network. Ethernet jacks should be physically secured and WiFi should be encrypted appropriately.
- Communication with the core that occurs over the interent should be encrypted via HTTPS, VPN, or similar secure technology.
- Gateways are responsible for secure communication with devices connected to them. The IoT and home automation space has a lot of insecure technologies. Products should be selected that allow for reasonable security given the type of devices connected. 

## Gateway Registration Flow

2. A human admin sends a `POST` request to the `/gateways` endpoint on the core containing the url and credentials (username and password) to access the new gateway.
4. The core sends a `GET` request to the gateway's root URL to get its name, endpoint urls, current status, etc. 
3. The core creates an auth for the new gateway sends and a `PUT` request to the `/core` endpoint to set itself as the core.
4. The new gateway sends `POST` requests to the core to update its status and register each of its devices. 

It is the gateway's responsibility to make `PATCH` requests to keep its state and devices up to date in the core. 
