# UNRELEASED

- Client: Ensure the hostname is used for a DNS request instead of the host:port
- Server: Do not spam log by not logging EOF, TCP reader or clean exits
- General: Rename Proxyguard to ProxyGuard

# 0.4.0

- Client: Make it a struct type for a nicer Go API
- Client: Set the default source port to 0
- Client + Server: Make the TCP reader timeout after 60 seconds
- Client: Ability to re-use source port on Linux, for reconnecting purposes
- Client: Do a DNS request inside Proxyguard if no --peer-ips are set and cache it so that reconnecting cannot fail

# 0.3.0 (2024-02-22)

- Client: Implement a --peer-ips flag to bypass DNS resolution
- Client: Log each IP that is being connected to in the --peer-ip case
- Client: Add a callback when the proxy is ready

# 0.2.0 (2024-02-13)

- Add a file descriptor callback mechanism to the client for eduvpn-common + Android integration
- Add a basic reconnect/retry mechanism to the client (fixes: #15)
- Update systemd files for the HTTP change
- Proxy over HTTP such that you can use it with a reverse proxy, e.g. Apache. This is done similar to Websockets with HTTP upgrade requests
- Add systemd files for Proxyguard client mode. Make sure it starts after the network is up

# 0.1.0 (2024-01-29)

- initial release of Proxyguard. A Go program and library that tunnels UDP over TCP, to use with WireGuard
