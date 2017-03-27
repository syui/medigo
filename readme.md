medium client

Download : [https://github.com/syui/medigo/releases](https://github.com/syui/medigo/releases)

See : [https://syui.github.io/medigo/](https://syui.github.io/medigo/)

```bash
$ git clone https://github.com/syui/medigo 
$ cd medigo
$ go build -ldflags "-X main.cid=$client_id -X main.secret=$secret_id"
$ ./medigo h
```

aur : https://aur.archlinux.org/packages/medigo/

```bash
$ yaourt -S medigo --noconfirm
```
