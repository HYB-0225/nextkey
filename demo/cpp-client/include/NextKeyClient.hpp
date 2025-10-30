/**
 * NextKey C++ 客户端封装类
 * 
 * 提供 RAII 资源管理和异常处理
 */

#ifndef NEXTKEY_CLIENT_HPP
#define NEXTKEY_CLIENT_HPP

#include "nextkey_sdk.h"
#include <string>
#include <memory>
#include <stdexcept>
#include <chrono>
#include <thread>
#include <atomic>
#include <functional>

namespace nextkey {

// 异常类
class NextKeyException : public std::runtime_error {
public:
    NextKeyException(const std::string& message, int32_t code)
        : std::runtime_error(message), error_code_(code) {}
    
    int32_t code() const noexcept { return error_code_; }
    
private:
    int32_t error_code_;
};

// 卡密信息
struct CardInfo {
    uint64_t id;
    std::string card_key;
    bool activated;
    int64_t duration;
    std::string custom_data;
    
    CardInfo(uint64_t i, const std::string& ck, bool act, int64_t dur, const std::string& cd)
        : id(i), card_key(ck), activated(act), duration(dur), custom_data(cd) {}
};

// 登录结果
struct LoginResult {
    std::string token;
    std::string expire_at;
    CardInfo card;
    
    LoginResult(const std::string& t, const std::string& e, const CardInfo& c)
        : token(t), expire_at(e), card(c) {}
};

// 项目信息
struct ProjectInfo {
    std::string uuid;
    std::string name;
    std::string version;
    std::string update_url;
    
    ProjectInfo(const std::string& u, const std::string& n, const std::string& v, const std::string& url)
        : uuid(u), name(n), version(v), update_url(url) {}
};

// NextKey 客户端类 (RAII)
class NextKeyClient {
public:
    /**
     * 构造函数
     * @param server_url 服务器URL
     * @param project_uuid 项目UUID
     * @param aes_key AES密钥
     * @throws NextKeyException 创建失败时抛出异常
     */
    NextKeyClient(const std::string& server_url,
                  const std::string& project_uuid,
                  const std::string& aes_key);
    
    // 禁止拷贝
    NextKeyClient(const NextKeyClient&) = delete;
    NextKeyClient& operator=(const NextKeyClient&) = delete;
    
    // 移动构造和移动赋值
    NextKeyClient(NextKeyClient&& other) noexcept;
    NextKeyClient& operator=(NextKeyClient&& other) noexcept;
    
    // 析构函数
    ~NextKeyClient();
    
    /**
     * 登录
     * @param card_key 卡密
     * @param hwid 设备码（可选）
     * @param ip IP地址（可选）
     * @return 登录结果
     * @throws NextKeyException 登录失败时抛出异常
     */
    LoginResult login(const std::string& card_key,
                     const std::string& hwid = "",
                     const std::string& ip = "");
    
    /**
     * 心跳验证
     * @throws NextKeyException 心跳失败时抛出异常
     */
    void heartbeat();
    
    /**
     * 获取云变量
     * @param key 变量名
     * @return 变量值
     * @throws NextKeyException 获取失败时抛出异常
     */
    std::string getCloudVar(const std::string& key);
    
    /**
     * 更新专属信息
     * @param data 专属数据
     * @throws NextKeyException 更新失败时抛出异常
     */
    void updateCustomData(const std::string& data);
    
    /**
     * 获取项目信息
     * @return 项目信息
     * @throws NextKeyException 获取失败时抛出异常
     */
    ProjectInfo getProjectInfo();
    
    /**
     * 解绑HWID
     * @param card_key 卡密
     * @param hwid 要解绑的设备码
     * @throws NextKeyException 解绑失败时抛出异常
     */
    void unbindHWID(const std::string& card_key, const std::string& hwid);
    
    /**
     * 启动自动心跳
     * @param interval 心跳间隔（秒）
     * @param on_error 错误回调函数（可选）
     */
    void startAutoHeartbeat(
        std::chrono::seconds interval = std::chrono::seconds(30),
        std::function<void(const NextKeyException&)> on_error = nullptr
    );
    
    /**
     * 停止自动心跳
     */
    void stopAutoHeartbeat();
    
    /**
     * 检查是否正在运行自动心跳
     * @return true表示正在运行
     */
    bool isHeartbeatRunning() const noexcept;
    
private:
    void* client_handle_;  // C API 句柄 (避免类型冲突)
    std::thread heartbeat_thread_;
    std::atomic<bool> heartbeat_running_;
    std::function<void(const NextKeyException&)> heartbeat_error_callback_;
    
    void heartbeat_loop(std::chrono::seconds interval);
    void throw_if_error(const char* operation, int32_t error_code);
};

} // namespace nextkey

#endif // NEXTKEY_CLIENT_HPP

