# Docker

This is a quick example of using ghost inside of a Docker container. Not all commands work due to the inner workings of a Docker container differeing from a full-fledged OS, but there are some useful examples below. This tool does not intend to to provide full support in Docker containers, the following are examples to illustrate a different context.

## Build

Assuming you are in the root of the project.

`docker build . -t mattwiater/ghost`

## Run

`docker run -it --rm --name ghost mattwiater/ghost bash`

Type: `exit` to leave the docker container and remove it.

## Interact

`./ghost sysinfo`

```
 sysInfoCmd
 SYSTEM INFO   VALUE
 OS            alpine 3.20.3
 Architecture  amd64
 Kernel        5.15.153.1-microsoft-standard-WSL2 
 Uptime        6m12s
 ```

`./ghost hostinfo`

 ```
 hostInfoCmd
 HOST INFO             VALUE
 Hostname              251068e18f11
 Uptime                438
 BootTime              1730940012
 Procs                 2
 OS                    linux
 Platform              alpine
 PlatformFamily        alpine
 PlatformVersion       3.20.3
 KernelVersion         5.15.153.1-microsoft-standard-WSL2   
 KernelArch            x86_64
 VirtualizationSystem  docker
 VirtualizationRole    guest
 HostID                2c3fb449-0afe-41cf-a815-cd08913c3471 
 ```

`./ghost networkinterfaces`

 ```
 networkInterfacesCmd
 NETWORK INTERFACES  VALUE
 Index               1
 MTU                 65536
 Name                lo
 HardwareAddr
 Flags               up
                     loopback
 Addrs               Addr: 127.0.0.1/8   
 ----------
 Index               17
 MTU                 1500
 Name                eth0
 HardwareAddr        02:42:ac:11:00:02   
 Flags               up
                     broadcast
                     multicast
 Addrs               Addr: 172.17.0.2/16 
 ----------
 ```