/**
 * NextKey C++ 客户端实现
 */

#include "NextKeyClient.hpp"
#include <cstring>
#include <iostream>

namespace nextkey {

// 使用 C API 的类型别名避免命名冲突
using CNextKeyClient = ::NextKeyClient;

// 辅助函数：抛出异常
void NextKeyClient::throw_if_error(const char* operation, int32_t error_code) {
    if (error_code != NEXTKEY_OK) {
        const char* last_error = nextkey_get_last_error();
        std::string message = std::string(operation) + " 失败";
        if (last_error) {
            message += ": " + std::string(last_error);
        }
        throw NextKeyException(message, error_code);
    }
}

// 构造函数
NextKeyClient::NextKeyClient(const std::string& server_url,
                             const std::string& project_uuid,
                             const std::string& aes_key)
    : client_handle_(nullptr), heartbeat_running_(false) {
    
    client_handle_ = reinterpret_cast<CNextKeyClient*>(::nextkey_client_new(
        server_url.c_str(),
        project_uuid.c_str(),
        aes_key.c_str()
    ));
    
    if (client_handle_ == nullptr) {
        const char* error = nextkey_get_last_error();
        throw NextKeyException(
            error ? std::string("创建客户端失败: ") + error 
                  : "创建客户端失败",
            NEXTKEY_ERR_UNKNOWN
        );
    }
}

// 移动构造函数
NextKeyClient::NextKeyClient(NextKeyClient&& other) noexcept
    : client_handle_(other.client_handle_),
      heartbeat_thread_(std::move(other.heartbeat_thread_)),
      heartbeat_running_(other.heartbeat_running_.load()),
      heartbeat_error_callback_(std::move(other.heartbeat_error_callback_)) {
    other.client_handle_ = nullptr;
    other.heartbeat_running_ = false;
}

// 移动赋值运算符
NextKeyClient& NextKeyClient::operator=(NextKeyClient&& other) noexcept {
    if (this != &other) {
        // 清理当前资源
        if (heartbeat_running_) {
            stopAutoHeartbeat();
        }
        if (client_handle_) {
            ::nextkey_client_free(reinterpret_cast<::NextKeyClient*>(client_handle_));
        }
        
        // 移动资源
        client_handle_ = other.client_handle_;
        heartbeat_thread_ = std::move(other.heartbeat_thread_);
        heartbeat_running_ = other.heartbeat_running_.load();
        heartbeat_error_callback_ = std::move(other.heartbeat_error_callback_);
        
        other.client_handle_ = nullptr;
        other.heartbeat_running_ = false;
    }
    return *this;
}

// 析构函数
NextKeyClient::~NextKeyClient() {
    if (heartbeat_running_) {
        heartbeat_running_ = false;
        if (heartbeat_thread_.joinable()) {
            heartbeat_thread_.detach();
        }
    }
    
    if (client_handle_) {
        ::nextkey_client_free(reinterpret_cast<::NextKeyClient*>(client_handle_));
        client_handle_ = nullptr;
    }
}

// 登录
LoginResult NextKeyClient::login(const std::string& card_key,
                                 const std::string& hwid,
                                 const std::string& ip) {
    char* token = nullptr;
    char* expire_at = nullptr;
    NextKeyCardInfo card_info = {0};
    
    int32_t result = ::nextkey_login(
        reinterpret_cast<::NextKeyClient*>(client_handle_),
        card_key.c_str(),
        hwid.empty() ? nullptr : hwid.c_str(),
        ip.empty() ? nullptr : ip.c_str(),
        &token,
        &expire_at,
        &card_info
    );
    
    if (result != NEXTKEY_OK) {
        throw_if_error("登录", result);
    }
    
    CardInfo card(
        card_info.id,
        card_info.card_key ? std::string(card_info.card_key) : "",
        card_info.activated == 1,
        card_info.duration,
        card_info.custom_data ? std::string(card_info.custom_data) : ""
    );
    
    LoginResult login_result(
        token ? std::string(token) : "",
        expire_at ? std::string(expire_at) : "",
        card
    );
    
    if (token) ::nextkey_free_string(token);
    if (expire_at) ::nextkey_free_string(expire_at);
    ::nextkey_free_card_info(&card_info);
    
    return login_result;
}

// 心跳
void NextKeyClient::heartbeat() {
    int32_t result = ::nextkey_heartbeat(reinterpret_cast<::NextKeyClient*>(client_handle_));
    throw_if_error("心跳", result);
}

// 获取云变量
std::string NextKeyClient::getCloudVar(const std::string& key) {
    char* value = nullptr;
    
    int32_t result = ::nextkey_get_cloud_var(
        reinterpret_cast<::NextKeyClient*>(client_handle_),
        key.c_str(),
        &value
    );
    
    if (result != NEXTKEY_OK) {
        throw_if_error("获取云变量", result);
    }
    
    std::string cloud_value = value ? std::string(value) : "";
    if (value) ::nextkey_free_string(value);
    
    return cloud_value;
}

// 更新专属信息
void NextKeyClient::updateCustomData(const std::string& data) {
    int32_t result = ::nextkey_update_custom_data(
        reinterpret_cast<::NextKeyClient*>(client_handle_),
        data.c_str()
    );
    
    throw_if_error("更新专属信息", result);
}

// 获取项目信息
ProjectInfo NextKeyClient::getProjectInfo() {
    char* uuid = nullptr;
    char* name = nullptr;
    char* version = nullptr;
    char* update_url = nullptr;
    
    int32_t result = ::nextkey_get_project_info(
        reinterpret_cast<::NextKeyClient*>(client_handle_),
        &uuid,
        &name,
        &version,
        &update_url
    );
    
    if (result != NEXTKEY_OK) {
        throw_if_error("获取项目信息", result);
    }
    
    ProjectInfo info(
        uuid ? std::string(uuid) : "",
        name ? std::string(name) : "",
        version ? std::string(version) : "",
        update_url ? std::string(update_url) : ""
    );
    
    if (uuid) ::nextkey_free_string(uuid);
    if (name) ::nextkey_free_string(name);
    if (version) ::nextkey_free_string(version);
    if (update_url) ::nextkey_free_string(update_url);
    
    return info;
}

// 心跳循环
void NextKeyClient::heartbeat_loop(std::chrono::seconds interval) {
    while (heartbeat_running_) {
        std::this_thread::sleep_for(interval);
        
        if (!heartbeat_running_) {
            break;
        }
        
        try {
            heartbeat();
            std::cout << "[心跳] 成功\n";
        } catch (const NextKeyException& e) {
            std::cerr << "[心跳] 失败: " << e.what() << " (错误码: " << e.code() << ")\n";
            
            if (heartbeat_error_callback_) {
                heartbeat_error_callback_(e);
            }
            
            // 认证失败时停止心跳
            if (e.code() == NEXTKEY_ERR_AUTH) {
                std::cerr << "[心跳] 认证失败，停止心跳\n";
                heartbeat_running_ = false;
                break;
            }
        }
    }
}

// 启动自动心跳
void NextKeyClient::startAutoHeartbeat(
    std::chrono::seconds interval,
    std::function<void(const NextKeyException&)> on_error
) {
    if (heartbeat_running_) {
        throw NextKeyException("心跳已在运行", NEXTKEY_ERR_UNKNOWN);
    }
    
    heartbeat_error_callback_ = on_error;
    heartbeat_running_ = true;
    
    heartbeat_thread_ = std::thread([this, interval]() {
        this->heartbeat_loop(interval);
    });
    
    std::cout << "自动心跳已启动 (间隔: " << interval.count() << " 秒)\n";
}

// 停止自动心跳
void NextKeyClient::stopAutoHeartbeat() {
    if (!heartbeat_running_) {
        return;
    }
    
    heartbeat_running_ = false;
    
    if (heartbeat_thread_.joinable()) {
        heartbeat_thread_.detach();
    }
    
    std::cout << "自动心跳已停止\n";
}

// 检查心跳是否运行
bool NextKeyClient::isHeartbeatRunning() const noexcept {
    return heartbeat_running_;
}

} // namespace nextkey
