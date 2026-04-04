# go-clean-architecture
Desafio 3 do curso Pós Go Expert - Full Cycle

----

# clean architecture: listagem de orders

Desafio Go Expert - Use case `ListOrders` exposto via **REST**, **gRPC** e **GraphQL** simultaneamente, partindo do codigo base do curso (20-CleanArch). Base repository: https://github.com/devfullcycle/goexpert/tree/main/20-CleanArch

## Como rodar

```bash
docker-compose up
```

Sobe o MySQL, RabbitMQ, aplica as migracoes e inicia a aplicacao automaticamente.

## Portas

| Servico  | Porta |
|----------|-------|
| REST     | 8000  |
| gRPC     | 50051 |
| GraphQL  | 8080  |

## Exemplos de uso

### REST

```bash
# Criar order
curl -s -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"6d68570a-54df-41c3-94dd-3c125db5e70f","price":100.0,"tax":10.0}'

# Listar orders
curl -s http://localhost:8000/order
```

### GraphQL

Disponivel em http://localhost:8080/

```bash
# Criar order
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation{createOrder(input:{id:\"eeb2a073-3871-4f99-9403-2329c6df2a00\",Price:200,Tax:20}){id Price Tax FinalPrice}}"}'

# Listar orders
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{listOrders{id Price Tax FinalPrice}}"}'
```

### gRPC (via grpcurl)

```bash
# Criar order
grpcurl -plaintext -d '{"id":"eea517cb-4b8d-4608-8d61-a17649579f91","price":300,"tax":30}' \
  localhost:50051 pb.OrderService/CreateOrder

# Listar orders
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```
