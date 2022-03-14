# Measure missing required headers - Envoy WASM filter

This is simple Envoy filter that measures missing headers in requests. This can be used for example to monitor how 
services abide microservice contract. 

## To test locally:
```bash
$ make build-example
$ make run-example
```

Try locally:

```bash
$ curl localhost:18000 -v -H "my-custom-header: foo" 

$ curl -s 'localhost:8001/stats/prometheus'| grep missing_required_headers
# TYPE envoy_missing_required_headers counter
envoy_missing_required_headers{header="x-request-id"} 0
envoy_missing_required_headers{header="x-site"} 1

```
