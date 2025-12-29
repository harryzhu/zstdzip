# zstdzip
compress and decompress via ZSTD

## Performance:
Faster than `7-Zip`.

`zstdzip`: `compress` 563GB (246,516 mp4 files): 10 minutes 6 seconds; `decompress`: 8 minutes 56 seconds;

`7-zip`: `compress` same data above: 63 minutes 24 seconds;



## Usage:
### compress:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip

# default: compress paralell
```

or:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --level=1  --threads=16 --password=1234
```

or:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --ignore-dot-file --ignore-empty-dir --regext=".(mp4|txt|png)" --min-age=20230101081520 --max-age=20230220153045 --min-size-mb=4 --max-size-mb=16

# 选择文件夹 /User/harryzhu/docs 中最后修改时间晚于2023年01月01日08:15:20, 且早于2023年02月20日15:30:45,
# 且文件大小大于4MB小于16MB，
# 且文件后缀名为 .mp4 或 .txt 或 .png 的文件
# 忽略空文件夹
# 忽略隐藏文件（点 . 开头的文件名）
# 保存为 /User/harryzhu/docs.zst.zip 
# 采用并行压缩，会生成16个文件，每个文件都是完整的压缩文件，可以单独解压缩
```

or:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --serial

# --serial: 传统压缩模式, 一个一个压缩，最后保存在一个单一的压缩包内.
```

`--speed=0/1/2/3` : 

        0: fastest without compression, 

        1: default compression,

        2: better compression,
        
        3: best compression but slowest speed. 


### decompress:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/t2
```

or:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/test --serial
```

`--serial` : 一次只解压一个压缩档，如果不指定该参数，默认会同时解压缩8个压缩档，在SSD上面解压缩性能极速提升.在机械硬盘上，应该显示指定`--serial`，该参数主要针对机械硬盘优化。



### hash sum:

```Bash
./zstdzip hash --source=/User/harryzhu/test.zip  --sum=sha256
```

`--sum` : sum algorithm: md5, sha1, sha256, blake3, xxhash; default is `xxhash`


