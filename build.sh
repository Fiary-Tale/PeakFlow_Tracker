#!/bin/bash

# 设置输出目录
OUTPUT_DIR="./make"
mkdir -p $OUTPUT_DIR

# 定义目标平台和架构
TARGETS=(
    "darwin arm64"
    "darwin amd64"
    "linux amd64"
    "linux arm64"
)

# 设置版本号
VERSION="1.0.6"

# 为每个目标构建项目
for TARGET in "${TARGETS[@]}"; do
    OS=$(echo $TARGET | cut -d ' ' -f 1)
    ARCH=$(echo $TARGET | cut -d ' ' -f 2)

    OUTPUT_NAME="PeakFlow_Tracker-$OS-$ARCH-$VERSION"

    echo "Building for $OS $ARCH..."

    # 设置环境变量
    GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT_DIR/$OUTPUT_NAME"

    if [ $? -ne 0 ]; then
        echo "Failed to build for $OS $ARCH"
        exit 1
    fi

    echo "Successfully built $OUTPUT_NAME"
done

echo "All builds completed successfully."
