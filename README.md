# FileCache

Uses a `ChunkedFile` to split a file into randomly sized chunks between 96 bytes and 1024 Kbytes to prevent contention on a single memcached slab.

`pkg/chunker/rand.chunkSizer` could be extended with other strategies other than random sizing for determining a chunk size.

Uses an MD5 hash of the file's content as a file identifier and for later content verification. With a LOT of data this could cause collisions. A potential enhancement could be to use a hashing function more suited for uniqueness.

## Getting Started

### pre-reqs

- make
- Docker

```sh
# run tests
make test

# setup memcached
make memcached

# put a file via CLI
go run cmd/main.go put -f path-to-file

# get a file via CLI
go run cmd/main.go get -k file-key

# run API server
go run main.go

# PUT a file via API
curl -vvv -XPOST http://localhost:8080/filecache --upload-file ./file.dat

# GET a file via API
curl -vvv http://localhost:8080/filecache/91388263e7c545ebea3952fb2637dffa --output file.dat

# destroy memcached
make teardown-memcached
```
