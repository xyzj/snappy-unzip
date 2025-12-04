package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath" // 用于处理文件路径和命名
	"strings"

	"github.com/golang/snappy"
)

// progressReader 包装一个 io.Reader，统计已读取的字节并在 stderr 打印百分比进度
type progressReader struct {
	r       io.Reader
	total   int64
	read    int64
	lastPct int
	name    string
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	if n > 0 {
		pr.read += int64(n)
		if pr.total > 0 {
			pct := int(pr.read * 100 / pr.total)
			if pct != pr.lastPct {
				pr.lastPct = pct
				fmt.Fprintf(os.Stderr, "\rDecompressing %s: %3d%%", pr.name, pct)
			}
		}
	}
	return n, err
}

// decompressFile 负责单个文件的解压缩逻辑
func decompressFile(inputFilename string) error {
	// 1. 打开输入文件
	file, err := os.Open(inputFilename)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", inputFilename, err)
	}
	defer file.Close()

	// 2. 确定输出文件名
	// 假设解压后的文件应去除 .snappy 后缀
	outputFilename := inputFilename
	if strings.HasSuffix(inputFilename, ".snappy") {
		outputFilename = strings.TrimSuffix(inputFilename, ".snappy")
	} else {
		// 如果没有 .snappy 后缀，我们给它添加 .uncompressed
		outputFilename = inputFilename + ".uncompressed"
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(outputFilename); err == nil {
		return fmt.Errorf("output file %s already exists. Skipping to prevent overwrite", outputFilename)
	}

	// 3. 创建输出文件
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	// 4. 读取输入文件大小并包装带进度显示的 reader
	var reader io.Reader
	var pr *progressReader
	if fi, err2 := file.Stat(); err2 == nil {
		pr = &progressReader{r: file, total: fi.Size(), lastPct: -1, name: filepath.Base(inputFilename)}
		reader = snappy.NewReader(pr)
	} else {
		reader = snappy.NewReader(file)
	}

	// 5. 复制数据
	if _, err = io.Copy(outputFile, reader); err != nil {
		// 复制失败时，尝试清理已创建的文件
		os.Remove(outputFilename)
		return fmt.Errorf("error during snappy decompression of %s: %w", inputFilename, err)
	}

	// 保证进度显示到 100% 并换行（如果我们有总大小信息）
	if pr != nil && pr.total > 0 {
		fmt.Fprintf(os.Stderr, "\rDecompressing %s: %3d%%\n", pr.name, 100)
	}

	fmt.Printf("Successfully decompressed %s to %s\n", inputFilename, outputFilename)
	return nil
}

func main() {
	// os.Args[0] 是程序名本身
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file1.snappy> [file2.snappy] ...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Example: snunzip abc*")
		os.Exit(1)
	}

	// 从 os.Args[1] 开始遍历所有文件参数
	for _, filename := range os.Args[1:] {
		// 使用 filepath.Glob 进一步处理每个参数中的通配符 (可选，但更健壮)
		// 注意：在 Bash 中，通配符通常已由 Shell 展开，但如果参数是硬编码的通配符，此步有用

		matches, err := filepath.Glob(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Globbing error for %s: %v. Treating as literal file.\n", filename, err)
			matches = []string{filename}
		}

		if len(matches) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No files found matching %s.\n", filename)
			continue
		}

		for _, match := range matches {
			if err := decompressFile(match); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to process %s: %v\n", match, err)
			}
		}
	}
}
