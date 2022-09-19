## go-elasticsearch



### SET up

```shell
$ git clone https://github.com/deviantony/docker-elk
$ cd docker-elk
$ docker-compose up -d
$ curl -XPOST -D- 'http://localhost:9200/_security/user/elastic/_password' \
    -H 'Content-Type: application/json' \
    -u elastic:changeme \
    -d '{"password":"KANG1823"}'
```