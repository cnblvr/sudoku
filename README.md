1. Create redis config and development environments
```shell
rm -f redis.conf && touch redis.conf
rm -f dev.env && touch dev.env
echo 'SECURECOOKIE_HASH_KEY='$(head -c 32 /dev/random | base64) >> dev.env
echo 'SECURECOOKIE_BLOCK_KEY='$(head -c 32 /dev/random | base64) >> dev.env
echo 'PASSWORD_PEPPER='$(head -c 32 /dev/random | base64) >> dev.env
# optional: set password for redis
echo 'requirepass <password>' >> redis.conf
echo 'REDIS_PASSWORD=<password>\n' >> dev.env
```

2. Run this application
```shell
sudo docker-compose up --build
```

###TODO

- easy: 33-37
- medium: 28-32
- hard: 23-27