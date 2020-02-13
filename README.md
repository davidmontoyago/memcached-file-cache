# FileCache

Uses a `ChunkedFile` to split a file into sized chunks between 96 bytes and 1024 Kbytes to prevent contention on a single memcached slab. The chunks sizing depends on the strategy passed as an implementation of `pkg/chunker.ChunkSizer`.

`pkg/chunker.ChunkSizer` allows using different strategies for determining a chunk size. Currently a random and the uniform distribution are implemented with the uniform as the default one in use.
Given that in memcached the last slab (>512 Kbytes) contains the majority of the distribution values between 96 bytes and 1024 Kbytes, it tends to contain an overallocation of values. A more elaborate `ChunkSizer` implementation could take this into consideration but it would come at the expense of more file fragmentation.

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

# optional - check slab allocation
watch docker run --rm --network=host koudaiii/memcached-tool localhost:11211 display

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
