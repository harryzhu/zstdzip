# zstdzip
compress and decompress via ZSTD with zip format, can keep file's permission and last-modified timestamp.

## Performance:
Faster than `7-Zip`

`zstdzip`: compress 563GB (246,516 mp4 files): `10 minutes 6 seconds`; decompress: `8 minutes 56 seconds`;

`7-Zip`: compress same data above: `63 minutes 24 seconds`;



## Usage:
### compress:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip

# default: 并行压缩，会自动生成 8 个压缩档
```

or:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --level=1  --threads=16 --password=1234

# --password=1234： 文件使用指定密码加密
# 对应的，解压（unzip）时，也需要提供该密码才能解压成功
# --level=0 ｜ 1 ｜ 2 ｜ 3` : 
#   0: fastest without compression, 
#   1: default compression,
#   2: better compression,
#   3: best compression but slowest speed. 

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
# 采用并行压缩，会生成8个文件，每个文件都是完整的压缩文件，可以单独解压缩
```

or:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --serial

# --serial: 传统压缩模式, 一个一个压缩，最后保存在一个单一的压缩包内.
```



### decompress:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/t2

# 默认并行解压缩，会自动解压同文件夹下的另外 7 个压缩档 test.zip.1, test.zip.2, test.zip.3 ... test.zip.7
```

or:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/test --serial

# --serial : 一次只解压一个压缩档，如果不指定该参数，默认会同时解压缩8个压缩档，在SSD上面解压缩性能极速提升.
# 但在机械硬盘上，应该显示指定`--serial`，io性能会更好。
```





### hash sum:

```Bash
./zstdzip hash --source=/User/harryzhu/test.zip  --sum=sha256

# 显示文件的哈希：
# 使用 --sum= 指定哈希算法，支持 md5, sha1, sha256, blake3, xxhash
```

`--sum` : sum algorithm: md5, sha1, sha256, blake3, xxhash; default is `xxhash`


