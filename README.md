# zstdzip
compress and decompress via ZSTD

## Usage:
compress:

`
./zstdzip zip --input=/User/harryzhu/docs  --output=/User/harryzhu/docs.zip   --speed=6 --logstatus=/User/harryzhu/result
`

--speed= can be 0/1/6/9

decompress:

`
./zstdzip unzip --input=/User/harryzhu/test.zip  --output=/User/harryzhu/test
`
or you can use https://github.com/mcmilk/7-Zip-zstd to unzip

## Benchmark:
