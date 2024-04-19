# Deployment
This section explains how to deploy ProxyGuard.

## Systemd
ProxyGuard comes with various systemd files to make deployment easier, these
are listed in the [systemd](./systemd) directory. These systemd files are also
installed with the DEB and RPM packages. On a server with systemd, we recommend
to use these files as they use certain hardening settings such as: dynamic user,
namespace restriction, no new privileges. For the full list of settings, see the
systemd unit files.

## Reverse proxy
As ProxyGuard uses the HTTP Upgrade mechanism for the handshake, you can put
ProxyGuard behind a reverse proxy in a similar way you would proxy WebSockets
traffic. We will document the Apache way.

## Apache

To proxy ProxyGuard behind Apache, you need the `proxy_http` module enabled:

```bash
sudo a2enmod proxy_http
```

Once that is enabled you can restart Apache (`systemctl restart httpd` for
Fedora/EL and `systemctl restart apache2` on Debian/Ubuntu).

Proxying the traffic is then done with the following line inside of your VirtualHost configuration:

```apache
ProxyPassMatch "^/proxyguard$" "http://127.0.0.1:51820/" upgrade=UoTLV/1
```

In this example, the HTTP path is /proxyguard and ProxyGuard itself listens on localhost port 51820.
