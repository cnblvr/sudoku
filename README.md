1. Create redis config and development environments
```shell
touch redis.conf
touch dev.env
# optional: set password for redis
echo "requirepass <password>" >> redis.conf
echo "REDIS_PASSWORD=<password>" >> dev.env
```

2. Run this application
```shell
sudo docker-compose up --build
```
