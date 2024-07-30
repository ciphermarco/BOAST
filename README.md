# BOAST

**BOAST** is the **B**OAST **O**utpost for **A**ppSec **T**esting: a server designed to receive and report Out-of-Band Application Security Testing (OAST) reactions.

```
            ┌─────────────────────────┐ 
            |          BOAST          ◄──┐
          ┌─┤ (DNS, HTTP, HTTPS, ...) |  |     
          │ └─────────────────────────┘  │     
          │                              │     
Reactions │                              │ Reactions
          │                              │     
          │                              │     
          │                              │     
   ┌──────▼──────────┐   Payloads   ┌────┴────┐
   │ Testing client  ├──────────────► Target  │
   └─────────────────┘              └─────────┘
```

Some application security tests will only trigger out-of-band reactions from
the tested applications. These reactions will not be sent as a response to
the testing client and, due to their nature, will remain unseen when the
client is behind a NAT. To clearly observe these reactions, another component
is needed. This component must be freely reachable on the Internet and capable
of communicating using various protocols across multiple ports for maximum
impact. BOAST is that component.

BOAST features DNS, HTTP, and HTTPS protocol receivers, each supporting multiple
simultaneous ports. Implementing protocol receivers for new protocols or customising
existing ones to better suit your needs is almost as simple as implementing the protocol
interaction itself.

## Used By

BOAST is used by projects such as:

- [Zed Attack Proxy (ZAP)](https://www.zaproxy.org/)

## Documentation

https://github.com/ciphermarco/boast/tree/master/docs
