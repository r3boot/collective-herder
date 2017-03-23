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

To use the client, see the examples below:
```bash
$ ./build/ch ping
PONG response from alita.local in 3.285842ms
PONG response from alita.local in 3.860543ms
PONG response from alita.local in 4.000727ms
PONG response from alita.local in 4.103193ms

Summary: min/avg/max = 3.285842ms/3.812576ms/4.103193ms

$ ./build/ch run uname -r
alita.local         stdout: 4.9.11-1-ARCH
alita.local         stdout: 4.9.11-1-ARCH
alita.local         stdout: 4.9.11-1-ARCH
alita.local         stdout: 4.9.11-1-ARCH

$ ./build/ch facts service_mgr
Discovered the following values for service_mgr:

alita.local         systemd
alita.local         systemd
alita.local         systemd
alita.local         systemd

```

## Writing custom plugins
TODO
