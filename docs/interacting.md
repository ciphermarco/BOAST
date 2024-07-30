# Interacting

Interacting with the server is quite simple. There's a small example bash client to show
how simple it is on
[../examples/bash_client/client.sh](https://github.com/ciphermarco/boast/blob/master/examples/bash_client/client.sh).
Don't mind its potential lack of elegance though; it was only made to show how little
preparation you need to interact with the server.

All you have to do is send a GET request containing an `Authorization` header with the
value `Secret <your base64 secret>` to the HTTPS API's `/events` endpoint. And that's
it. You'll receive a test `id` which can be used as an unique domain for your payloads
(e.g. `<id>.example.com`). And when you wish to retrieve new possibly existing events,
you just need to do the same to check if the `events` array was updated.

Of course, the experience is intended to be better with a client such as
[ZAP](https://github.com/zaproxy/zaproxy/issues/3022) (when/if it comes to be supported)
or at least a script better than the example bash client.

Now let's see a terminal interaction example step-by-step so it's clear. This
walkthrough will assume a server controlled by a third-party where you don't have
control over the events' TTL and other parameters.

## Registration

### 1. Generate a base64 random secret

You can generate it as you wish. The only limitations are that it must be a valid base64
and not longer than 44 bytes decoded.

```
$ openssl rand -base64 32
kkMrhv3ic2Em63PH6duIejNVRiqyOYpfBZHkjTDswBk=
```

### 2. Register

```
$ curl -H "Authorization: Secret kkMrhv3ic2Em63PH6duIejNVRiqyOYpfBZHkjTDswBk=" https://example.com:2096/events
{"id":"cxcjyaf5wahkidrp2zvhxe6ola","canary":"x7ilthx62hx2kfyvsioydd43da","events":[]}
```

This will give you a test `id` that can be used on your tests and a `canary` token that
can be used by protocol receivers when responding to the target application to aid some
kinds of test results detection. The most useful form to use it is to send
`cxcjyaf5wahkidrp2zvhxe6ola.example.com` in your payloads. There are other reasons for
this, but a good one is that, by using this unique domain in your payloads, even if the
target application is behind some firewall or whatnot and other protocol communications
fail (e.g. an HTTP request is blocked), you could be lucky (not uncommon) and receive
the DNS query for this unique domain in case the outbound restrictions do not apply to
the DNS protocol and port.

Both the `id` and `canary` are deterministically generated from the sent secret using a
cryptographic hash function. This means that if the server maintains the `hmac_key`
configuration parameter, you can expect to always receive the same `id` and `canary` for
a given secret. This is nice to have because, even if the server is restarted, your
previously sent payloads will still be valid and, after repeating this step 2, the
server will be able to recognise any late reactions from the target application. This is
one of the most important reasons for why these values are deterministically generated.

Note that, depending on the server's configuration, the API domain could be different
from the protocol receivers' domain.

### 3. Generate an event (and optionally check it's populating `events`)

```
$ curl http://example.com/cxcjyaf5wahkidrp2zvhxe6ola
<html><body>x7ilthx62hx2kfyvsioydd43da</body></html>
```

This steps does more than just to check if the server is working. It is advisable to do
this periodically because tests without any events will be deleted according to the
server's
[configured](https://github.com/ciphermarco/boast/blob/master/docs/boast-configuration.md)
`checkInterval` which may be too soon for your test reactions to happen and be recorded
by the server. This is a limitation you don't have to worry if you control the server as
you can change the configuration parameters to best suit your needs, but it is important
to have this in mind in the case of using a third-party server.

## Retrieving events

For retrieving events, you only need to repeat step 2 from the last section. And if some
new event has been generated and not expired yet, the `events` array will be populated
with it.

```
% curl -k -H "Authorization: Secret kkMrhv3ic2Em63PH6duIejNVRiqyOYpfBZHkjTDswBk=" https://example.com:2096/events
{"id":"cxcjyaf5wahkidrp2zvhxe6ola","canary":"x7ilthx62hx2kfyvsioydd43da","events":[{"id":"fbb6osymic6llzuiw7f7ylwix4","time":"2020-09-16T16:31:05.183124969+01:00","testID":"cxcjyaf5wahkidrp2zvhxe6ola","receiver":"HTTP","remoteAddress":"127.0.0.1:57770","dump":"GET /cxcjyaf5wahkidrp2zvhxe6ola HTTP/1.1\r\nHost: localhost:8080\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\nAccept-Language: en-GB,en;q=0.5\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0\r\n\r\n"}]}
```
