# Go Simple-HTTP

This repo implemented simple http client with builtin "net" and "crypto/tls" (for https request) library.
It also supports profile websites, you can use it to stress test websites with concurrency.

## Install

With ```make```:

```$ make ```

Or build manually:

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

## Discovery

I seperate the project into two parts: one is the simple implementation of http client in http package; the other is profile & stress tests in profile, stat, stress package.

First, for http package, I use ```net``` builtin library to build a TCP connection. To support https connection, I use ```crypto/tls``` library as well. As this project's goal is to make requests to the endpoints we created in the General Assignment, I don't implement any methods except ```GET```.

To deal with a response, I parse raw bytes and follow the [HTTP/1.1 Protocol](https://tools.ietf.org/html/rfc2616#section-6.1) carefully. However, I don't handle the data carefully enough to avoid attacks such as [CRLF injection](https://www.acunetix.com/websitesecurity/crlf-injection/). But I think it's enough for most common sites. :)
I made the package support ```Transfer-Encoding: chunked``` response. I've tested it with some sites and it worked like a charm. However I think more tests are needed to find out bugs since the logic is prone to errors (in my opinion).
Speaking of tests, please forgive me for I didn't write tests in this project. I have spent too much time for this project to catch up with school work.

The second part is profile & stress tests. I use goroutines to make concurrent requests. I run ```go -race``` to check if there's any race condition, and it only showed that data race in reading statistics. I think it's fine considering that the worst case is showing incorrect statistics at a moment (and it will output correct statistics eventually).

By the way, to get median efficiently with a stream of statistics, I use ```running median``` technique which I learned in *LeetCode No.295. Find Median from Data Stream*, and it's interesting to see that advanced algorithm we learn can be used in real world. :)