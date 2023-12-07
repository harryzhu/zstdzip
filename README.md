# zstdzip
compress and decompress via ZSTD

## Usage:
compress:

`
./zstdzip zip --input=/User/harryzhu/docs  --output=/User/harryzhu/docs.zip
`

or:

`
./zstdzip zip --input=/User/harryzhu/docs  --output=/User/harryzhu/docs.zip   --speed=6 --logstatus=/User/harryzhu/result
`

--speed=0/1/6/9 : 0~fastest without compression, 9~slowest with most compression. default is 1
--logstatus=/path/of/logfile.txt : log the global result(json format) to a file for monitoring 

decompress:

`
./zstdzip unzip --input=/User/harryzhu/test.zip  --output=/User/harryzhu/test
`
or you can use https://github.com/mcmilk/7-Zip-zstd to unzip

## Benchmark:
