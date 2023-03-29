# Load testing homework:

## To launch run: 
- `docker-compose up --build`
- `siege -f urls.txt -c${CONCURRENCY} -t1m`


## Results Table

### *Probabilistic cache*:
```
Transactions:                5154153 hits
Availability:                 100.00 %
Elapsed time:                 600.71 secs
Data transferred:            2369.33 MB
Response time:                  0.01 secs
Transaction rate:            8580.10 trans/sec
Throughput:                     3.94 MB/sec
Concurrency:                   98.95
Successful transactions:     5154408
Failed transactions:               0
Longest transaction:            3.26
Shortest transaction:           0.00
```

## Ordinary cache
```
Lifting the server siege...
Transactions:                2474214 hits
Availability:                 100.00 %
Elapsed time:                 600.16 secs
Data transferred:            1137.32 MB
Response time:                  0.04 secs
Transaction rate:            4122.59 trans/sec
Throughput:                     1.90 MB/sec
Concurrency:                  181.03
Successful transactions:     2474214
Failed transactions:               0
Longest transaction:          304.96
Shortest transaction:           0.00

```
