docker build -t localmapper .
docker-compose up -d
docker exec -it localmapper-server_server_1 /bin/bash