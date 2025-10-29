@echo off
echo Building NextKey Rust SDK...

REM 动态生成 .cargo/config.toml 使用 NDK_HOME 环境变量
if defined NDK_HOME (
    echo Configuring Android NDK from environment: %NDK_HOME%
    setlocal enabledelayedexpansion
    set "NDK_PATH=!NDK_HOME:\=/!"
    (
        echo [target.aarch64-linux-android]
        echo linker = "!NDK_PATH!/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android21-clang.cmd"
        echo ar = "!NDK_PATH!/toolchains/llvm/prebuilt/windows-x86_64/bin/llvm-ar.exe"
        echo.
        echo [build]
        echo target-dir = "target"
        echo.
        echo [env]
        echo CC_aarch64-linux-android = "!NDK_PATH!/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android21-clang.cmd"
        echo CXX_aarch64-linux-android = "!NDK_PATH!/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android21-clang++.cmd"
        echo AR_aarch64-linux-android = "!NDK_PATH!/toolchains/llvm/prebuilt/windows-x86_64/bin/llvm-ar.exe"
    ) > .cargo\config.toml
    endlocal
)

echo.
echo [1/4] Building CLI tool for Windows...
cargo build --release --bin nextkey-cli
if %errorlevel% neq 0 (
    echo Failed to build CLI tool
    exit /b 1
)

echo.
echo [2/4] Building static library for Windows...
cargo build --release --lib
if %errorlevel% neq 0 (
    echo Failed to build static library
    exit /b 1
)

echo.
echo [3/4] Building static library for Android ARM64...
echo (Optional - requires Android NDK configured correctly)
cargo build --release --lib --target aarch64-linux-android
if %errorlevel% neq 0 (
    echo Warning: Android library build failed (NDK may not be configured)
    echo Continuing anyway - Windows build is complete
)

echo.
echo [4/4] Syncing header file to C++ client...
if exist "include\nextkey_sdk.h" (
    copy /Y include\nextkey_sdk.h ..\cpp-client\include\nextkey_sdk.h > nul
    echo Header file synced to cpp-client
)

echo.
echo Build completed successfully!
echo.
echo Output files:
echo   CLI (Windows): target\release\nextkey-cli.exe
echo   Library (Windows): target\release\nextkey_sdk.lib
echo   Library (Android ARM64): target\aarch64-linux-android\release\libnextkey_sdk.a
echo   Header: include\nextkey_sdk.h (synced to cpp-client)

