# 定义程序名称和源文件
PROGRAM = snunzip
SOURCE = main.go
BUILD_DIR = _dist

# 默认目标：编译所有平台
all: linux windows

# 创建输出目录
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)
	@echo "Build directory created: $(BUILD_DIR)"

# ----------------------------------------------------
# 目标 1: Linux (AMD64)
# GOOS=linux GOARCH=amd64
# ----------------------------------------------------
linux: $(BUILD_DIR)
	@echo "=================================================="
	@echo "Building Linux/amd64 executable..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(PROGRAM) $(SOURCE)
	@echo "Linux build SUCCESS: $(BUILD_DIR)/$(PROGRAM)"

# ----------------------------------------------------
# 目标 2: Windows (AMD64)
# GOOS=windows GOARCH=amd64
# ----------------------------------------------------
windows: $(BUILD_DIR)
	@echo "=================================================="
	@echo "Building Windows/amd64 executable..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(PROGRAM).exe $(SOURCE)
	@echo "Windows build SUCCESS: $(BUILD_DIR)/$(PROGRAM).exe"

# ----------------------------------------------------
# 目标 3: 清理
# ----------------------------------------------------
clean:
	@echo "=================================================="
	@echo "Cleaning up binaries..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleanup complete."

.PHONY: all linux windows clean