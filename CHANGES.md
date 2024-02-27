# 0.4.0

- Client: Make it a struct type for a nicer Go API
- Client: Set the default source port to 0

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
