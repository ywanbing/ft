# ft
big file transfer, support various network protocols

### Use

```text
NAME:
   ft - big file transfer, support various network protocols

USAGE:
   ft [global options] command [command options] [arguments...]

COMMANDS:
   client, cli  start an upload client.
   server, srv  start a server that receives files and listens on a specified port.
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### example

启动服务
```text
./ft srv|server -d 保存文件的路径 --addr 服务监听的地址 --nw 可选的网络协议


NAME:
   ft server - start a server that receives files and listens on a specified port.

USAGE:
   ft server [command options] [arguments...]

OPTIONS:
   --addr value                 specify a listening port (default: "0.0.0.0:9988")
   --dir value, -d value        upload dir or save dir (default: "./data")
   --network value, --nw value  choose a network protocol(tcp|udp) (default: "tcp")
```


启动客户端
```text
./ft cli|client -d 文件所在的文件夹 --addr 服务器地址 --nw 可选的网络协议  [需要传输的文件名,可以多个]

NAME:
   ft client - start an upload client.

USAGE:
   ft client [command options] [arguments...]

OPTIONS:
   --addr value                 specify a server address (default: "127.0.0.1:9988")
   --dir value, -d value        upload dir or save dir (default: "./")
   --network value, --nw value  choose a network protocol(tcp|udp) (default: "tcp")

```
