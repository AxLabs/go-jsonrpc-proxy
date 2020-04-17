# go-jsonrpc-proxy

Simple JSON-RPC proxy based on RPC methods.

## Use Case

Imagine the following set-up:

- Node1: serves many JSON-RPC methods on the running application, but 
only `X` method should be publicly exposed.
- Node2: serves specific JSON-RPC methods in the running application,
but only `Y` and `Z` methods should be publicly exposed.
- Node3: should serve all JSON-RPC methods publicly, even if some overlap 
with Node1 or Node2.

In this scenario it would be ideal that there's a proxy which could 
receive a request and forward it based on rules:
- If a request with method `X` comes, forward to `Node1`
- If a request with method `Y` comes, forward to `Node2`
- If a request with method `Z` comes, forward to `Node2`
- Else, forward to `Node3`

That's exactly what `go-jsonrpc-proxy` solves. :smiley:

## ToDos

- [ ] Rate limits for each declared method (using Redis)
- [ ] Rate limits based on API keys (e.g., API key `X` 
specified in the `Authorization` HTTP header can perform more 
requests than API key `Y`)
- [ ] More load balancer strategies. Nowadays forwarding hosts 
are randomly chosen. 
- [ ] `docker-compose` file
- [ ] Integration tests