/**
 * NextKey C++ 客户端示例
 * 
 * 演示现代 C++ 风格的使用方式
 */

#include "NextKeyClient.hpp"
#include <iostream>
#include <thread>
#include <chrono>

using namespace std::chrono_literals;

int main() {
    // 配置信息（实际使用时请替换为真实值）
    const std::string server_url = "http://localhost:8080";
    const std::string project_uuid = "fe402b23-a193-47eb-9d7f-9c0a168e3cb3";
    const std::string aes_key = "78e54210cc4bdf4e6955a5e916f7000631d583e8dccc7ffb93525f53fdcbf061";
    const std::string card_key = "spFtLiotz8bTpYrr";
    const std::string hwid = "test-device-cpp-001";
    
    std::cout << "=== NextKey C++ 客户端示例 ===\n\n";
    
    try {
        // 1. 创建客户端 (RAII自动管理资源)
        std::cout << "[步骤 1] 创建 NextKey 客户端...\n";
        auto client = std::make_unique<nextkey::NextKeyClient>(server_url, project_uuid, aes_key);
        std::cout << "✓ 客户端创建成功\n\n";
        
        // 2. 登录
        std::cout << "[步骤 2] 登录中...\n";
        auto login_result = client->login(card_key, hwid);
        std::cout << "✓ 登录成功\n";
        std::cout << "  令牌: " << login_result.token << "\n";
        std::cout << "  过期时间: " << login_result.expire_at << "\n";
        std::cout << "  卡密信息:\n";
        std::cout << "    ID: " << login_result.card.id << "\n";
        std::cout << "    卡密: " << login_result.card.card_key << "\n";
        std::cout << "    已激活: " << (login_result.card.activated ? "是" : "否") << "\n";
        std::cout << "    时长: " << login_result.card.duration << " 秒\n";
        std::cout << "    专属信息: " << login_result.card.custom_data << "\n\n";
        
        // 3. 手动心跳测试
        std::cout << "[步骤 3] 测试手动心跳...\n";
        client->heartbeat();
        std::cout << "✓ 心跳正常\n\n";
        
        // 4. 启动自动心跳（带错误回调）
        std::cout << "[步骤 4] 启动自动心跳（5秒间隔）...\n";
        client->startAutoHeartbeat(5s, [](const nextkey::NextKeyException& e) {
            std::cerr << "[心跳错误] " << e.what() 
                     << " (错误码: " << e.code() << ")\n";
        });
        std::cout << "✓ 自动心跳已启动\n\n";
        
        // 5. 获取云变量
        std::cout << "[步骤 5] 获取云变量 'notice'...\n";
        try {
            auto value = client->getCloudVar("notice");
            std::cout << "✓ 云变量值: " << value << "\n\n";
        } catch (const nextkey::NextKeyException& e) {
            std::cerr << "✗ " << e.what() << "\n\n";
        }
        
        // 6. 更新专属信息
        std::cout << "[步骤 6] 更新专属数据...\n";
        std::string custom_data = "这是一个测试";
        client->updateCustomData(custom_data);
        std::cout << "✓ 专属数据已更新: " << custom_data << "\n\n";
        
        // 7. 获取项目信息
        std::cout << "[步骤 7] 获取项目信息...\n";
        auto proj_info = client->getProjectInfo();
        std::cout << "✓ 项目信息:\n";
        std::cout << "  UUID: " << proj_info.uuid << "\n";
        std::cout << "  名称: " << proj_info.name << "\n";
        std::cout << "  版本: " << proj_info.version << "\n";
        std::cout << "  更新地址: " << proj_info.update_url << "\n\n";
        
        // 8. 解绑HWID示例
        std::cout << "[步骤 8] 测试解绑HWID功能...\n";
        std::cout << "提示：此操作需要项目启用解绑功能\n";
        try {
            client->unbindHWID(card_key, hwid);
            std::cout << "✓ HWID解绑成功\n";
            std::cout << "  注意：解绑后需要重新登录才能在此设备使用\n\n";
        } catch (const nextkey::NextKeyException& e) {
            std::cerr << "✗ 解绑失败: " << e.what() << "\n";
            std::cerr << "  可能原因：项目未启用解绑、冷却期内、或HWID未绑定\n\n";
        }
        
        // 9. 运行一段时间观察心跳
        std::cout << "[步骤 9] 运行 10 秒以观察心跳...\n";
        std::cout << "（自动心跳将在后台运行）\n\n";
        
        for (int i = 10; i > 0; --i) {
            std::cout << "\r剩余时间: " << i << " 秒..." << std::flush;
            std::this_thread::sleep_for(1s);
        }
        std::cout << "\n\n";
        
        // 10. 清理资源（RAII自动完成，心跳会立即停止不阻塞）
        std::cout << "[步骤 10] 清理资源...\n";
        client->stopAutoHeartbeat();
        std::cout << "✓ 心跳已停止（detach模式，立即返回）\n";
        std::cout << "✓ 资源将自动清理（RAII）\n\n";
        
        std::cout << "=== 示例成功完成 ===\n";
        
    } catch (const nextkey::NextKeyException& e) {
        std::cerr << "\n[严重错误] " << e.what() 
                  << " (错误码: " << e.code() << ")\n";
        return 1;
    } catch (const std::exception& e) {
        std::cerr << "\n[异常] " << e.what() << "\n";
        return 1;
    }
    
    return 0;
}
