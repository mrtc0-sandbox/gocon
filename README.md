# gocon

---

```bash
docker export $(docker create alpine /bin/sh) > alpine.tar
tar xf alpine.tar
go run main.go /bin/sh
```
