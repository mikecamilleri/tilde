# First steps

1. Read device definitions from configuration file and register devices.
2. Accept data from registered devices via a RESTful API.
3. Emit events from the API to a pub/sub type service.
4. Read events from the pub/sub service and store in database (InfluxDB).

- use `unittest`, `http.server` or `urllib`, `json`, PyYAML, `logging`
- Flask or Falcon or Bottle instead of `http.server`, or something else?
- package management?
- `blinker`?

# Desired State Management

## Situations the state manager should be able to resolve

1. During the day, I am watching a movie. Starting the movie has changed living room lighting settings. The movie runs into the night. When the movie ends the room should convert to night time lighting, not day time. 

2. The soil is too dry, but it is going to rain later so irrigation shouldn't be activated.

3. I manually turn a light on in the bathroom, then I leave the room without turning it off. It should turn off autmatically after a delay.

4. An intruder alarm is activated. If I am home sleeping, it should react differnetly than if I am not home. 

5. The house is too warm (due to cooking), it is cooler outside than inside, the windows should open to ventilate instead of using the air conditioner, unless it is raining.

6. Although there may be a default, manual adjustments should be respected. If the default for a light is off, the system shouldn't fight me when I turn it on. 

7. Different temperature setting for all combinations home/not-home and day/night 

8. If a window is open, climate control should be off.

## State manager architecture

1. When an event is received
    - update the current state
    - take prescribed actions
        - could update desired state
2. When current state is changed
    - test whether current state out of bounds
3. When desired state is changed
    - test whether current state out of bounds
4. When current state out of bounds
    - emit events that say what to do, but not how to do it
        (i.e. raise temperature in living room, not turn on heater in living room)
        - sometimes nothing should be done

## Structure of desired state
- a quasi stack in which the attributes in each layer override the same attributes in layers below

0. Default - this group contains safe defaults for all devices. 
    - (e.g. lighting off, temperature above freezing, windows closed)
1. Base - this group contains the base state layers which re set by conditions such as
    - day/night, presence/absense, arbitrary boolean conditions (gardening).
2. Temporary - this group contains temporary, non critical state overrides 
    - watching a movie
3. Critical - this group contains critical overrides
    - if it it raining, don't open windows 
4. Emergency -  this layer has absolute highest priority
    - fire alarm
    - entry alarm
    - natural disaster 
    
## Example of state stack

0. Default
    - LR lighting: off
    - BR lighting: off
    - outdoor lighting: on
    - temperature: 55 - 85 F
    - soil moisture: >= 0
    - perimiter alarm: on
    
1. Base
    - is day
        - outdoor lighting: off
        - temperature: 68 - 72
    - is night
        - outdoor lighting: on
        - temperature: 65 - 75
        
2. Temporary
    - am watching movie
        - LR lighting: low
        - LR shades: closed

## Problems

1. What happens when current state hasn't caught up to desired state? Are actions repeatedy taken?

## 