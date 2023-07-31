# envoy-go-basic-auth

This is a simple basic auth filter for envoy written in go. Only requests that pass the configuration authentication will be proxied to the upstream service.

## Status

The plugin has been completed and may be updated in the future.

## Usage

The client sets credentials in `Authorization` header in the following format:

```Plaintext
credentials := Basic base64(username:password)
```

An example of the `Authorization` header is as follows (`Zm9vOmJhcg==`, which is the base64-encoded value of `foo:bar`):

```Plaintext
Authorization: Basic Zm9vOmJhcg==
```

Configure your [envoy.yaml](envoy.yaml) to set pairs of username and password.

```yaml
http_filters:
- name: envoy.filters.http.golang
typed_config:
  "@type": type.googleapis.com/envoy.extensions.filters.http.golang.v3alpha.Config
  library_id: example
  library_path: /etc/envoy/libgolang.so
  plugin_name: basic-auth
  plugin_config:
    "@type": type.googleapis.com/xds.type.v3.TypedStruct
    value:
      users:
        - username: "foo"
          password: "bar"
        - username: "lobby"
          password: "niu"
```

Then, you can start your filter.

```bash
make build
make run 
```

## Test

This test case is based on a local Envoy. Run it with the example config file.

```bash
make test
```
