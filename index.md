---
layout: default
---

```bash
$ git clone https://github.com/syui/medigo 
$ cd medigo
$ go build -ldflags "-X main.cid=$client_id -X main.secret=$secret_id"
$ ./medigo h
```

