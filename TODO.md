# Todo

A running list of things todo to keep me on track ...

- [ ] build prototype registration api (HTTP)
- [x] decide on database (SQLite https://github.com/mattn/go-sqlite3)
- [ ] use database to store auths

- [ ] set up RabbitMQ (docker or local? CLI config or HTTP API?)
- [ ] add RabbitMQ topic creation to registration procedure

- [ ] design a format for data exchange between gateways and core (GraphQL-like?) Gateways need to declare their and their devices' interfaces to the core. Gateways need to send device state updates to the core. The core needs to send device state change requests to the Gatways. 
