# FileCache

Uses a `ChunkedFile` to split a file into sized chunks between 96 bytes and 1024 Kbytes to prevent contention on a single memcached slab. The chunks sizing depends on the strategy passed as an implementation of `pkg/chunker.ChunkSizer`.

`pkg/chunker.ChunkSizer` allows using different strategies for determining a chunk size. Currently a random distribution, a uniform, a skewed distribution sizer, and a cumulative probability distribution by slab are implemented.

Given that in memcached the last slab (>512 Kbytes) contains the majority of the distribution values between 96 bytes and 1024 Kbytes, a random or a uniform distribution tend to cause an overallocation of values. 

With the skewed chunk sizer, 80% of the time we allocate chunks between 96bytes and 246KB, and the rest of the time, chunks are allocated to the more spacious slabs after 246KB. This helps to even out the contention across the bigger slabs.

With the cumulative probability distribution by slab, we randomly select a slab range and generate a random chunk size between the lower and upper bounds of the slab size.

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
