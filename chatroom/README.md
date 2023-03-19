# Chatroom

## basic room
### Features
- a room is only be in 1 server
- a room can only publish to clients connecting to same server
### How to run ?
1. `make start-server`
2. `make start-client`


## room with pubsub broadcast
### Features
- a room can be on multiple server (TODO)
- a room can publish to different clients on different server

### How to run ?
1. `make start-redis` 
2. `make start-server-pubsub` (server1, to localhost `8000`)
3. open new ternminal, `make start-server-pubsub2` (server2, to localhost `8001`)
4. open new ternminal, `make start-client` (to server1)
    - create a room with `name`
5. open new ternminal, `make start-client2` (to server2)
    - join the same room created by step #4

### Reset Database 
`make reset`