# Core Archtecture (brainstorming)

## State Manager

The state manager is the part of the core responsible for responding to data from and controlling devices using a set of rules. This is the "automation" part of the RESTful Home platform.

### Situations the state manager should be able to handle

1. During the day, I'm watching a movie. Starting the movie has changed living room lighting settings. The movie runs into the night. When the movie ends the room should convert to night time lighting, not day time. 
2. The soil is too dry, but it is going to rain later so irrigation shouldn't be activated.
3. I manually turn a light on in the bathroom, then I leave the room without turning it off. It should turn off autmatically after a delay.
4. An intruder alarm is activated. If I am home sleeping, it should react differently than if I am not home. 
5. The house is too warm (due to cooking), it is cooler outside than inside, the windows should open to ventilate instead of using the air conditioner, unless it is raining.
6. Although there may be a default, manual adjustments should be respected. If the default for a light is off, the system shouldn't fight me when I turn it on. 
7. Different temperature setting for all combinations home/not-home and day/night 
8. If a window is open, climate control should be off.

### Gross architecture

A simple, but effective design for a state manager might consist of layers of desired state and a method to update those layers and make the necessary changes to the connected devices. 

### The desired state layers

A layered data structure in which the attributes in each layer override the same attributes in layers below.

**Note:** The list below should be in reverse chronological order from four (emergency) to zero (default) but GitHub's markdown parser to parses it as an HTML `<ol>`.

4. Emergency -  this layer has absolute highest priority
    - fire alarm, entry alarm, natural disaster 
3. Critical - this group contains critical overrides
    - if it's raining, the windows should be closed
2. Temporary - this group contains temporary, non critical state overrides 
    - watching a movie
1. Base - this group contains the base state layers which are set by conditions such as
    - day/night, presence/absense, arbitrary boolean conditions (gardening).
0. Default - this group contains safe defaults for all devices. 
    - (e.g. lighting off, temperature above freezing, windows closed)

### The update method

This will likely be called each time a device's state or gateways's status are updated. In addition to using the layered desired state above, this loop will require a data strcuture for the current state, and the last state.

1. Move the `current_state` to the `last_state` (using veriable style names here for clarity).
2. Get the `current_state`.
3. For each value in the `current_state` that has changed since `last_state`, check whether any layers in `desired_state` should be changed. (This will require a mapping of current state to desired state conditions. Logic along the lines of: `IF fire_alarm.value == "on" ADD LAYER emergency_fire AT LEVEL 4`)
4. For each value in the `current_state`, check whether it is within the bounds defined by the highest layer in desired state. If not, ask the appropriate gateway to make the change. 
5. Done.

**Note:** This is a very preliminary procedure and as described likely has major bugs. Concerns that I already have include:

- Some gateways or devices might respond slowly. Do we hang until we get a response to all of our requests? Do we wait but potentially send duplicate requests on the next round?
- Is there the potential for any infinite loops? 
- How are changes made at devices such as manually turning on a light handles? My home automation system shouldn't fight me. Maybe change the desired state for that device at the highest layer in which it is currently set?

Perhaps triggering events should be treated separately from status events?

## Gateway and Device Registration

Gateways must be registered by a user. 

Devices must be registered by their gateway. The gateway field is set automatically using the auth.

Devices may act as their own gateway by being registered as a gateway and then registering themselves as a device.
