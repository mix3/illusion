# illusion

[mirage](https://github.com/acidlemon/mirage) にインスパイアされて出来た、Dockerでコンテナを立ち上げるとコンテナ名で動的にリバースプロキシしてくれる君。

## Usage

config.toml を用意して
```toml
domain           = "example.net"
listen_addr      = "127.0.0.1:8080"
forward_port     = 5000
ignore_subdomain = []
docker_endpoint  = "unix:///var/run/docker.sock"
```

```
$ go get github.com/mix3/illusion
$ illusion
```

などとして起動する

## config.toml

### domain

待ち受けるドメイン名を指定する。サブドメインの判定にも使っている。

### listen\_addr

Webサーバの ホスト,IP:ポート を指定する

### forward\_port

Dockerコンテナへリバースプロキシするときのポートを指定する

コンテナ側はここで指定されたポートで待ち受けてもらう

### ignore\_subdomain

illusionをDockerで起動した場合にループするのを回避するため、サブドメインのマッチング対象外を指定可能にしている。

### docker\_endpoint

コンテナからホストのDockerAPIに触れる必要があるのでエンドポイントを指定する。

## LICENSE

MIT
