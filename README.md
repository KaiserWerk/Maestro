# Maestro
A slim app to cover your service discovery needs.

* A service to register/deregister an app, do health pings and query for service addresses
* Every party uses the same auth token
* Can be run with HTTP or HTTPS
* If an app does not send a ping within a configurable interval, it is considered dead
and removed from the registry

Maestro comes as a single binary for the popular operating system and has no dependencies.
It just needs an open TCP port to be reachable.

## Configuration

```yaml
bind_address: 'http://localhost:9200'
auth_token: 123abc # change this to something secure
die_after: 5m # supply time units, e.g. 2h12m52s. default is 5 minutes
certificate_file:
key_file:
```

The bind address specifies which host, if any, and which port Maestro binds to.
When the bind address is actually ``http://localhost:9200``, Maestro is only reachable
via localhost. If you want to bind any IP address, use ``http://:9200``.
If you also want to use TLS, supply the paths to a certificate and key file in PEM
format and use the bind address ``https://:9200``.

All configuration values can be overwritten using environment variables; they have higher
precedence than configuration file values, if they are set.

Available environment variables are

* MAESTRO_BIND_ADDRESS
* MAESTRO_AUTH_TOKEN
* MAESTRO_DIE_AFTER
* MAESTRO_CERTIFICATE_FILE
* MAESTRO_KEY_FILE

## API

The API route prefix is ``/api/v1``.

You can also use the [Maestro Go SDK](https://github.com/KaiserWerk/Maestro-Go-SDK) and the [Maestro .NET SDK](https://github.com/KaiserWerk/Maestro-DotNet-SDK).

Use the custom ``X-Registry-Token`` header to supply the authentication token.

#### Registration

``POST /register``

with the request body

```json
{
    "id": "some-service-handle",
    "address": "http://localhost:9001"
}
```

registers an app with the given ID under the specified address which is then queryable
by other parties.

#### Deregistration

``DELETE /deregister?id=some-service-id``

removes the entry with the supplied service ID from the registry.

#### Query

``GET /query?id=some-service-id``

queries the registry for the entry with the supplied ID.
The response body looks like this:

```json
{
    "id": "some-service-id",
    "address": "http://localhost:9001",
    "last_ping": "2021-12-05T23:54:07.844640791+01:00"
}
```

#### Ping

``PUT /ping?id=some-service-id``

signals that the app which sends out the ping is still alive and working.
