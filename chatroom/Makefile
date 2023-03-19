start-server:
	go build && ./chatroom

start-client:
	go build ../client-app/. && ./client-app

start-server-pubsub:
	go build && REDIS_PUBSUB=true ./chatroom

start-redis:
	 docker-compose up -d