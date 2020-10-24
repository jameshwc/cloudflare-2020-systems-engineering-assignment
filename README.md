# Go Simple-HTTP

This repo implemented simple http client with builtin "net" and "crypto/tls" (for https request) library.
It also supports profile websites, you can use it to stress test websites with concurrency.

## Install

With ```make```:
```$ make ```

or build manually:
```$ go build -o simple-http```

## Usage

You can run demo with make:

```bash
$ make single # single http/https request to the default url
$ make profile # profile the default website with no concurrency
$ make stress # stress test the default website with concurrency
```

Or manually run:
```bash
$ ./simple-http -url $URL # single http/https-request, url must with protocol prefix
$ ./simple-http -url $URL -profile 100 # send 100 requests to the url
$ ./simple-http -url $URL -profile 100 -c 100 # send 100 requests to the url with 100 concurrency, so total requests are 100*100=10000
```
