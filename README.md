# dyndnsdo

A basic application that I can run on a virtual machine to perform dynamic DNS
updates from my [UniFi Security Gateway](https://www.ui.com/unifi-routing/usg/).

## How it works?

Using the `namecheap` backend in the USG, I can pass values and a custom endpoint
address where the `ddclient` utility will periodically submit.

The `/update` endpoint is appended to the server parameter and is requested by
HTTPS. It is passed four form values, encoded as the query string to a `GET`
request:

1. host
2. domain
3. password
4. ip

I want to update the naked domain name record, so we pass the domain name as the
`domain` value, `@` as the `host`, our target IP as identified by the USG is `ip`
and `password` is a pre-shared secret.

## Certificates

The application will require the `-cert` and `-key` switches to allow the service
to make use of TLS.

## Socket activation

I wanted to experiment with a few technologies. Using Go was first, but socket
activation means that the application will be spawned when traffic is received.

```
# /etc/systemd/system/dyndns.socket
[Socket]
ListenStream=443

[Install]
WantedBy=sockets.target
```

```
# /etc/systemd/system/dyndns.service
[Service]
Environment=DO_API_TOKEN=TOP_SECREY
ExecStart=/usr/local/bin/dyndnsdo -cert /etc/pki/tls/private/cert1.pem -key /etc/pki/tls/private/privkey1.pem -password sharedsecret
```