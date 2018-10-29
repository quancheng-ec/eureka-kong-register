# eureka-kong-register

## USAGE

1.0: for Kong < 0.12

2.0: for Kong >= 0.12

```bash

docker run -d \
  --link eureka01:eureka01 \
  --link eureka02:eureka02 \
  -e KONG_HOST=http://kong:8001 \
  -e EUREKA_URLS=http://eureka01:8761/eureka|http:eureka02:8761/eureka \ 
  zephyrdev/eureka-kong-register
  
```

## Run Local

first step into your $GOPATH/src and create new folders github.com/zephyrpersonal

then run following in the path

```bash
git clone https://github.com/quancheng-ec/eureka-kong-register.git

go get github.com/tools/godep

godep restore
```
