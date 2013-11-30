go-swift-get-nodes
==================

swift-get-nodes...in :sparkles:Go:sparkles: ...hence the name

#### Install 
  
```
go get github.com/pandemicsyn/go-swift-get-nodes
```
  
#### Running it

Usage:
```
fhines@47:~$ go-swift-get-nodes 
Usage: 
	 go-swift-get-nodes /path/to/ring.gz /account/[container]/[object]
```

Success:
```
fhines@47.ronin.io:~$ go-swift-get-nodes /etc/swift/object.ring.gz /AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/ff.jpg
Ring File:	/etc/swift/object.ring.gz
Target:		/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/ff.jpg
Partion:	113446
Hash:		6ec9aa9170d6858cc8d53a7b220bf168

Node 1: 200 OK
Url: http://127.0.0.1:6010/sdb1/113446/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/ff.jpg
 X-Timestamp: [1382728047.59973]
 Last-Modified: [Fri, 25 Oct 2013 19:07:27 GMT]
 Content-Length: [384845]
 Etag: ["1d17f9281079a3fb59593b4e76918d86"]
 Content-Type: [image/jpeg]
 Date: [Sat, 30 Nov 2013 06:16:20 GMT]

Node 2: 200 OK
Url: http://127.0.0.1:6030/sdb3/113446/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/ff.jpg
 X-Timestamp: [1382728047.59973]
 Last-Modified: [Fri, 25 Oct 2013 19:07:27 GMT]
 Content-Length: [384845]
 Etag: ["1d17f9281079a3fb59593b4e76918d86"]
 Content-Type: [image/jpeg]
 Date: [Sat, 30 Nov 2013 06:16:20 GMT]

Node 3: 200 OK
Url: http://127.0.0.1:6020/sdb2/113446/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/ff.jpg
 X-Timestamp: [1382728047.59973]
 Last-Modified: [Fri, 25 Oct 2013 19:07:27 GMT]
 Content-Length: [384845]
 Etag: ["1d17f9281079a3fb59593b4e76918d86"]
 Content-Type: [image/jpeg]
 Date: [Sat, 30 Nov 2013 06:16:20 GMT]
```

A failure:
```
fhines@47:~/go/src/github.com/pandemicsyn$ go-swift-get-nodes /etc/swift/object.ring.gz /AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/cats.gif
Ring File:	/etc/swift/object.ring.gz
Target:		/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/cats.gif
Partion:	202558
Hash:		c5cf8f74abf7616c381a49578e1a0e8e

Node 1: 404 Not Found
Url: http://127.0.0.1:6040/sdb4/202558/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/cats.gif
 Content-Type: [text/html; charset=UTF-8]
 Content-Length: [0]
 Date: [Sat, 30 Nov 2013 06:16:32 GMT]

Node 2: 404 Not Found
Url: http://127.0.0.1:6030/sdb3/202558/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/cats.gif
 Content-Type: [text/html; charset=UTF-8]
 Content-Length: [0]
 Date: [Sat, 30 Nov 2013 06:16:32 GMT]

Node 3: 404 Not Found
Url: http://127.0.0.1:6010/sdb1/202558/AUTH_f1317546-589c-4ab2-a136-bf6a658b19f6/omgmonkey/cats.gif
 Content-Type: [text/html; charset=UTF-8]
 Content-Length: [0]
 Date: [Sat, 30 Nov 2013 06:16:32 GMT]
（╯°□°）╯ ~:74@sǝuıɥɟ$
```

