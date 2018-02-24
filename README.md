# eureka-kong-register

## USAGE

```bash

docker run -d \
  --link eureka01:eureka01 \
  --link eureka02:eureka02 \
  -e KONG_HOST=http://kong:8001 \
  -e EUREKA_URLS=http:eureka01:8761/eureka|http:eureka02:8761/eureka \
  zephyrdev/eureka-kong-register
  
```
