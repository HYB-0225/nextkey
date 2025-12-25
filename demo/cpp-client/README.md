# NextKey C++ Client

基于 Rust FFI 的现代 C++ 客户端封装，提供 RAII 资源管理和异常处理。

## 特性

- ✅ 现代 C++17 风格
- ✅ RAII 自动资源管理
- ✅ 异常处理机制
- ✅ 异步心跳支持
- ✅ 智能指针和移动语义
- ✅ 跨平台（Windows、Linux、Android）

## 编译

### 前置要求

1. CMake 3.15+
2. C++17 编译器
3. Rust 静态库（从 `../rust-sdk/` 编译）

### Linux

```bash
# 复制 Rust 静态库
cp ../rust-sdk/target/release/libnextkey_sdk.a lib/

# CMake 编译
mkdir build
cd build
cmake ..
cmake --build . --config Release
```

### Android (使用 NDK)

```bash
# 复制 Android ARM64 静态库
cp ../rust-sdk/target/aarch64-linux-android/release/libnextkey_sdk.a lib/

# 使用 CMake Android 工具链
mkdir build
cd build
cmake ..
cmake --build . --config Release
```

## 使用示例

### 基本流程

```cpp
#include "NextKeyClient.hpp"
#include <iostream>

using namespace nextkey;

int main() {
    try {
        // 1. 创建客户端 (RAII 自动管理资源)
        auto client = std::make_unique<NextKeyClient>(
            "http://localhost:8080",
            "your-project-uuid",
            "your-aes-key"
        );
        
        // 2. 登录
        auto result = client->login("your-card-key", "device-001");
        std::cout << "Token: " << result.token << "\n";
        std::cout << "Expire at: " << result.expire_at << "\n";
        
        // 3. 心跳
        client->heartbeat();
        std::cout << "Heartbeat OK\n";
        
        // 4. 获取云变量
        auto value = client->getCloudVar("notice");
        std::cout << "Cloud var: " << value << "\n";
        
        // 5. 更新专属信息
        client->updateCustomData(R"(这是一个测试)");
        std::cout << "Custom data updated\n";
        
        // 6. 获取项目信息
        auto info = client->getProjectInfo();
        std::cout << "Project: " << info.name << " v" << info.version << "\n";
        
        // 资源自动清理（RAII）
        
    } catch (const NextKeyException& e) {
        std::cerr << "Error: " << e.what() 
                  << " (code: " << e.code() << ")\n";
        return 1;
    }
    
    return 0;
}
```

### 自动心跳

```cpp
#include "NextKeyClient.hpp"
#include <chrono>

using namespace nextkey;
using namespace std::chrono_literals;

int main() {
    try {
        auto client = std::make_unique<NextKeyClient>(...);
        client->login("card-key", "hwid");
        
        // 启动自动心跳（30秒间隔）
        client->startAutoHeartbeat(30s, [](const NextKeyException& e) {
            std::cerr << "Heartbeat error: " << e.what() << "\n";
        });
        
        // 运行你的程序逻辑
        std::this_thread::sleep_for(std::chrono::minutes(5));
        
        // 停止心跳（或析构时自动停止）
        client->stopAutoHeartbeat();
        
    } catch (const NextKeyException& e) {
        std::cerr << "Error: " << e.what() << "\n";
        return 1;
    }
    
    return 0;
}
```

### 移动语义

```cpp
#include "NextKeyClient.hpp"

NextKeyClient createClient() {
    return NextKeyClient(
        "http://localhost:8080",
        "uuid",
        "key"
    );
}

int main() {
    // 移动构造
    auto client = createClient();
    
    // 移动赋值
    auto client2 = std::move(client);
    
    // client 已失效，client2 拥有资源
    client2.heartbeat();
    
    return 0;
}
```

## API 参考

### NextKeyClient 类

#### 构造函数
```cpp
NextKeyClient(
    const std::string& server_url,
    const std::string& project_uuid,
    const std::string& aes_key
);
```
创建客户端实例。失败抛出 `NextKeyException`。

#### login
```cpp
LoginResult login(
    const std::string& card_key,
    const std::string& hwid = "",
    const std::string& ip = ""
);
```
登录获取 Token。返回 `LoginResult` 结构体。

#### heartbeat
```cpp
void heartbeat();
```
发送心跳验证。失败抛出异常。

#### getCloudVar
```cpp
std::string getCloudVar(const std::string& key);
```
获取云变量值。失败抛出异常。

#### updateCustomData
```cpp
void updateCustomData(const std::string& data);
```
更新专属信息。失败抛出异常。

#### getProjectInfo
```cpp
ProjectInfo getProjectInfo();
```
获取项目信息。返回 `ProjectInfo` 结构体。

#### startAutoHeartbeat
```cpp
void startAutoHeartbeat(
    std::chrono::seconds interval = 30s,
    std::function<void(const NextKeyException&)> on_error = nullptr
);
```
启动自动心跳线程。

#### stopAutoHeartbeat
```cpp
void stopAutoHeartbeat();
```
停止自动心跳线程。

### 数据结构

#### LoginResult
```cpp
struct LoginResult {
    std::string token;
    std::string expire_at;
    CardInfo card;
};
```
**说明**: free 模式下 `card` 字段为默认空值（ID 为 0，字符串为空）。

#### ProjectInfo
```cpp
struct ProjectInfo {
    std::string uuid;
    std::string name;
    std::string version;
    std::string update_url;
};
```

### 异常类

#### NextKeyException
```cpp
class NextKeyException : public std::runtime_error {
public:
    int32_t code() const noexcept;
};
```

继承自 `std::runtime_error`，提供错误码访问。

## 项目结构

```
cpp-client/
├── include/
│   ├── NextKeyClient.hpp  # C++ 封装类
│   └── nextkey_sdk.h      # C FFI 头文件
├── src/
│   └── NextKeyClient.cpp  # 实现
├── lib/
│   └── libnextkey_sdk.a   # Rust 静态库
├── examples/
│   └── example.cpp        # 示例程序
└── CMakeLists.txt
```

## 最佳实践

### 1. 使用智能指针

```cpp
// 推荐
auto client = std::make_unique<NextKeyClient>(...);

// 不推荐
NextKeyClient* client = new NextKeyClient(...);
// ... 容易忘记 delete
```

### 2. 异常处理

```cpp
try {
    client->login("card-key");
    client->heartbeat();
} catch (const NextKeyException& e) {
    // 处理特定异常
    std::cerr << "NextKey error: " << e.what() 
              << " (code: " << e.code() << ")\n";
} catch (const std::exception& e) {
    // 处理其他异常
    std::cerr << "Error: " << e.what() << "\n";
}
```

### 3. 心跳错误回调

```cpp
client->startAutoHeartbeat(30s, [&](const NextKeyException& e) {
    if (e.code() == NEXTKEY_ERR_AUTH) {
        // 认证失败，重新登录
        try {
            client->login("card-key", "hwid");
        } catch (...) {
            std::cerr << "Re-login failed\n";
        }
    }
});
```

### 4. RAII 资源管理

```cpp
{
    NextKeyClient client(...);
    client.login(...);
    client.startAutoHeartbeat(30s);
    
    // 运行程序逻辑
    // ...
    
} // 离开作用域，自动停止心跳并释放资源
```

## 故障排查

### 编译错误

**C++17 不支持**:
```
error: 'make_unique' is not a member of 'std'
```
确保使用 C++17: `set(CMAKE_CXX_STANDARD 17)`

**链接错误**:
```
undefined reference to `nextkey_client_new'
```
确保正确链接 Rust 静态库和系统库。

### 运行错误

**异常捕获失败**:
确保使用 `const NextKeyException&` 引用捕获。

**心跳线程未停止**:
确保在析构前调用 `stopAutoHeartbeat()` 或让 RAII 自动处理。

## License

MIT

