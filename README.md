# Onetable

Simple persistent key-value store

## How does it work?
To use the `OneTable` you need to choose an index.
The index is an in memory data structure used for fast lookups
of of value metadata (offset and length).

- The values are stored in an append only file, which does not make
much sense without the index
- The index data is stored also in an append only csv file in 
format `{key: string},{offset: int},{length: int}`

This allows for fast lookups and inserts without loading the
entire file content to memory
