# delivery

### http server generation command
```
oapi-codegen -config configs/server.cfg.yaml https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/openapi.yml
```

### http client generation command
```
protoc --go_out=./internal/generated/clients --go-grpc_out=./internal/generated/clients ./api/proto/geo_service.proto
```

### kafka message generation command
```
protoc --go_out=./internal/generated/events ./api/proto/basket_confirmed.proto

protoc --go_out=./internal/generated/events ./api/proto/order_status_changed.proto
```