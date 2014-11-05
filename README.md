revproxy
========

Reverse Proxy for tcp layer, Use it to caculate network flow.

```
go get -v github.com/codeskyblue/revproxy
cd $GOPATH/src/github.com/codeskyblue/revproxy
go build
```

edit `_config.yml` to set proxies.

Run `./revproxy` will open a interaction cli.

```
Proxy engine started...
        127.0.0.1:9000 --> 10.0.0.1:4000
        127.0.0.1:8000 --> 10.0.0.1:4000

Welcome to reverse proxy console

ReverseProxy $ h
Usage:
exit              Exit program
help/h            Show help information
print/p           Show netstat info

ReverseProxy $ 
```

use `print` to see network flow.
