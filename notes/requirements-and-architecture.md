# General Requirements and Architecture

_Tilde_ is an extensible, multi-purpose, IoT automation platform with a focus on automating the home. At a very high level, such a platform must be able to ingest data both from connected devices and internet sources, and send commands to connected devices that cause their state to change (i.e. light switch on/off).

## User Stories

These user stories are selected to illustrate interaction with both physical devices and web resources, and various types of challenging logic. 

1. I come home from work, the lights automatically come on.
2. An intruder alarm is activated. All of the lights come on, a siren sounds, a text message is sent, and surveillance cameras begin recording at their highest resolution. 
3. There is a severe weather alert, widows close, a text message is sent. 
4. The soil is too dry, but it is going to rain later so irrigation shouldn't be activated.
5. I manually turn a light on in the bathroom, then I leave the room without turning it off. The light turns off automatically after a delay.
6. I open a window because it's a nice day. Climate control systems turn off or adjust their settings appropriately until it is closed. 
7. The house is too warm (due to cooking), it is cooler outside than inside, the windows open to ventilate instead of using the air conditioner, unless it is raining.
8. During the day, I'm watching a movie. Starting the movie has changed living room lighting settings. The movie runs into the night. When the movie ends the room should convert to night time lighting, not day time. 

## Gross Architecture

A complete _Tilde_ system consists of four types of components.

- **Core:** The system has a single _core_ which provides interfaces for the other components, and basic functionality such an automation engine and data storage.
- **Devices:** _Devices_ can be connected items in the home such as light switches, thermostats, and door sensors. _Devices_ can also be entities on the internet such as a weather report for a certain location or a Twitter feed. 
- **Gateways:** _Gateways_ connect to both the _core_ and _devices_ and act as an intermediary between them. _Gateways_ will typically implement a single IoT protocol such as Z-Wave, or interact with one or several closely related web services.
- **User Interfaces:** The core will provide an API to which several _user interfaces_ may connect.

## Core Architecture

### General Requirements

In order to perform control and automation functions the _core_ must be able to:

1. Receive information about the state of _devices_.
2. Store and update a representation of the current state of _devices_.
3. Store and update a representation of the desired state of the _devices_.
4. Accept user input regarding the desired state of _devices_. 
5. Based on the current state of _devices_, determine what the desired state should be. This is "automation."
6. Control _devices_. 

Additionally, the history of device state should be stored in a database. This will facilitate things such as machine learning later. 

### External Interfaces

The _core_ has two interfaces: the Gateway API and the User API. The User API will be a RESTful HTTP/JSON API and will allow for the configuration and control of the _core_, _gateways_, and _devices_. The User API should be the only way the users need to interact with the system. The Gateway API will likely be connection oriented and the protocol designed to be as easy to implement as possible. A major goal of its design is to allow _gateways_ to be implemented easily. One possibility that is being considered is WebSocket/JSON with the _core_ being the server, and _gateways_ being the clients. 

### Internal Devices

The _core_ will have _internal devices_. These may include things such as a clock or metronome, network connection information, system diagnostic information (CPU, memory, etc.), and support for notifications.

### Current State

A data structure that stores the current state of _devices_.

### Layers of Desired State

A layered data structure that stores the state of _devices_ and in which the attributes in each layer override the same attributes in layers below.

0. **Emergency:** this layer has absolute highest priority
    - fire alarm, entry alarm, natural disaster 
1. **Critical:** contains critical overrides
    - if it's raining, the windows should be closed
2. **Temporary:** contains temporary, non critical state overrides 
    - watching a movie
3. **Base:** contains the normal priority state layers which are set by conditions such as
    - day/night, presence/absence
4. **Default:** contains safe defaults for all devices. 
    - lighting off, temperature above freezing, windows closed

Consideration needs to be given to the handling of manual changes to desired state. For example, what should happen if someone turns on a light using a switch? Obviously, the light would come on and the associated _gateway_ would send a message to the core updating the current state of that light, but how do we determine when that light should turn off? Should it never turn off? Probably not. Should it turn off the next time its state would normally be set? Is there a particular layer of desired state in which manual changes should be applied? I extreme cases such as a burglar alarm, should we contradict manual changes such as by turning lights back on immediately after they are turned off? 

### Rules

Rules are used to determine what the desired state should be based on the current state. 

### Automation Engine

Using rules, the Automation Engine will set the layers of desired state based on the current state. Measure will need to be taken to prevent infinite loops.

### State Automation Data Flow

1. A _gateway_ receives an update from a _device_ regarding its state. 
    - it is now raining
2. The _gateway_ sends a message containing the _device_ state change to the _core_ Device API.
3. The Device API updates the Current State data structure.
4. The Automation Engine applies rules based on the new Current State and updates the appropriate layers of Desired State.
5. The Device API sends messages the appropriate _gateways_ requesting that they change the state of one or more their _devices_.
    - close windows, turn off irrigation
6. The _gateways_ send messages to the _core_ Device API when their _device_ states change.
    - windows closed, irrigation off

There is a potential for infinite loops of the above and measures should be taken in the Automation Engine to prevent such loops. 
