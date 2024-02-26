# zstdzip
compress and decompress via ZSTD

## Usage:
### compress:

`
./zstdzip zip --input=/User/harryzhu/docs  --output=/User/harryzhu/docs.zst.zip
`

or:

`
./zstdzip zip --input=/User/harryzhu/docs  --output=/User/harryzhu/docs.zst.zip  --speed=6  --threads=16 --log=/User/harryzhu/result.log
`

`--speed=0/1/6/9` : 

        0: fastest without compression, 

        1: default compression,

        6: better compression,
        
        9: best but slowest compression. 


`--log=/path/of/logfile.log` : log the global result(json format) to a file for monitoring 

### decompress:

`
./zstdzip unzip --input=/User/harryzhu/test.zip  --output=/User/harryzhu/test
`

or:

`
./zstdzip unzip --input=/User/harryzhu/test.zip  --output=/User/harryzhu/test --async --threads=32
`

`--async` : better performance for too many small files decompression, default is `false`.

or you can download https://github.com/mcmilk/7-Zip-zstd to unzip


## Benchmark:
