# Chatroom

## features

### 1. basic room
- a room can only be in 1 server
- a room can only publish to clients connecting to same server

### 2. room with pubsub broadcast
- a room can be on multiple server (TODO)
- a room can only publish to clients connecting to same server

## How to run
1. `make start-redis` (for pubsub only)
2. `make start-server` or `make start-server-pubsub`
3. `make start-client`