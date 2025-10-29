#!/bin/bash

echo "正在构建 NextKey Rust SDK..."

# 动态生成 .cargo/config.toml 使用 NDK_HOME 环境变量
if [ -n "$NDK_HOME" ]; then
    echo "从环境变量配置 Android NDK: $NDK_HOME"
    
    # Termux 中的 NDK 路径通常是 Unix 风格，不需要转换
    NDK_PATH="$NDK_HOME"
    
    # 创建 .cargo 目录（如果不存在）
    mkdir -p .cargo
    
    # 生成 config.toml 文件 - 针对 Termux 环境调整工具链路径
    cat > .cargo/config.toml << EOF
[target.aarch64-linux-android]
linker = "$NDK_PATH/toolchains/llvm/prebuilt/linux-aarch64/bin/aarch64-linux-android21-clang"
ar = "$NDK_PATH/toolchains/llvm/prebuilt/linux-aarch64/bin/llvm-ar"

[build]
target-dir = "target"

[env]
CC_aarch64-linux-android = "$NDK_PATH/toolchains/llvm/prebuilt/linux-aarch64/bin/aarch64-linux-android21-clang"
CXX_aarch64-linux-android = "$NDK_PATH/toolchains/llvm/prebuilt/linux-aarch64/bin/aarch64-linux-android21-clang++"
AR_aarch64-linux-android = "$NDK_PATH/toolchains/llvm/prebuilt/linux-aarch64/bin/llvm-ar"
EOF
fi

echo
echo "[1/2] 正在为 Android ARM64 构建静态库..."
cargo build --release --lib --target aarch64-linux-android
if [ $? -ne 0 ]; then
    echo "错误: Android 库构建失败"
    exit 1
fi

echo
echo "[2/2] 正在为 Android ARM64 构建 CLI 工具..."
cargo build --release --bin nextkey-cli --target aarch64-linux-android
if [ $? -ne 0 ]; then
    echo "错误: CLI 工具构建失败"
    exit 1
fi

echo
echo "构建成功完成!"
echo
echo "输出文件:"
echo "  CLI (Android ARM64): target/aarch64-linux-android/release/nextkey-cli"
echo "  库 (Android ARM64): target/aarch64-linux-android/release/libnextkey_sdk.a"
echo "  头文件: include/nextkey_sdk.h"