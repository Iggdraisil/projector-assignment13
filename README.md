# Queue benchmark homework:

## To launch run: 
- `docker-compose up --build`
- `siege -b -t1m -c50 'http://127.0.0.1:9000/produce'`
### To chose queue type and consumer number edit .env file
## Results Table
### Single consumer
| concurrency | redis   | redis rdb | redis aof | beanstalkd | beanstalkd persist |
|-------------|---------|-----------|-----------|------------|--------------------|
| 50          | 411/266 | 431/278   | 423/219   | 425/298    | 417/210            |
| 100         | 785/215 | 786/220   | 769/161   | 752/277    | 675/192            | 
| 300         | 1881/79 | 1845/79   | 1716/57   | 1457/315   | 905/315            | 
### 8 consumers (1 consumer per cpu)

| concurrency | redis    | redis rdb | redis aof | beanstalkd | beanstalkd persist |
|-------------|----------|-----------|-----------|------------|--------------------|
| 50          | 404/404  | 396/396   | 404/404   | 440/440    | 414/414            |
| 100         | 740/740  | 736/736   | 707/706   | 728/728    | 591/530            |
| 300         | 1608/465 | 1552/455  | 1578/373  | 1226/1001  | 705/673            |
