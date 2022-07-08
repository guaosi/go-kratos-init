# go-kratos-init

## 组件
- 微服务框架：go-kratos
- 日志：zap
- 数据库：gorm
- 链路追踪：jaeger
- 缓存：redis
- 认证：jwt
- 配置中心：consul
- 注册发现：consul
- 文档：swagger
- 参数校验：protoc-gen-validate

## 配置中心配置参考
### shop
文件名 `dev/shop.yaml`
```
app:
  name: frontend.guaosi.shop.service # 本服务的服务名，在注册中心上体现
  version: v1 # 本服务的版本号
  debug: false # 是否启用debug模式
  log_path: /Users/guaosi/coding/log/ # 日志存放路径
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/test
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
trace:
  endpoint: http://127.0.0.1:14268/api/traces # 链路追踪地址
service:
  user: backend.guaosi.user.service # 连接用户服务，用户服务所用的服务名
auth:
  user_service_key: aaa # 连接用户服务，用的jwt加密key
  api_key: bbb # 本服务用的jwt加密key
  timeout: 3600 # 本服务颁发jwt，jwt存活时间，单位s
  tls_ca_crt_path: /Users/guaosi/coding/go/cert/ca.crt # 连接其他服务所用的tls客户端证书
  tls_server_name: www.kratos.com # 连接其他服务所用的tls证书里的服务名
```

### user
文件名 `dev/user.yaml`
```
app:
  name: backend.guaosi.user.service # 本服务的服务名，在注册中心上体现
  version: v1 # 本服务的版本号
  debug: false # 是否启用debug模式
  log_path: /Users/guaosi/coding/log/ # 日志存放路径
server:
  grpc:
    addr: 0.0.0.0:9001
    timeout: 1s
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
trace:
  endpoint: http://127.0.0.1:14268/api/traces # 链路追踪地址
auth:
  api_key: aaa # 本服务用的jwt加密key
  tls_server_crt_path: /Users/guaosi/coding/go/cert/server.crt # 本服务grpc安全认证用的服务端证书
  tls_server_key_path: /Users/guaosi/coding/go/cert/server.key # 本服务grpc安全认证用的服务端key
```

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/backend/server/server.proto

# Generate the proto code
kratos proto client api/backend/server/server.proto

# Generate the source code of service by proto file
mkdir -p app/backend/server/internal/service
kratos proto server api/backend/server/server.proto -t app/backend/server/internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```