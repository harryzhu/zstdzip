# zstdzip
compress and decompress via ZSTD with zip format, can keep file's permission and last-modified timestamp.

采用ZSTD算法压缩文件，保存在zip格式包中，`压缩速度更快`、`压缩率更高`、还能`保留源文件的权限`和`最后修改时间`。

## Performance:
### Much Faster than `7-Zip`

#### 数据量：`380,457`个 `.png` 文件，共`319GB`


压缩方式                      | 压缩时间(越短越好)   | 解压缩时间(越短越好)
-----------------------------|-------------------|-------------------
zstdzip(默认：串行)            |   27m 26s         | 5m 28s  
zstdzip(并行：--serial=false) | **6m 31s**        | **4m 4s** 
7zip ( .zip)                 |   12m 16s         |  

`zstdzip`: 

* compress（压缩） : 并行需要 `6 minutes 31 seconds`，串行需要 `27 minutes 26 seconds`; 
* decompress（解压缩）: 并行 `4 minutes 4 seconds`(宏碁GM7（PCIe 4.0），SanDisk（PCIe 3.0）) 或者 串行 `5 minutes 28 seconds`; 
* 在并行模式下，压缩/解压缩 的速度完全依赖于数据所在SSD的写入速度;


`7-Zip`: 

* compress same data above: `12 minutes 16 seconds`;


#### 数据量：`56,572`个 .txt 文件，共`21GB`

压缩方式 | 时间(越短越好) | 压缩后大小(越小越好)
--------|--------------|--------
zstdzip(并行)       |  **7s**    | 2.41GB
zstdzip(--level=2) |  13s       | 2.20GB
zstdzip(--level=3) |  34s       | 2.12GB
zstdzip(--serial)  |  37s       | 2.41GB
7zip ( .7z)        | 249s       | **1.46GB**
7zip ( .zip)       |  68s       | 3.56GB


## Usage:
### Compress（压缩）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zstd.zip 

# default: 串行压缩，会自动生成 1 个压缩档
# 添加 --serial=false： 并行压缩，会自动生成 8 个压缩档， 适用于SSD硬盘，超多文件压缩
```

or（加密、压缩级别）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zstd.zip  --level=1  --password=123

# 压缩文件夹 /User/harryzhu/docs 里面的所有文件，保存为 /User/harryzhu/docs.zst.zip
# --password=123： 文件使用指定密码加密
# 对应的，解压（unzip）时，也需要提供该密码才能解压成功
# --level=0 ｜ 1 ｜ 2 ｜ 3` : 
#   0: fastest without compression, 
#   1: default compression,
#   2: better compression,
#   3: best compression but slowest speed. 

```

or（按大小、时间、后缀名过滤文件）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zstd.zip  --ignore-dot-file --ext=".(mp4|txt|png|jpg)" --min-age=2023-01-01,08:15:20 --max-age=2023-02-20,15:30:45 --min-size-mb=4 --max-size-mb=16

# 选择文件夹 /User/harryzhu/docs 中最后修改时间晚于2023年01月01日08:15:20, 且早于2023年02月20日15:30:45,
# 且文件大小大于4MB小于16MB，
# 且文件后缀名为 .mp4 或 .txt 或 .png 或 .jpg 的文件（不区分大小写）
# 忽略文件名以点 . 开头的文件
# 保存为 /User/harryzhu/docs.zst.zip 
# 默认采用串行压缩，会生成1个文件，
# 如果文件非常多（且在固态SSD硬盘上），可以添加 --serial=false 参数，并行压缩，会生成8个文件，每个文件都是完整的压缩文件，可以单独解压缩（如果有一个损坏，不影响其他压缩包）
```

or（机械硬盘中：将并行压缩改为传统的串行压缩，保存在单一的压缩包中）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zstd.zip  --serial

# --serial: 传统压缩模式, 一个一个压缩，最后保存在 1个 单一的压缩包内， 不会自动生成 8个 压缩档
```

### Decompress（解压缩）:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zstd.zip  --target=/User/harryzhu/t2

# 默认串行解压缩，会自动解压同文件夹下的另外 7 个压缩档 test.zstd.zip.1, test.zstd.zip.2, test.zstd.zip.3 ... test.zstd.zip.7
```

or（机械硬盘中：串行解压缩）:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/test --serial=false

# --serial=false: 同时解压缩8个压缩档，在SSD上面解压缩性能极速提升.
# 但在机械硬盘上，可以显式指定`--serial`，io性能会更好。
```





### Hash（哈希）:

```Bash
./zstdzip hash --source=/User/harryzhu/test.zip  --sum=sha256

# 显示文件的哈希：
# 使用 --sum= 指定哈希算法，支持 md5, sha1, sha256, blake3, xxhash
# 多用于文件校验
# 对于超大文件哈希，推荐使用 blake3 或 xxhash 算法
```

`--sum` : sum algorithm: md5, sha1, sha256, blake3, xxhash; default is `sha256`


```Bash
./zstdzip hash --source=/User/harryzhu/folder --target=/User/harryzhu/folder_hash  --sum=xxhash

# 如果 --source= 是文件夹：
# 使用 --sum= 指定哈希算法，支持 md5, sha1, sha256, blake3, xxhash
# 递归对文件夹中的所有文件进行hash
# 并将结果保存在 --target= 中指定的文件
# 会保存两种格式， json 格式（格式：路径：hash值），以及 txt 格式（一行一个文件，格式： 哈希值：路径）
```

