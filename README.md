## Introduction
Collective-herder is a framework for building server orchestration or parallel job execution. It is designed to let users run commands on a diverse set of servers based on the presence of facts. The transport between the client and server is provided by AMQP using a publish-subscribe model.

## Building
First, download the code
```bash
go get -v github.com/r3boot/collective-herder
```

Then, proceed into the build directory and build the various commands
```bash
cd $GOPATH/src/github.com/r3boot/collective-herder
make
```

## Usage
To start a server, run the following command:
```bash
./build/chd -d
```

To start the client, run the following command:
```bash
./build/ch -d
```

## Writing custom plugins
TODO
