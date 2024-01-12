# Local development

You need running ftp server on port `20021`. With the following credentials:

```
USERNAME=user
PASSWORD=password
```

If you don't have one, use `docker-compose.yaml` in `deployment/dev`.

To run tests, start the ftp server and run the following command:

```shell
go test -v ./...
```