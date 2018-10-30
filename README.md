# eureka-kong-register

Scheduled task that keep syncing kong upstreams with eureka instances

## INTRO

say we have an application with 2 instances on Eureka: (I try to describe by JSON)

```json
{ "name": "APP:1.0.0",
  "instances": [{
      ipAddr: 192.168.0.1
   },{
      ipAddr: 192.168.0.2
   }]
}
```

and the register will work as follow:

- polling fetch http://{eurekaHost}/eureka/apps to get application and instance information
- create a kong upstream object named "APP-1-0-0.eureka.internal" if it's not existed
- create new targets by newly registered instance. Target's ip equals to instance's ipAddr
- remove unregistered ones

## WHICH VERSION TO USE

1.x: for Kong < 0.12

2.x: for Kong >= 0.12

## USAGE

### DOCKER

```bash

docker run -d \
  --link eureka01:eureka01 \
  --link eureka02:eureka02 \
  -e KONG_HOST=http://kong:8001 \
  -e EUREKA_URLS=http://eureka01:8761/eureka|http:eureka02:8761/eureka \ 
  zephyrdev/eureka-kong-register
  
```

### WITHOUT DOCKER

first step into your $GOPATH/src and create new folders github.com/zephyrpersonal

then cd to it and tun

```bash
git clone https://github.com/quancheng-ec/eureka-kong-register.git

go get github.com/tools/godep

godep restore

go build main.go
```

