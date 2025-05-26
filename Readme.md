docker service create --name web --replicas 3 -p 8080:80 nginx
docker service ls
docker service ps web
docker service scale web=5
docker service rm web
