start-server:
	go build && ./chatroom

start-client:
	go build ../client-app/. && ./client-app

start-server-pubsub:
	go build && REDIS_PUBSUB=true ./chatroom

start-redis:
	 docker-compose up -d

start-server-pubsub2:
	go build && REDIS_PUBSUB=true WEB_HOST=8001 ./chatroom

start-client2:
	go build ../client-app/. && WEB_HOST=8001 ./client-app

reset:
	rm -f chatdb.db