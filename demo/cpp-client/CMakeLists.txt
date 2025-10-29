# 设置要求的 CMake 最低版本
cmake_minimum_required(VERSION 3.10)

# 定义项目名称
project(NextKeySDK)

# --- Android NDK 配置 ---
# 使用 NDK_HOME 环境变量设置 Android 工具链
if(DEFINED ENV{NDK_HOME})
    set(CMAKE_TOOLCHAIN_FILE "$ENV{NDK_HOME}/build/cmake/android.toolchain.cmake" CACHE STRING "Android Toolchain File")
else()
    message(FATAL_ERROR "NDK_HOME environment variable is not set. Please set it to your Android NDK path.")
endif()

# 设置目标 Android ABI 和 API 级别
set(ANDROID_ABI "arm64-v8a")
set(ANDROID_NATIVE_API_LEVEL 21)

set(CMAKE_CXX_FLAGS "\
-Wno-format-security \
-Wno-narrowing \
-fvisibility=hidden \
-fvisibility-inlines-hidden \
-w \
-O3 \
-s \
-g0 \
-fno-rtti \
-Wl,--strip-all \
-fno-asynchronous-unwind-tables \
-fno-stack-protector \
-fno-common \
-frtti \
${CMAKE_CXX_FLAGS}")

# --- 静态库产物 ---
# 从 src/ 目录中的源文件创建静态库
add_library(nextkey_client STATIC src/NextKeyClient.cpp)

# 为 nextkey_client 库指定公共头文件目录
# PUBLIC 关键字可确保链接到此库的任何目标都将继承这些包含目录
target_include_directories(nextkey_client PUBLIC 
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)

# --- 示例可执行文件产物 ---
# 从 examples/ 目录中的源文件创建可执行文件
add_executable(example examples/example.cpp)

# 将示例可执行文件链接到 nextkey_client 库和预编译的 libnextkey_sdk.a
target_link_libraries(example PRIVATE 
    nextkey_client
    ${CMAKE_CURRENT_SOURCE_DIR}/lib/libnextkey_sdk.a
)

# --- 输出信息 ---
message(STATUS "Project: ${PROJECT_NAME}")
message(STATUS "Configuring for Android ABI: ${ANDROID_ABI}")
message(STATUS "Using NDK Toolchain: ${CMAKE_TOOLCHAIN_FILE}")
