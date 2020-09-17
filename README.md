# home24

## Backend requirements

* [golang](https://golang.org/) - The Go Programming Language
* [golang mod](https://github.com/golang/go/wiki/Modules) - Go dependency management tool 

## Versions requirements
* golang **>=1.15.X**

### Setup Linux

```bash
git clone git@github.com:donutloop/home24.git
cd ./home24
go build -o home24 ./cmd/home24/main.go
```

#### example call for bin

```bash
./home24
```

#### example http call

```bash
curl -i --header "Content-Type: application/json" \
  --request POST \
  --data '{"website_url":"https://google.de"}' \
  http://localhost:8080/websitestats
```
