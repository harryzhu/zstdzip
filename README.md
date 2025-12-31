# zstdzip
compress and decompress via ZSTD with zip format, can keep file's permission and last-modified timestamp.

采用ZSTD算法压缩文件，保存在zip格式包中，`压缩速度更快`、`压缩率更高`、还能`保留源文件的权限`和`最后修改时间`。

## Performance:
### Much Faster than `7-Zip`

数据量： `200GB` (`86,929` mp4 files)

`zstdzip`: 

* compress（压缩） : `2 minutes 46 seconds`; 
* decompress（解压缩）: `1 minutes 48 seconds`(宏碁GM7（PCIe 4.0）) 或者 `3 minutes 34 seconds`（SanDisk（PCIe 3.0））; 
* 在并行模式下，压缩/解压缩 的速度完全依赖于数据所在SSD的写入速度;
* 在传统模式下（添加 `--serial`参数，开启串行模式）, `11 minutes 8 seconds`（PCIe 4.0）;

`7-Zip`: 

* compress same data above: `22 minutes 34 seconds`，在 PCIe 4.0 和 3.0上速度几乎一样;
* 采用“不压缩、仅存储”选项: `4 minutes 43 seconds`(PCIe 4.0)


## Usage:
### Compress（压缩）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip

# default: 并行压缩，会自动生成 8 个压缩档
```

or（加密、压缩级别）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --level=1  --threads=16 --password=1234

# 压缩文件夹 /User/harryzhu/docs 里面的所有文件，保存为 /User/harryzhu/docs.zst.zip
# --password=1234： 文件使用指定密码加密
# 对应的，解压（unzip）时，也需要提供该密码才能解压成功
# --level=0 ｜ 1 ｜ 2 ｜ 3` : 
#   0: fastest without compression, 
#   1: default compression,
#   2: better compression,
#   3: best compression but slowest speed. 

```

or（按大小、时间、后缀名过滤文件）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --ignore-dot-file --ext=".(mp4|txt|png|jpg)" --min-age=2023-01-01,08:15:20 --max-age=2023-02-20,15:30:45 --min-size-mb=4 --max-size-mb=16

# 选择文件夹 /User/harryzhu/docs 中最后修改时间晚于2023年01月01日08:15:20, 且早于2023年02月20日15:30:45,
# 且文件大小大于4MB小于16MB，
# 且文件后缀名为 .mp4 或 .txt 或 .png 或 .jpg 的文件（不区分大小写）
# 忽略文件名以点 . 开头的文件
# 保存为 /User/harryzhu/docs.zst.zip 
# 默认采用并行压缩，会生成8个文件，每个文件都是完整的压缩文件，可以单独解压缩
#
```

or（机械硬盘中：将默认并行压缩改为传统的串行压缩，保存在单一的压缩包中）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --serial

# --serial: 传统压缩模式, 一个一个压缩，最后保存在 1个 单一的压缩包内， 不会自动生成 8个 压缩档
```

or（ dry-run 模式：只查看即将 压缩/解压缩 的文件列表，不会实际写入文件）:

```Bash
./zstdzip zip --source=/User/harryzhu/docs  --target=/User/harryzhu/docs.zst.zip  --dry-run

./zstdzip unzip --source=/User/harryzhu/docs.zst.zip  --target=/User/harryzhu/docs --dry-run

# --dry-run: 不写入文件，仅显示即将压缩或解压缩的文件列表，用于查看过滤的文件是否符合预期
#
```

### Decompress（解压缩）:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/t2

# 默认并行解压缩，会自动解压同文件夹下的另外 7 个压缩档 test.zip.1, test.zip.2, test.zip.3 ... test.zip.7
```

or（无需全部解压缩，可以挑选文件解压缩）:
```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/t2 --min-size-mb=4 --min-age=2023-02-15,14:30:12 --ext=".mp4" --ignore-empty-dir

# 默认并行解压缩，会自动解压同文件夹下的另外 7 个压缩档 test.zip.1, test.zip.2, test.zip.3 ... test.zip.7
# 可以用参数指定仅解压符合条件的文件，
# 上面表示： 仅解压文件大小超过4MB，文件最后修改时间晚于 2023-02-15 14:30:12 的后缀名为 .mp4 的文件， 
# --ignore-empty-dir 默认会忽略空文件夹，避免从大量文件中解压小部分文件会生成大量空文件夹，如果需要这些空文件夹，设置 --ignore-empty-dir=false 即可
```

or（机械硬盘中：串行解压缩）:

```Bash
./zstdzip unzip --source=/User/harryzhu/test.zip  --target=/User/harryzhu/test --serial

# --serial : 一次只解压一个压缩档，如果不指定该参数，默认会同时解压缩8个压缩档，在SSD上面解压缩性能极速提升.
# 但在机械硬盘上，应该显式指定`--serial`，io性能会更好。
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


