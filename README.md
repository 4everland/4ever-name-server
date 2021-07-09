# 4ever-Name-Server
CoreDNS plugin for authoritative server.  
**This project needs to paired with a Resource Record (RR) repository.**

## Syntax
```
fns {
  api-url URL
}
```
* ```api-url``` is the query URL of the Resource Record (RR).

## Examples
```
. {
  fns {
    api-url https://ns.example.com/resouce-records
  }
}  
```

## Compilation from Source
```shell    
sh build.sh
```
## Compilation with Docker

```shell
docker build -t fns .
```