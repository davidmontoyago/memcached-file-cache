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