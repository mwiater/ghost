# ghost

![Ghost Gopher Mascot](ghostpher-sm.png)

_Golang + Host Inspection = Ghost!_

This project provides a suite of tools for system and network management, all packaged as a versatile CLI application built with Go and Cobra. It offers features like network scanning, system info retrieval, and more, all of which are compiled and released using GoReleaser.

---

## Table of Contents

1. [Cloning and Setting Up the Project](#cloning-and-setting-up-the-project)
2. [Building the Project with GoReleaser](#building-the-project-with-goreleaser)
3. [Cobra Commands](#cobra-commands)
   - [Command List](#command-list)
   - [Command Details](#command-details)

---

## Cloning and Setting Up the Project

To get started, follow these steps to clone the repository and install dependencies:

### Step 1: Clone the Repository

Clone the repository using `git`:

```bash
git clone git@github.com:mwiater/ghost.git
cd ghost
```

### Step 2: Install Dependencies

This project relies on several Go modules specified in `go.mod`. Install them with the following command:

```bash
go mod tidy
```

This command will fetch all required dependencies, ensuring a smooth build process.

---

## Building the Project with GoReleaser

We use **GoReleaser** to manage builds and releases. Follow these instructions to build a snapshot release, which is ideal for local testing.

### Step 1: Install GoReleaser

If you haven’t already installed GoReleaser, you can do so by following the [GoReleaser installation instructions](https://goreleaser.com/install/).

### Step 2: Build a Snapshot Release

To build `ghost` locally, use **GoReleaser** with snapshot mode. This setup builds the project without creating a formal release, making it ideal for local testing. 

### Local Build Command

Run the following command to build the application locally:

```bash
goreleaser release --snapshot --clean --skip archive
```

This command:
- Builds the application for all specified platforms in the `.goreleaser.yml` file.
- Cleans the `dist` directory before each build.
- Skips creating an archive or release, building only a snapshot.

For this example application and configuration, this command builds the following files:

```
├── dist
│   ├── artifacts.json
│   ├── config.yaml
│   ├── linux64
│   │   └── ghost
│   ├── linuxarm64
│   │   └── ghost
│   ├── metadata.json
│   └── win64
│       └── ghost.exe
```

For more detailed setup instructions, see the [GORELEASER.md](GORELEASER.md) file. This document includes installation steps, configuration examples, and validation commands to set up GoReleaser for `ghost`. 

---

## Cobra Commands

This project includes several commands that allow users to interact with system and network utilities. Below is a detailed list of each command, including usage examples and flag definitions.

### Command List

- `arpscanner`: Scans the network for active devices.
- `cpuinfo`: Retrieves detailed CPU information.
- `diskusage`: Shows disk usage statistics.
- `envvars`: Lists all environment variables.
- `find`: Searches for files or directories based on the specified parameters.
- `fsinfo`: Displays information about the file system.
- `getservices`: Lists active services on the system.
- `hostinfo`: Provides general information about the host.
- `largestdirs`: Finds the largest directories.
- `largestfiles`: Finds the largest files.
- `localip`: Shows the local IP address.
- `meminfo`: Retrieves memory usage information.
- `netstat`: Shows network status and connections.
- `networkinterfaces`: Lists all network interfaces.
- `portscanner`: Scans for open ports on the network.
- `subnetcalc`: Calculates subnet information.
- `treeprint`: Prints directory structure in a tree format.

---

### Command Details

Here is a detailed breakdown of each command, including examples and flag definitions.

---

####  `arpscanner`

**Description:** Scans the network for active devices and shows their IP addresses.

```bash
./ghost arpscanner --interface eth0
```

**Flags:**
- `--interface`: Specifies the network interface to scan (e.g., `eth0`).

Example Output:

```
 Windows ARP Table
 OUTPUT

 Interface: 192.168.213.1 --- 0x5
   Internet Address      Physical Address      Type       
   192.168.213.254       00-50-56-f1-e2-70     dynamic    
   192.168.213.255       ff-ff-ff-ff-ff-ff     static     
   224.0.0.9             01-00-5e-00-00-09     static     
   224.0.0.22            01-00-5e-00-00-16     static     
   224.0.0.251           01-00-5e-00-00-fb     static     
   224.0.0.252           01-00-5e-00-00-fc     static     
   239.255.255.250       01-00-5e-7f-ff-fa     static     
   255.255.255.255       ff-ff-ff-ff-ff-ff     static     

 Interface: 169.254.67.233 --- 0x6
   Internet Address      Physical Address      Type       
   169.254.255.255       ff-ff-ff-ff-ff-ff     static     
   224.0.0.9             01-00-5e-00-00-09     static     
   224.0.0.22            01-00-5e-00-00-16     static     
   224.0.0.251           01-00-5e-00-00-fb     static     
   224.0.0.252           01-00-5e-00-00-fc     static     
   239.255.255.250       01-00-5e-7f-ff-fa     static     
   255.255.255.255       ff-ff-ff-ff-ff-ff     static     

 Interface: 192.168.142.244 --- 0xd
   Internet Address      Physical Address      Type       
   192.168.142.113       ca-52-50-74-df-ee     dynamic    
   192.168.142.255       ff-ff-ff-ff-ff-ff     static     
   224.0.0.9             01-00-5e-00-00-09     static     
   224.0.0.22            01-00-5e-00-00-16     static     
   224.0.0.251           01-00-5e-00-00-fb     static     
   224.0.0.252           01-00-5e-00-00-fc     static     
   239.255.255.250       01-00-5e-7f-ff-fa     static     
   255.255.255.255       ff-ff-ff-ff-ff-ff     static     

 Interface: 192.168.116.1 --- 0xe
   Internet Address      Physical Address      Type       
   192.168.116.254       00-50-56-e4-8f-4c     dynamic    
   192.168.116.255       ff-ff-ff-ff-ff-ff     static     
   224.0.0.9             01-00-5e-00-00-09     static     
   224.0.0.22            01-00-5e-00-00-16     static     
   224.0.0.251           01-00-5e-00-00-fb     static     
   224.0.0.252           01-00-5e-00-00-fc     static     
   239.255.255.250       01-00-5e-7f-ff-fa     static     
   255.255.255.255       ff-ff-ff-ff-ff-ff     static     

 Interface: 169.254.76.71 --- 0x1d
   Internet Address      Physical Address      Type       
   169.254.255.255       ff-ff-ff-ff-ff-ff     static     
   224.0.0.9             01-00-5e-00-00-09     static     
   224.0.0.22            01-00-5e-00-00-16     static     
   224.0.0.251           01-00-5e-00-00-fb     static     
   224.0.0.252           01-00-5e-00-00-fc     static     
   239.255.255.250       01-00-5e-7f-ff-fa     static     
   255.255.255.255       ff-ff-ff-ff-ff-ff     static  
```

---

####  `cpuinfo`

**Description:** Displays CPU information such as model, cores, and usage.

```bash
./ghost cpuinfo
```

**Flags:** None

Example Output:

```
 cpuInfoCmd
 CPU INFO   VALUE
 CPU 1
 ModelName  Intel(R) Core(TM) i7-4710HQ CPU @ 2.50GHz
 Cores      8
 Frequency  2.50 GHz
```

---

####  `diskusage`

**Description:** Shows disk usage statistics for the system.

```bash
./ghost diskusage --path / --limit 10
```

**Flags:**
- `--path`: Specifies the path to check disk usage.
- `--limit`: Limits the number of entries displayed.

Example Output:

```
 diskUsageCmd
 MOUNT POINT  TOTAL SPACE  USED SPACE  FREE SPACE  USED PERCENT 
 C:           466.61 GB    353.14 GB   113.47 GB   75.68%       
 D:           26.84 GB     3.05 GB     23.80 GB    11.34%
 G:           466.61 GB    358.82 GB   107.79 GB   76.90%  
```

---

####  `envvars`

**Description:** Lists all environment variables.

```bash
./ghost envvars
```

**Flags:** None

Example Output:

```
envVarsCmd
VARIABLE                             VALUE
ALLUSERSPROFILE                      C:\ProgramData
APPDATA                              C:\Users\Matt\AppData\Roaming
COMMONPROGRAMFILES                   C:\Program Files\Common Files
COMSPEC                              C:\WINDOWS\system32\cmd.exe
CONFIG_SITE                          C:\Program Files\Git\etc\config.site    
ChocolateyInstall                    C:\ProgramData\chocolatey
```

---

####  `find`

**Description:** Searches for files or directories based on specified parameters.

```bash
./ghost find .go
```

**Flags:**
- `--name`: Pattern to match file or directory names.
- `--path`: Directory path to search in.

Example Output:

```
 findCmd
 MATCHING FILES
 C:\Users\Matt\projects\ghost\cmd\arpscanner.go        
 C:\Users\Matt\projects\ghost\cmd\cpuInfo.go
 C:\Users\Matt\projects\ghost\cmd\demo.go
 C:\Users\Matt\projects\ghost\cmd\diskUsage.go
 C:\Users\Matt\projects\ghost\cmd\envVars.go
```

---

####  `fsinfo`

**Description:** Provides information about the file system.

```bash
./ghost fsinfo
```

**Flags:** None

Example Output:

```
 fsInfoCmd                                                                   
 FILESYSTEM  TYPE  TOTAL SPACE  USED SPACE  AVAILABLE SPACE  USED PERCENT    
 C:                466.61 GB    353.14 GB   113.47 GB        75.68%          
 D:                26.84 GB     3.05 GB     23.80 GB         11.34%
 G:                466.61 GB    358.82 GB   107.79 GB        76.90%   
```

---

####  `hostinfo`

**Description:** Displays general information about the host, including OS and architecture.

```bash
./ghost hostinfo
```

**Flags:** None

Example Output:

```
 hostInfoCmd
 HOST INFO             VALUE
 Hostname              Laptop08
 Uptime                1129014
 BootTime              1729639261
 Procs                 281
 OS                    windows
 Platform              Microsoft Windows 11 Home
 PlatformFamily        Standalone Workstation
 PlatformVersion       10.0.19045 Build 19045
 KernelVersion         10.0.19045 Build 19045
 KernelArch            x86_64
 VirtualizationSystem
 VirtualizationRole
 HostID                1a8f4fa2-a7c7-a7c6-a9c9-7cbf6af1a994 
```

---

####  `largestdirs`

**Description:** Finds the largest directories within a specified path.

```bash
./ghost largestdirs --path /home --limit 5
```

**Flags:**
- `--path`: Path to search for large directories.
- `--limit`: Maximum number of entries to display.

Example Output:

```
 largestDirsCmd
 DIRECTORY PATH  SIZE (MB) 
 .                   71.68 
 .git                49.90
 dist                21.66 
 cmd                  0.07
 utils                0.01 
```

---

####  `largestfiles`

**Description:** Finds the largest files within a specified path.

```bash
./ghost largestfiles --path /var/log --limit 5
```

**Flags:**
- `--path`: Path to search for large files.
- `--limit`: Maximum number of entries to display.

Example Output:

```
 largestFilesCmd
 FILE PATH                                                             SIZE (MB) 
 .git\objects\pack\pack-0eb34eb50c013b3d04285bf919e9d6ae897e937b.pack  40.40     
 dist\win64\ghost.exe                                                  7.43
 dist\linux64\ghost                                                    7.23      
 dist\linuxarm64\ghost                                                 7.00
 .git\objects\f8\9e9afc27cde6241c1fd8612b8023f7e63cf998                3.29      
 .git\objects\1d\f0ca22765a29aa0d9cda84b90c31755fef4362                3.22
 .git\objects\bc\8a2f62a3a992711d907205f87a2f1a466e0c91                2.95      
 go.sum                                                                0.01
 .git\objects\pack\pack-0eb34eb50c013b3d04285bf919e9d6ae897e937b.idx   0.01      
 cmd\portscanner.go                                                    0.01
 cmd\largestDirs.go                                                    0.01      
 README.md                                                             0.01
 utils\packageList.go                                                  0.01      
 GORELEASER.md                                                         0.00
 .git\hooks\pre-rebase.sample                                          0.00      
 .git\hooks\fsmonitor-watchman.sample                                  0.00
 dist\config.yaml                                                      0.00      
 cmd\networkInterfaces.go                                              0.00
 .git\index                                                            0.00      
 SCRIPTS.md                                                            0.00
```

---

####  `localip`

**Description:** Shows the local IP address of the system.

```bash
./ghost localip
```

**Flags:** None

Example Output:

```
 localIPCmd
 LOCAL IP    VALUE
 IP Address  192.168.116.1 
```

---

####  `loggedin`

**Description:** Shows the local IP address of the system.

```bash
./ghost loggedin
```

**Flags:** None

Example Output:

```
 loggedInCmd
 USER              TERMINAL  HOST      
 Laptop12\Matt  N/A       localhost 
```

---

####  `meminfo`

**Description:** Displays memory usage statistics.

```bash
./ghost meminfo
```

**Flags:** None

Example Output:

```
 memInfoCmd
 MEMORY INFO  VALUE    
 Total        17.09 GB 
 Used         9.41 GB
 Free         7.68 GB  
 UsedPercent  55.00%
```

---

####  `netstat`

**Description:** Shows network status, including open connections.

```bash
./ghost netstat
```

**Flags:** None

Example Output:

```
Active Network Connections
 PROTOCOL  LOCAL ADDRESS                                  REMOTE ADDRESS                             STATE       
 TCP       0.0.0.0:135                                    0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:445                                    0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:902                                    0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:912                                    0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:5040                                   0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:5357                                   0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:49664                                  0.0.0.0:0                                  LISTEN      
 TCP       0.0.0.0:49665                                  0.0.0.0:0                                  LISTEN  
```

---

####  `networkinterfaces`

**Description:** Lists all network interfaces on the system.

```bash
./ghost networkinterfaces
```

**Flags:** None

Example Output:

```
networkInterfacesCmd
 NETWORK INTERFACES  VALUE
 Index               16
 MTU                 1500
 Name                Ethernet
 HardwareAddr        f8:a9:5e:3c:7d:42
 Flags               broadcast
                     multicast
 Addrs               Addr: fe80::92c7:819:a487:b76d/64
                     Addr: 169.254.245.133/16
 ----------
 Index               6
 MTU                 1500
 Name                VirtualBox Host-Only Network #2
 HardwareAddr        0a:00:27:00:00:06
 Flags               up
                     broadcast
                     multicast
 Addrs               Addr: fe80::87db:86ed:1d04:29dd/64
                     Addr: 169.254.67.233/16
 ----------
 Index               29
 MTU                 1500
 Name                VirtualBox Host-Only Network
 HardwareAddr        0a:00:27:00:00:1d
 Flags               up
                     broadcast
                     multicast
 Addrs               Addr: fe80::b583:79ce:27fc:72c6/64
                     Addr: 169.254.76.71/16
```

---

####  `portscanner`

**Description:** Scans for open ports on the specified host.

```bash
./ghost portscanner --host 192.168.1.1 --ports 20-80
```

**Flags:**
- `--host`: Host IP to scan.
- `--ports`: Range of ports to scan.

Example Output:

```
Port Scan Results
 PORT  PROTOCOL  LOCAL ADDRESS  FOREIGN ADDRESS  STATE      PROCESS           PID   OWNER 
 135   TCP       0.0.0.0:135    0.0.0.0:0        LISTENING  svchost.exe       1184  N/A   
 445   TCP       0.0.0.0:445    0.0.0.0:0        LISTENING  System            4     N/A
 902   TCP       0.0.0.0:902    0.0.0.0:0        LISTENING  vmware-authd.exe  5272  N/A   
 912   TCP       0.0.0.0:912    0.0.0.0:0        LISTENING  vmware-authd.exe  5272  N/A
```

---

####  `services`

**Description:** Lists active services on the system.

```bash
./ghost services
```

**Flags:**
- `--filter`: Filter services based on a pattern.

Example Output:

```
 servicesCmd
 SERVICE NAME  STATUS       MEMORY USAGE 
 Energy        Manager      44.04 MB     
 Memory        Compression  137.59 MB
 Secure        System       39.40 MB     
```

---

####  `subnetcalc`

**Description:** Calculates subnet information based on input IP and mask.

```bash
./ghost subnetcalc --ip 192.168.1.10 --mask 24
```

**Flags:**
- `--ip`: IP address.
- `--mask`: Subnet mask.

Example Output:

```
 Subnet Calculation Results
 FIELD              VALUE
 Network Address    192.168.1.0
 Broadcast Address  192.168.1.255
 IP Range           192.168.1.0 - 192.168.1.255 
```

---

####  `sysinfo`

**Description:** Calculates subnet information based on input IP and mask.

```bash
./ghost sysinfo
```

**Flags:** None

Example Output:

```
 sysInfoCmd
 SYSTEM INFO   VALUE
 OS            Microsoft Windows 11 Home 10.0.19045 Build 19045 
 Architecture  amd64
 Kernel        11.0.19045 Build 19025
 Uptime        313h45m48s
```

---

####  `treeprint`

**Description:** Prints the directory structure in a tree format.

```bash
./ghost treeprint --path /home
```

**Flags:**
- `--path`: Root directory to print.

Example Output:

```.
├── .gitignore
├── .goreleaser.yaml
├── GORELEASER.md
├── LICENSE
├── README.md
├── SCRIPTS.md
├── SETUP.md
├── cmd
│   ├── arpscanner.go
│   ├── cpuInfo.go
│   ├── demo.go
│   ├── diskUsage.go
│   ├── envVars.go
│   ├── find.go
│   ├── fsInfo.go
│   ├── getServices.go
│   ├── getServices_linux.go
│   ├── getServices_windows.go
│   ├── hostInfo.go
│   ├── largestDirs.go
│   ├── largestFiles.go
│   ├── list.go
│   ├── localIP.go
│   ├── loggedIn.go
│   ├── loggedin_linux.go
│   ├── loggedin_windows.go
│   ├── memInfo.go
│   ├── netstat.go
│   ├── networkInterfaces.go
│   ├── portscanner.go
│   ├── root.go
│   ├── subnetcalc.go
│   ├── sysinfo.go
│   └── treePrint.go
├── dist
│   ├── artifacts.json
│   ├── config.yaml
│   ├── linux64
│   │   └── ghost
│   ├── linuxarm64
│   │   └── ghost
│   ├── metadata.json
│   └── win64
│       └── ghost.exe
├── go.mod
├── go.sum
├── main.go
├── release.sh
├── updateMain.sh
└── utils
    ├── packageList.go
    ├── tables.go
    └── terminal.go

```