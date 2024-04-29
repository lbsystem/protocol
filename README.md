# Decode net packet
* It's a super fast decode program.<br>
* It's ZeroCopy data you shall be handle this.<br>
---
- [X] VLAN
- [X] ICMP
- [X] ICMPv6
- [X] IPv6
- [X] IPv4
- [X] ARP
- [X] UDP
- [X] TCP
- [ ] DHCP

### example
```go
eth := protocol.NewEthernet()
eth.UnmarshalBinary(data)
ip,ok:=eth.Data(*protocol.IPv4)
if ok{
    udp,ok:=ip.Data(*porotocol.UDP)
    if ok{
        udp.anyFeild
        udp.anyMod
    }
}

```