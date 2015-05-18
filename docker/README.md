## docker build

```
$ docker build --rm -t mix3/illusion:latest .
```

## docker run

```
$ docker run -d --name illusion -v /var/run/docker.sock:/var/run/docker.sock -p 80:8080 -t mix3/illusion:latest
```
