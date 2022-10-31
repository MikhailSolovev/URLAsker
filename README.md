# URLAsker

### Get Started

API prefix: `/api/v1`

Service started with default interval of 5s and empty set of urls. For get service information, use `/info` handler.
For change settings use: `/setInterval` `/setUrls` `/addUrls` `/deleteUrls`. For get results of asker work: `/list`
`/listLatest`

### Contracts
Project has several public contracts: [swagger](https://github.com/MikhailSolovev/URLAsker/blob/main/api/swagger.yaml),
[interfaces](https://github.com/MikhailSolovev/URLAsker/blob/main/internal/interfaces/asker.go)

### Make

`run` - build and start service

`stop` - stop service

`stop-clean-volumes` - stop service and clean persistence volumes

`restart` - make new image of asker and restart service

`test-local` - run local tests

`gen-swagger` - generate swagger file in api directory

### Code architecture
![Alt text](./img/clean_code_architecture.png)

### Some thoughts
Maybe it's more convenient to use clickhouse instead postgres,
therefore clickhouse is better on aggregations, which can be useful
for future analytics.
