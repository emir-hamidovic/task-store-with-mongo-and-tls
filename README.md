# Task server with TLS (mongo/in memory implementations)

This uses mTLS so you need server and client PEM (key and cert) files.

There are three ways to generate server and client PEM files:
1. Create a Golang app to generate those for you similar to: [tls-self-signed-cert.go](https://github.com/eliben/code-for-blog/blob/master/2021/tls/tls-self-signed-cert.go) (run on client and server as well)
2. Run: `go run /usr/local/go/src/crypto/tls/generate_cert.go` (had issues with this on the client side)
3. Install [mkcert](https://github.com/FiloSottile/mkcert) and use that with `mkcert localhost` (for client and server)

So, example of how to run this server would be:
`go run rest.go -cert localhost.pem -key localhost-key.pem -clientcert ../restclient/clientcert.pem`

`clientcert` flag presents the location where I had the client certificate
`cert` and `key` are server's certificate and private key generated by mkcert

Using this client: [emir.hamidovic/restclient](https://github.com/emir-hamidovic/task-store-with-mongo-and-tls-client) this server can be tested. So far, only the simplest /task endpoint is available with *POST*, *GET*, *DELETE* functionalities. Also, it's possible to use *in memory* and *mongodb* implementations of the server.

When this server is running, to access it with a browser (for example *Chrome*), this can be used: [MTLS/Mutual TLS Authentication Chrome](https://velmuruganv.wordpress.com/2020/04/27/mtls-mutual-tls-authentication-chrome/) (certificates and keys mentioned here would be the ones from the client or in my case from `../restclient`)

As far as the other files, due to the implementation of this mTLS, `manual.sh` (which was used for testing before), is no longer valid. `middleware` package has one function for panic recovery which is not used.
