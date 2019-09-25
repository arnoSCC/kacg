Keepalived Configuration Generator
==================================

Use case
--------

This simple go script create a simple Keepalived config based on a desired VIP value.

It first parse all interfaces to find the NIC with the correct CIDR for the VIP then output a config in Stdout using:

- The NIC found
- The VIP with a /32
- A priority computed from multiplying the last two byte of the IP, modulo 255.
- A router ID from the last byte of the VIP.

Demo
----

```
# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:d9:26:a2 brd ff:ff:ff:ff:ff:ff
    inet 192.168.71.216/24 brd 192.168.71.255 scope global noprefixroute dynamic eth0
       valid_lft 1639sec preferred_lft 1639sec
    inet6 fe80::c8b3:575d:5647:6cd0/64 scope link noprefixroute
       valid_lft forever preferred_lft forever

# ./kacg 10.0.0.1
No valid interface found for 10.0.0.1

# ./kacg 192.168.71.2
! Configuration File for keepalived
global_defs {
   vrrp_version 2
}
vrrp_instance VI_2 {
    state BACKUP
    interface eth0
    virtual_router_id 2
    priority 207
    advert_int 1
    nopreempt
    virtual_ipaddress {
        192.168.71.2/32
    }
}
```

Note
----

This script was made for a very specific use case and will likely not be maintained. It is also made only for IPv4.
