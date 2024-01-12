
This is a consistency checker for Storagnodes of the Storj network.

It reads all the blob files (...sj1 files), and recalculates the checksum to compare it with the stored one.

Run it from the storage dir (the dir where you see at least one blobs subfolder).

Install:

```
github.com/elek/storagenode-checker@latest
```

Example:

```
storagenode-checker /storj/storj01/data/storage
```

Example of rotten data:

```
cd /storj/storj01/data/storage
storagenode-checker

checking namespace  7b2de9d72c2e935f1918c058caaf8ed00f0581639008707317ff1bd000000000
Hash mismatch: d6b5d209213990998978e2e3c8d455cbc394a069092a88b822a6977686ba0f0b
Hash mismatch: d6b5d509f9c818936b72fd93c30dbecf2edf1744f42d71905bdb8ac594dca6fc
Hash mismatch: d6b5e4a6eca64fdd0cd2ad2d45dcbfb5d7e9745259775763e456e2ac620f9d6c
```
