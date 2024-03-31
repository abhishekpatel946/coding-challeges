
### Valid request for getting 200 response

```sh
$ nc localhost 3000
    GET / HTTP/1.1 \
    Host: localhost:3000

    HTTP/1.1 200 OK
```

### Invalid request for getting 400 response

```sh
$ nc localhost 3000
    GET /index.html HTTP/1.1 \
    Host: localhost:3000

    HTTP/1.1 400 Not Found
```
