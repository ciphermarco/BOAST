# Deploying

## Makefile

The Makefile can be found
[here](https://github.com/ciphermarco/boast/blob/master/Makefile) and, together with the
Go module files, should make it very easy to build BOAST by yourself. You most likely
want to issue `make` for building and `make test` for running the package's tests.

## BOAST configuration

The server's configuration file is described on
[boast-configuration.md](https://github.com/ciphermarco/boast/blob/master/docs/boast-configuration.md)
and example configurations can be found in the [config examples
directory](https://github.com/ciphermarco/boast/tree/master/examples/config).

## Log level

The default log level is INFO (1) which must not disclose any details about the
reactions events. The log level can be changed to DEBUG (0) passing the `-log_level=0`
flag to the binary. I may implement this flag with their mnemonics instead of numbers
soon to make it more obvious, but this will always be a flag and never a parameter in
the configuration file or any other somewhat implicit way. The reason for this is that
avoiding the mistake of unintentionally logging possibly sensitive testing information
is paramount.

## Deploying with Docker

A Dockerfile, a BOAST configuration file (`boast.toml`), and `certbot` pre validation
and renew hooks can all be found [in the build
directory](https://github.com/ciphermarco/boast/tree/master/build). They are meant to
work together and you must edit some parameters in the `boast.toml` file. Additionally,
you may need to edit other files if you want a different setup.

Also note that the steps listed below may be followed with a variety of divergences
depending on your on preferences that will not be exhaustively detailed on this part of
the document since this aims to be a simple tutorial. In general, this only means that
some pieces like the exact DNS records may be configured in slightly different ways and
still be valid (or not even be used at all for less functionality), but the overall
process should remain very similar for any case.

This tutorial assumes you have cloned the repository and is at the project's root
directory.

### 1. DNS configuration

For full functionality, BOAST runs its own DNS server to respond and record queries
about the domain used for the protocol receivers. Thus, you have to dedicate an internal
or external domain or subdomain for this use. If your domain is `example.com`, you DNS
configuration should look something like this:

```
example.com.      IN      NS      ns1.example.com.
example.com.      IN      NS      ns2.example.com.
```

You also need to configure the glue-records for the NS domains. As these depend on how
your domain registrar exposes them on their interface, you should search their
documentation or contact support for more details.

### 2. Edit `boast.toml`

If you change any of the uncommented parameters, you may have to change one or more
parts of the remaining steps. For now, you only have to be careful about the ports as
this will change the port parameters for the `docker run` command.

Using the `boast.toml`, these are the values you need to uncomment and possibly change:

* `storage` section: `max_events`, `max_events_by_test`, `max_dump_size`, `hmac_key`.
* `storage.expire` subsection: `ttl`, `check_interval`, `max_restarts`. 
* `dns_receiver` section: `domain`, `public_ip`.

The other commented parameters are optional and may be changed at will.

For more details, have a look [at the configuration
section](https://github.com/ciphermarco/boast/blob/master/docs/boast-configuration.md).

### 3. Build the docker image and run BOAST with the `-dns_only` flag

```
$ docker build . -t boastimg -f build/Dockerfile
$ docker run -d --name boastdns -p 53:53/udp boastimg /go/src/github.com/ciphermarco/BOAST/boast -dns_only
```

This will build the BOAST's Docker image and run it in a container named `boastdns` with
the option flag `-dns_only`. As this option will only be used for the ACME DNS-01
challenge, there's no need to put this in the Dockerfile directly.

Having the DNS server running before the Let's Encrypt ACME DNS-01 challenge is
necessary or else it fails. The `-dns_only` flag is used so only the DNS receiver and
its dependencies are run and you don't have to worry about the TLS files not being in
place yet or anything else. Make sure the DNS receiver is configured to at least listen
on port 53 as this will be used for the ACME challenge as expected.

### 4. Wildcard TLS certificate

As BOAST will freely and dynamically use subdomains for its operations, it needs an
wildcard TLS certificate for the configured domain. You need to perform some variation
of this step even if you choose to not use the HTTPS receiver (by not configuring its
TLS ports) as the API only supports HTTPS. Of course, the certificate can be self-signed
or acquired by other means. The only requirement is that the TLS files must be PEM
encoded as [documented here](https://golang.org/pkg/crypto/tls/#LoadX509KeyPair).

To perform a Let's Encrypt ACME DNS-01 challenge to acquire a wildcard certificate, you
need [`certbot`](https://github.com/certbot/certbot) and a little help from BOAST to
respond to the challenge. Assuming the domain is `example.com` and the hook script has
execute permission, you may use this command:

```
$ certbot certonly --agree-tos --manual --preferred-challenges=dns -d *.example.com -d example.com --manual-auth-hook ./build/certbot-dns-01-pre-hook.sh
```

This command will attempt a wildcard certificate issuance from Let's Encrypt using the
provided script as a pre-validation hook.

An ACME DNS-01 challenge will be initiated and the validation string will be available
to the pre-validation hook script as `$CERTBOT_VALIDATION`. Any container named
`boastdns` or `boast` will be stopped and a new `-dns_only` BOAST container will be run
using the flag `-dns_txt` with the validation string as value. As the DNS receiver
responds the same TXT record for any subdomain, this will make sure that Let's Encrypt
will find the validation TXT record for the `_acme_challenge` subdomain with a record
similar to this:

```
_acme-challenge.example.com.     300     IN      TXT     "SIbmWivQZsP9RyoEb4KFjNYPywSU-YsDUlPtWc1xDWg"
```

If everything worked correctly and the validation was successful, the script output will
let you know with a congratulations message and information about your certificate files.

Now, copy the certificate files to a directory for this use so it can be reliably used
to mount a volume inside the container without the other files or problems with symlinks:

```
$ cp /etc/letsencrypt/live/example.com/fullchain.pem ./tls
$ cp /etc/letsencrypt/live/example.com/privkey.pem ./tls
```

### 5. Run

The only thing you need to do now is run a container with the exposed ports and a volume
containing the TLS files at the right container's path:

```
$ docker run -d --name boastmain -p 53:53/udp -p 80:80 -p 443:443 -p 2096:2096 -p 8080:8080 -p 8443:8443 \
		-v $PWD/tls:/go/src/github.com/ciphermarco/BOAST/tls boastimg 
```

And you can [start using it](https://github.com/ciphermarco/boast/blob/master/docs/interacting.md).

### 6. Automate the certificate renewal

This part of the documentation will be improved for better reproducibility, but, for
now, the renew hook script may need some editing to work on your end.  Make sure to test
it before delegating it to a `certbot` cron job.

For automating the certificate renewal process, you can use the `cerbot` pre validation
and renew hooks found in [the build
directory](https://github.com/ciphermarco/boast/tree/master/build). You just have to put
them in the right hook directories to be run by `certbot` when renewing or by using the
flags `--manual-auth-hook` and `--renew-hook` to run the hooks non-interactively like
this:

```
certbot certonly -n --agree-tos --manual-public-ip-logging-ok --manual --preferred-challenges=dns -d *.example.com  \
		 --manual-auth-hook $HOME/boast/build/certbot-dns-01-pre-validation-hook.sh \
		 --renew-hook $HOME/boast/build/certbot-dns-01-renew-hook.sh
```

Using the hook directories is recommended, but, when using the flags, make sure the
paths to the hooks are right.

In both cases, you just have to run the `certbot` command from a cron job with your
preferences and these two hooks. The hooks will only be run if a renewal is due so, with
attention to [Let's Encrypt's rate limits](https://letsencrypt.org/docs/rate-limits/),
it's safe to run the cron job routinely for automatic certification renewal and server
restart with the new certificate.

To make customization easier, here's the minimum operations the pre validation and renew
hooks or alternatives should do:

**Pre validation hook:**

1. Stop any conflicting BOAST containers or restart it without binding it to port 53.

2. Start a DNS-only BOAST container with the right validation TXT record.

**Renew hook:**

1. Stop any conflicting BOAST containers.

2. Start the main BOAST container with the new certificates accessible to BOAST.

### Using a different domain for the API

One possibility not yet covered by this document is to configure the API's `domain`
parameter in the [configuration
file](https://github.com/ciphermarco/boast/edit/master/docs/boast-configuration.md).
Doing this will allow you to protect the API with a proxy or what else may need a domain
not dedicated to BOAST's DNS receiver. It will not be possible to perform the ACME
DNS-01 challenge using the BOAST's DNS receiver as a helper and you'll need to configure
the API's TLS file paths, but you may issue a self-signed certificate for the API only
(hence without the ACME challenge) if that fits your requiremets.

### Possible improvements

1. This can be made more automated and reproducible by pushing the whole ACME DNS-01
   challenge validation to Docker with `certbot` (or alternative) included. This way,
   the challenge container can be orchestated to perform the whole challenge process
   without host dependencies and save certificates to a volume to be shared with the
   main BOAST container. But I'm yet to document it :).

2. Automate most of the process with an installation script.
