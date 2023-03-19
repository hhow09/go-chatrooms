start-server:
	go build && ./chatroom-channel

start-client:
	go build ../client-app/. && ./client-app

start-server-pubsub:
	go build && REDIS_PUBSUB=true ./chatroom-channel

start-redis:
	 docker-compose up -d