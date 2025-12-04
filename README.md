# snunzip

一个小型命令行工具，用于解压 `.snappy` 文件（基于 Go + `github.com/golang/snappy`）。已在解压过程中添加压缩输入字节的百分比进度显示。

**Build**:

```bash
go build
# 或使用 Makefile（若可用）
make
```

**Usage**:

```bash
./snunzip file1.snappy [file2.snappy] ...
# 支持通配符（例如在没有被 shell 展开的场景下）
./snunzip *.snappy
```

运行时将向 `stderr` 打印类似的进度输出（逐步更新同一行）：

```
Decompressing file.snappy:  23%
```

完成时会确保显示 `100%` 并在下一行打印成功信息，例如：

```
Decompressing file.snappy: 100%
Successfully decompressed file.snappy to file
```

**实现说明**:
- **进度依据**: 进度为已读取的压缩输入字节数除以输入文件总大小的百分比（基于 `stat` 获取的文件大小）。这是对解压过程的合理估计，但不是基于解压后输出字节的精确百分比。
- **回退行为**: 如果无法获取输入文件大小，则不会显示百分比进度（仍然会进行解压）。
- **安全**: 如果目标输出文件已存在，程序会跳过该文件以防止覆盖。

**示例**:

```bash
# 解压单个文件并观察进度
./snunzip /path/to/log-2025-12-03.snappy

# 在当前目录解压所有 .snappy 文件
./snunzip *.snappy
```

如需更详细的显示（速率、已解压字节、预计剩余时间）或添加一个开关来开启/关闭进度显示，请提出需求，我可以继续扩展。

---
项目文件：`main.go`（包含 `progressReader` 实现）
