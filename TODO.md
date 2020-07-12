# Todo

A running list of things todo to keep me on track ...

- [ ] set up RabbitMQ (docker or local? CLI config or HTTP API?)
- [x] decide on database (SQLite https://github.com/mattn/go-sqlite3)
- [ ] configure database to store auths
- [ ] build auth engine in core 
- [ ] build rudementary gateway for testing
- [ ] get core-RabbitMQ-gateway pairing procedure working including RabbitMQ topic creation.
- [ ] design a format for data exchange between gateways and core (GraphQL-like?) Gateways need to declare their and their devices' interfaces/capabilities to the core. Gateways need to send device state updates to the core. The core needs to send device state change requests to the Gatways. 
- [ ] build out structs and DB schema for gateways, devices, and features.
