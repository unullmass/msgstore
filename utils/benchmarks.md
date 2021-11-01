```shell
$ go run ../gobench/gobench.go -c 100  -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=true
Dispatching 100 clients
Waiting for results...

Requests:                             4222 hits
Successful requests:                  4222 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:               70 hits/sec
Read throughput:                     14002 bytes/sec
Write throughput:                    48478 bytes/sec
Test time:                              60 sec

$ go run ../gobench/gobench.go -c 200  -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=true
Dispatching 200 clients
Waiting for results...

Requests:                             5416 hits
Successful requests:                  2061 hits
Network failed:                        524 hits
Bad requests failed (!2xx):           2831 hits
Successful requests rate:               34 hits/sec
Read throughput:                     19386 bytes/sec
Write throughput:                    62992 bytes/sec
Test time:                              60 sec
```


mullas@mullas-MOBL MINGW64 ~/learning/msg-store (master)
$ go run ../gobench/gobench.go -c 10  -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=true
Dispatching 10 clients
Waiting for results...

Requests:                             1623 hits
Successful requests:                  1623 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:               27 hits/sec
Read throughput:                      5193 bytes/sec
Write throughput:                    18316 bytes/sec
Test time:                              60 sec

mullas@mullas-MOBL MINGW64 ~/learning/msg-store (master)
$ go run ../gobench/gobench.go -c 100  -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=true
Dispatching 100 clients
Waiting for results...

Requests:                             2744 hits
Successful requests:                  2744 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:               45 hits/sec
Read throughput:                      8780 bytes/sec
Write throughput:                    31900 bytes/sec
Test time:                              60 sec


 go run ../gobench/gobench.go -c 200   -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=true
Dispatching 200 clients
Waiting for results...

Requests:                             3100 hits
Successful requests:                  2924 hits
Network failed:                        176 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:               48 hits/sec
Read throughput:                      9356 bytes/sec
Write throughput:                    37015 bytes/sec
Test time:                              60 sec




$ go run ../gobench/gobench.go -c 500   -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=true
Dispatching 500 clients
Waiting for results...

Requests:                             5521 hits
Successful requests:                    35 hits
Network failed:                       5459 hits
Bad requests failed (!2xx):             27 hits
Successful requests rate:                0 hits/sec
Read throughput:                       224 bytes/sec
Write throughput:                    58517 bytes/sec
Test time:                              60 sec



$ go run ../gobench/gobench.go -c 200   -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=false -tr 20000
Dispatching 200 clients
Waiting for results...

Requests:                             2673 hits
Successful requests:                  2673 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:               44 hits/sec
Read throughput:                      9400 bytes/sec
Write throughput:                    31974 bytes/sec
Test time:                              60 sec


$ go run ../gobench/gobench.go -c 400   -t 60  -d test/createdoc.json -u http://localhost:8080/mydata -k=false -tr 20000
Dispatching 400 clients
Waiting for results...

Requests:                             2485 hits
Successful requests:                  2485 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:               41 hits/sec
Read throughput:                      8738 bytes/sec
Write throughput:                    32119 bytes/sec
Test time:                              60 sec

