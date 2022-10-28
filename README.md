# URLAsker

### Get Started


### Contracts
Project has several public contracts: [swagger](https://github.com/MikhailSolovev/URLAsker/blob/main/api/swagger.yaml),
[interfaces](https://github.com/MikhailSolovev/URLAsker/blob/main/internal/interfaces/asker.go)

### Make

`run`

`stop`

`stop-clean-volumes`

`restart`

`test-local`

### Code architecture
![Alt text](./img/clean_code_architecture.png)

### Some thoughts
Maybe it's more convenient to use clickhouse instead postgres,
therefore clickhouse is better on aggregations, which can be useful
for future analytics.
