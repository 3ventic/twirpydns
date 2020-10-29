# TwirpyDNS

This is an implementation of DNS over Twirp. Point clients at the client, configure server to forward the request to desired upstream.

## Usage

1. Download a release for your platform or build using `go build {client,server}/main.go`
2. Create server.ini and/or client.ini using the example files provided in the source client and server folders.

### Server

3. Run the server binary and expose the port or use your favorite reverse proxy to expose the twirp endpoint.

### Client

3. Run the client binary on the desired computer and set your system DNS to 127.0.0.1:53.

Tip: on Windows you can install the binary as a service using [nssm](https://nssm.cc/)

## Contributing

To update the protobuf schema, run `protoc --proto_path=. --twirp_out=. --go_out=. ./rpc/twirpydns/twirpydns.proto`. This is only required when making changes to the proto file. You can install protoc from [its repository](https://github.com/protocolbuffers/protobuf/releases).

Added configurability, testing, and CI is welcome as long as the codebase remains simple. KISS applies here.
