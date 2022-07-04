# Overview
This is an instruction on how to setup a test environment for the ztp-dhcp server purely based on network namespaces.

```
#  Add a client and a relay namespace
ip netns add client
ip netns add relay

# add two veth pairs
# client <-> relay
# relay <-> host
ip link add cif type veth peer name rifc
ip link add hif type veth peer name rifh

# assign veth interfaces to namespaces
ip l set dev rifc netns relay
ip l set dev rifh netns relay
ip l set dev cif netns client

# assign IPs
ip netns exec relay ip a add dev rifc 192.168.50.1/24
ip netns exec relay ip a add dev rifh 192.168.51.10/24
ip a add dev hif 192.168.51.1/24

# bring up the interface
ip netns exec relay ip l set dev rifh up
ip netns exec relay ip l set dev rifc up
ip netns exec client ip l set dev cif up
ip l set dev hif up

# add a route to the client subnet on the host
ip r add 192.168.50.0/24 via 192.168.51.10

# add route to dhcp ip in relay namespace
ip netns exec relay ip r add default via 192.168.51.1

# execute the dhcp relay agent in the relay namespace
# (needs to listen on uplink as well as downlink interface)
ip netns exec relay dhcrelay -a -d -i rifc -i rifh 172.24.100.101

# MANUAL STEP HERE:
## start the ztp-dhcp on the host in the root namespace

# start the dhcpclient in the client namespace
ip netns exec client dhclient -d -lf /dev/null -i cif
```