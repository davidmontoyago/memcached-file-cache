# Getting Started

### pre-reqs

- make
- Docker

```sh
# setup memcached
make memcached

# destroy memcached
make teardown-memcached

# run tests
make test
```


# TODO 
- Make chunkSizer an interface to swap strategies for chunk sizing
- MD5 based keys with a LOT of data could cause collisions - add a timestamp to key