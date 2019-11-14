# Log Parser

Parse a log file to obtain valuable info. See the docs for more detailed description of the requirements



## Getting started

### Requirements

You need one of the two (or both):

* Docker
* Go 12

You can run the app in docker or compile it yourself to run it as a binary file.

### Makefile

A Makefile is provided to automate most of the stuff

```bash
$ make
help                           Print this message.
build                          Build the binary. Place it under bin
test                           Run tests
build-docker                   Build the docker image
tag                            Tag the docker images with version and commit
push                           Push to the docker repository
```

### Build

**Docker**

```sh
$ make build-docker
building image log-parser 0.1.0 5c2054e20c0f64f5a858c90c923bae96041264d8
docker build --build-arg VERSION=0.1.0 --build-arg GIT_COMMIT=5c2054e20c0f64f5a858c90c923bae96041264d8 --label version.type=dev -t "pablosjv/log-parser":local .
Sending build context to Docker daemon  16.28MB
[...]
Successfully tagged pablosjv/log-parser:local
```
Now you can use the image `pablosjv/log-parser:local`

**Go**

```sh
$ make build
building log-parser 0.1.0
GOPATH=/Users/pablosjv/go
go build -o bin/log-parser cmd/*.go
```

Now you can find the binary file in the `./bin` directory

### Run

```sh
$ ./bin/log-parser -h
Usage of ./bin/log-parser:
  -e int
        End time in unix format
  -f string
        Log filename to parse
  -h    Print this help message
  -i int
        Init time in unix format
  -originHost string
        Hostname that connect to others
  -period int
        Period in seconds to report (default 3600)
  -t    Tail the log file and report
  -targetHost string
        Hostname to wich other connects
  -v    Print the version and exits
```

The tool expose to execution types

1. Get all the connections to a given host in a period of time
2. Report every time period the following data:
    * a list of hostnames connected to a given (configurable) host during the last hour
    * a list of hostnames received connections from a given (configurable) host during the last hour
    * the hostname that generated most connections in the last hour

For the case 1:

```sh
$ ./bin/log-parser -f test/input-file-10000.txt -targetHost Sidney -i 1565647204350 -e 1665647328946
2019/11/14 22:57:36 Processing connected hosts
2019/11/14 22:57:36 End of file
2019/11/14 22:57:36 Connected hosts to Sidney : map[Adalhi:1 Alexionna:1 Anishka:1 Dyshawn:1 Kya:1 Trestyn:1 Wadie:1 Zyleah:1]
2019/11/14 22:57:36 Processing time: 20.4148ms
```

For the case 2:

```sh
$ ./bin/log-parser -f test/input-file-1.txt --targetHost Zyla -t -period 10 -originHost Aadvik
```

If you use Docker, the commands have similar structure. However you have to mount the file to process in the docker filesystem

```sh
$ docker run --rm -it -v $(pwd):/data pablosjv/log-parser:local -f data/test/input-file-10000.txt -targetHost Sidney -i 1565647204350 -e 1665647328946
2019/11/14 22:57:36 Processing connected hosts
2019/11/14 22:57:36 End of file
2019/11/14 22:57:36 Connected hosts to Sidney : map[Adalhi:1 Alexionna:1 Anishka:1 Dyshawn:1 Kya:1 Trestyn:1 Wadie:1 Zyleah:1]
2019/11/14 22:57:36 Processing time: 20.4148ms
```
