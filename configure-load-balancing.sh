#!/usr/bin/env bash

# create upstream
curl -X POST http://localhost:8001/upstreams --data "name=price.v1.service"

# add two targets to the upstream
curl -X POST http://localhost:8001/upstreams/price.v1.service/targets \
    --data "target=kong-grpc-lb_api_1:50051"

curl -X POST http://localhost:8001/upstreams/price.v1.service/targets \
    --data "target=kong-grpc-lb_api_2:50051"

curl -XPOST http://localhost:8001/services/ \
  --data name=grpc \
  --data protocol=grpc \
  --data "host=price.v1.service" \
  --data port=15002

curl -XPOST http://localhost:8001/services/grpc/routes \
  --data protocols=grpc \
  --data name=catch-all \
  --data paths=/