# BOAST configuration

## Flags

* `-config` | _(string)_ | TOML configuration file (default "boast.toml")
* `-dns_only` | Run only the DNS receiver and its dependencies
* `-dns_txt` | DNS receiver's TXT record
* `-log_file` | _(string)_ | Path to log file
* `-log_level` | _(int)_ | Set the logging level (0=DEBUG|1=INFO) (default 1)
* `-v` | Print program version and quit

## Configuration file

By default, BOAST will look for a file called `boast.toml` in the working directory but this behaviour can be changed with the `-config` flag.
Example configuration files may be found in [the config directory](https://github.com/ciphermarco/boast/tree/master/examples/config).

Here's a brief description of each configuration section and its parameters:

### Temporary storage

The `[storage]` section is required as the server is useless without it.

Only the `hmac_key` parameter is optional, but you should be aware of the implications. (TODO: explain the implications :])
  
* `[storage]`: Section for the temporary in-memory events storage.
  * `max_events` _(int)_ | The maximum number of events to be held by the server in a given moment | Example value: `1_000_000`
  * `max_events_by_test` _(int)_ | The maximum number of events by test | Example value: `"80KB"`
  * `hmac_key` _(string)_ | The HMAC key to be used by the server's HMAC algorithm | Example value: `"TJkhXnMqSqOaYDiTw7HsfQ=="`
  * `[storage.expire]`: Section for the storage's expiration feature.
    * `ttl` _(string)_ | Time to live for the stored events | Example value: `"24h"`
    * `check_interval` _(string)_ | Interval for checking and deleting expired events according to `ttl` | Example value: `"1h"`
    * `max_restarts` _(int)_ | Maximum attempts to restart the expiration routine before crashing | Example value: `100`

### API

The `[api]` section is required as the server is useless without it.

The `domain` parameter is optional and can be used for setting a different domain for
the API if needed (e.g. the API is behind a proxy for protection). If `domain` is not set, the API will be respond in any subdomain (e.g. anything.example.com).

The `[api.status]` subsection is optional and it will just deactivate the status page if not set.

* `[api]`: Section for the web API.
  * `domain` _(string)_ | The domain name for the API | Example value: `"proxied.example.com"`
  * `host` _(string)_ | The host for the API | Example value: `"0.0.0.0"`
  * `tls_port` _(int)_ | The TLS port for the API | Example value: `2096`
  * `tls_cert` _(string)_ | The TLS certificate file for the API | Example value: `"/path/to/tls/fullchain.pem"`
  * `tls_key` _(string)_ | The TLS private key file for the API | Example value: `"/path/to/tls/privkey.pem"`
  * `[api.status]`: section for the server's status page
    * `url_path` _(string)_ | The secret URL path for the satus page | Example value: `"rzaedgmqloivvw7v3lamu3tzvi"`
    
### HTTP receiver
    
The `[http_receiver]` is optional.

The `host` parameter is required if any of the `ports` parameters is set.

The `[http_receiver.tls]`'s `cert` and `key` are required if any TLS `ports` is set.

`real_ip_header` is always optional.

* `[http_receiver]`: Section for the HTTP protocol receiver.
  * `host` _(string)_ | The host for the HTTP receiver | Example value: `"0.0.0.0"`
  * `ports` _([]int)_ | The ports for the HTTP receiver | Example value: `[80, 8080]`
  * `real_ip_header` _(string)_ | The client's real IP header to be recorded when proxied | Example: `"X-Real-IP"`
  * `[http_receiver.tls]`: Section for the HTTP receiver's TLS configuration.
    * `ports` _([]int)_ | The TLS ports for the HTTP protocol receiver | Example value: `[443, 8443]`
    * `cert` _(string)_ | The TLS certificate file for the HTTP protocol receiver | Example value: `"/path/to/tls/fullchain.pem"`
    * `key` _(string)_ | The TLS private key file for the HTTP protocol receiver | Example value: `"/path/to/tls/privkey.pem"`
    
### DNS receiver

The `[dns_receiver]` is optional.

If the `ports` parameter is set, the only optional parameter is the `txt`. All the other
parameters are required for the correct functioning.
    
* `[dns_receiver]`: Section for the DNS protocol receiver.
  * `host` _(string)_ | The host for the DNS receiver | Example value: `"0.0.0.0"`
  * `ports` _([]int)_ | The ports for the DNS receiver | Example value: `[53]`
  * `domain` _(string)_ | The domain name for the server | Example value: `"example.com"`
  * `public_ip` _(string)_ | The server's publicly accessible IP | Example value: `"203.0.113.77`
  * `txt` _([]string)_ | An arbitrary TXT DNS record | Example value: `["testing", "TXT"]`
