[![Build Status](https://magnum.travis-ci.com/JNPRAutomate/gogoDoS.svg?token=Taq81d9PL7keqp96e9qu&branch=master)](https://magnum.travis-ci.com/JNPRAutomate/gogoDoS)

#gogoDoS

This tool helps in testing DDoS/DoS prevention techniques. While this tool could support multiplatform as of today it is Linux only.

##Configuring Linux
When generating an extreme amount of new packets the kernel needs a bit of tuning. If the kernel is not tuned then it will limit the max PPS of the tool. There are two files that need to be tuned.

### Edit /etc/sysctl.conf

At the end of the file place
```
fs.file-max = 900000
```

### Edit /etc/security/limits.conf

At the end of the file place the following entries. This will allow root to open many files. Since the tool works in raw packets we only allow root to open the filehandles.
```
root          soft     nofile         900000
root          hard     nofile         900000
```

### After file changes

Reboot the device to take effect. To ensure the limits are correctly imposed check with the following commands.
```
rcameron-mbp15:~ rcameron$  sysctl -a | grep fs.file-max
fs.file-max = 900000

rcameron-mbp15:~ rcameron$ su - 
Password:
rcameron-mbp15:~ root# ulimit -a
core file size          (blocks, -c) 0
data seg size           (kbytes, -d) unlimited
file size               (blocks, -f) unlimited
max locked memory       (kbytes, -l) unlimited
max memory size         (kbytes, -m) unlimited
open files                      (-n) 900000
pipe size            (512 bytes, -p) 1
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) 709
virtual memory          (kbytes, -v) unlimited
```

##Building:
```
rcameron-mbp15:gogoDoS rcameron$ xport GOPATH=\`pwd\`
rcameron-mbp15:gogoDoS rcameron$ cd src
rcameron-mbp15:src rcameron$ go build gogoDoS.go
rcameron-mbp15:src rcameron$ ls -la gogoDoS
-rwxr-xr-x  1 rcameron  JNPR\prodmktg  6770972 Mar 31 16:06 gogoDoS
rcameron-mbp15:src rcameron$ 
```

##Running gogoDoS

```
Usage of ./gogoDoS:
  -D=60: Specify the total duration of the test
  -F=false: Specifies if the dns request should be flooded statelessly
  -P=53: Specify destination port for DoS
  -R=false: If set to true the specified then specify the source IPs to spoof the requests from. In this case the source IPs are destination IPs and the destination is the source.
  -d="127.0.0.1": Specify an single host or a list of destination hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5)
  -i="": Specify which interface name to eject packets from
  -p="dns": Specify protocol to use for DoS (dns)
  -r=1: Specify the amount of protocol requests per second
  -s="127.0.0.1": Specify an single host or a list of source hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5), only used for reflection attacks

rcameron-mbp15:~ rcameron$ ./gogoDoS -D 60 -P 53 -i eth0
 ````