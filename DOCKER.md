# Docker

This is a quick example of using ghost inside of a Docker container.

## Build

`docker build . -t mattwiater/ghost`

## Run

```
docker run -it --rm --name ghost \
    --memory="512m" \
    --cpus="1.5" \
    mattwiater/ghost \
    bash
```

## Interact

The docker run command above will drop you into an interactive shell inside the container.

`./ghost`

```
A versatile toolkit for network diagnostics and system information gathering, offering developers a suite of commands to scan networks, retrieve system details, and perform IP and port analyses.

Usage:
  ghost [command]

Available Commands:
  arpscan           Scans the local network using ARP to find devices
  completion        Generate the autocompletion script for the specified shell
  cpuinfo           Displays detailed CPU information such as model, cores, and frequency.
  diskusage         Displays disk usage information, including total, used, and free space.
  envvars           Displays all environment variables.
  find              Finds files with names containing a specified substring.
  fsinfo            Displays filesystem information, including type, total space, and available space.
  help              Help about any command
  hostinfo          Delivers comprehensive host system information.
  largestdirs       Lists the largest directories in a specified directory, sorted by size.
  largestfiles      Lists the largest files in a specified directory, sorted by size.
  localip           Finds an internal IPv4 address.
  loggedin          Displays currently logged-in users.
  meminfo           Displays memory usage statistics, including total, used, and free memory.
  netstat           Displays active network connections on the system
  networkinterfaces Lists all network interfaces on the host.
  portscanner       Scans a range of ports on a specified host
  services          Lists running services with their status and memory usage.
  subnetcalc        Calculates network details for a given IP address and subnet (CIDR)
  sysinfo           Displays system information such as OS, architecture, and uptime.
  treeprint         Displays a tree-like structure of files and directories.

Flags:
  -h, --help   help for ghost

Use "ghost [command] --help" for more information about a command.

```

Type: `exit` to leave the docker container and remove it.