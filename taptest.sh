#!/bin/bash
# Build binary.
go test -c -tags=taptest -o=dgrams.test .

# Set CAP_NET_ADMIN capabilities.
sudo setcap cap_net_admin=+ep ./dgrams.test
echo "start tests"
./dgrams.test & # Send test job to background and continue linking tun interface.
pid=$! # Get PID of our test.

DEV=tap0
# Set tunnel's IP:
sudo ip addr add 192.168.0.2/24 dev $DEV # ip addr will now show our tun0 interface after this command.

sudo ip link set up dev $DEV # This links our tun0 with another interface causing it to read data.
ping -I $DEV 192.168.0.2:9090 & # Start pinging on the tun0 device.

trap "kill $pid" INT TERM
wait $pid