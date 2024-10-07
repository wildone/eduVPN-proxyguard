## Client

**NOTE**: you do NOT need to read this documentation when you are using the 
latest eduVPN or Let's Connect! applications! Those VPN clients handle 
WireGuard over TCP automatically. This documentation is _only_ for manual 
WireGuard over TCP configuration and was used extensively for testing the 
WireGuard over TCP server.

### Linux

In order to use WireGuard over TCP, you need to modify the WireGuard 
configuration on the client and run `proxyguard-client`.

Install `proxyguard-client`, available as source code 
[here](https://codeberg.org/eduVPN/proxyguard), or available from the (server) 
[Repository](https://docs.eduvpn.org/server/v3/repo.html) for Fedora / EL and 
Debian / Ubuntu. Configure the repository first, and then continue below.

#### Fedora / EL

```bash
$ sudo dnf -y install proxyguard-client
```

Enable and start the service:

```bash
$ sudo systemctl enable --now proxyguard-client
```

#### Debian / Ubuntu

```bash
$ sudo apt -y install proxyguard-client
```

On Debian and Ubuntu the service is started and enabled automatically.

#### Configuration

Download a VPN configuration from your VPN server that supports WireGuard over
TCP. Make sure you select the "WireGuard (TCP)" variant of the profile. Store 
the configuration in `vpn.conf`.

Now run the `vpn_setup.sh` script from this folder:

```bash
$ sh ./vpn_setup.sh vpn.conf
```

This should set everything up!
