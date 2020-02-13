# FileCache

Splits a `ChunkedFile` in randomly sized chunks between 96bytes and 1024Kbytes to prevent contention on a single memcached slab. `pkg/chunker/rand.chunkSizer` could be extended with other approaches for chunking the file.

## Getting Started

### pre-reqs

- make
- Docker

```sh
# setup memcached
make memcached

# run tests
make test

# put a file via CLI
go run cmd/main.go put -f path-to-file

# get a file via CLI
go run cmd/main.go get -k file-key

# run API server
go run main.go

# PUT a file
curl -vvv -XPUT http://localhost:8080/filecache --upload-file ./file.dat

# GET a file
curl -vvv http://localhost:8080/filecache/91388263e7c545ebea3952fb2637dffa --output file.dat

# destroy memcached
make teardown-memcached
```

# TODO 
- Make chunkSizer an interface to swap strategies for chunk sizing
- MD5 based keys with a LOT of data could cause collisions - add a timestamp to key