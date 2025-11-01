use crate::client::NextKeyClient;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::ptr;

/// 成功返回码
/// 
/// 当API调用成功完成时返回此值
pub const NEXTKEY_OK: i32 = 0;

/// 参数无效错误码
/// 
/// 当传入的参数为空指针、格式错误或不满足要求时返回此值
pub const NEXTKEY_ERR_INVALID_PARAM: i32 = -1;

/// 网络错误码
/// 
/// 当网络请求失败、超时或服务器不可达时返回此值
pub const NEXTKEY_ERR_NETWORK: i32 = -2;

/// 认证失败错误码
/// 
/// 当卡密无效、已过期、Token失效或权限不足时返回此值
pub const NEXTKEY_ERR_AUTH: i32 = -401;

/// 解密失败错误码
/// 
/// 当响应数据解密失败或加密方案不匹配时返回此值
pub const NEXTKEY_ERR_DECRYPT: i32 = -3;

/// 未知错误码
/// 
/// 当遇到未预期的错误或内部异常时返回此值
pub const NEXTKEY_ERR_UNKNOWN: i32 = -999;

// 全局错误消息
static mut LAST_ERROR: Option<String> = None;

fn set_last_error(err: String) {
    unsafe {
        LAST_ERROR = Some(err);
    }
}

fn clear_last_error() {
    unsafe {
        LAST_ERROR = None;
    }
}

// 辅助函数：从C字符串转换
unsafe fn c_str_to_string(c_str: *const c_char) -> Option<String> {
    if c_str.is_null() {
        return None;
    }
    CStr::from_ptr(c_str).to_str().ok().map(|s| s.to_string())
}

// 辅助函数：转换为C字符串（调用者负责释放）
fn string_to_c_char(s: String) -> *mut c_char {
    match CString::new(s) {
        Ok(c_str) => c_str.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

/// 创建NextKey客户端（使用默认AES-256-GCM加密）
/// 
/// # 功能
/// 
/// 创建一个新的NextKey客户端实例，使用AES-256-GCM加密方案进行通信加密。
/// 这是最安全的加密方案，推荐用于生产环境。
/// 
/// # 参数
/// 
/// * `server_url` - 服务器地址，例如 "http://localhost:8080"，不能为NULL
/// * `project_uuid` - 项目UUID，从管理后台获取，不能为NULL
/// * `aes_key` - AES密钥（64位十六进制字符串，对应32字节），从项目配置获取，不能为NULL
/// 
/// # 返回值
/// 
/// * 成功: 返回客户端实例指针，需要调用 `nextkey_client_free` 释放
/// * 失败: 返回NULL，可通过 `nextkey_get_last_error` 获取错误信息
/// 
/// # 示例
/// 
/// ```c
/// const char* url = "http://localhost:8080";
/// const char* uuid = "550e8400-e29b-41d4-a716-446655440000";
/// const char* key = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037";
/// 
/// NextKeyClient* client = nextkey_client_new(url, uuid, key);
/// if (client == NULL) {
///     printf("创建客户端失败: %s\n", nextkey_get_last_error());
///     return -1;
/// }
/// 
/// // 使用客户端...
/// 
/// nextkey_client_free(client);
/// ```
/// 
/// # 线程安全性
/// 
/// 此函数是线程安全的，但返回的客户端实例不是线程安全的，
/// 不应在多个线程间共享同一个客户端实例。
/// 
/// # 另见
/// 
/// * `nextkey_client_new_with_scheme` - 指定加密方案创建客户端
/// * `nextkey_client_free` - 释放客户端资源
#[no_mangle]
pub extern "C" fn nextkey_client_new(
    server_url: *const c_char,
    project_uuid: *const c_char,
    aes_key: *const c_char,
) -> *mut NextKeyClient {
    // 使用默认加密方案（AES-256-GCM）调用带scheme参数的函数
    let scheme_cstr = std::ffi::CString::new("aes-256-gcm").unwrap();
    nextkey_client_new_with_scheme(server_url, project_uuid, aes_key, scheme_cstr.as_ptr())
}

/// 创建NextKey客户端（指定加密方案）
/// 
/// # 功能
/// 
/// 创建一个新的NextKey客户端实例，可自定义加密方案。
/// 支持多种加密方案以兼容不同的安全需求和性能要求。
/// 
/// # 参数
/// 
/// * `server_url` - 服务器地址，例如 "http://localhost:8080"，不能为NULL
/// * `project_uuid` - 项目UUID，从管理后台获取，不能为NULL
/// * `aes_key` - 加密密钥（格式根据scheme不同而不同），不能为NULL
/// * `scheme` - 加密方案名称（不区分大小写），支持以下值:
///   - `"aes-256-gcm"` - AES-256-GCM（推荐，高安全）
///   - `"chacha20-poly1305"` - ChaCha20-Poly1305（高安全，移动端友好）
///   - `"rc4"` - RC4（中等安全，快速）
///   - `"xor"` - XOR（低安全，仅用于测试）
///   - `"custom-base64"` - 自定义Base64编码（低安全，仅用于调试）
/// 
/// # 返回值
/// 
/// * 成功: 返回客户端实例指针，需要调用 `nextkey_client_free` 释放
/// * 失败: 返回NULL，可通过 `nextkey_get_last_error` 获取错误信息
/// 
/// # 示例
/// 
/// ```c
/// // 使用RC4加密方案
/// NextKeyClient* client = nextkey_client_new_with_scheme(
///     "http://localhost:8080",
///     "550e8400-e29b-41d4-a716-446655440000",
///     "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037",
///     "rc4"
/// );
/// 
/// if (client == NULL) {
///     printf("创建客户端失败: %s\n", nextkey_get_last_error());
///     return -1;
/// }
/// 
/// nextkey_client_free(client);
/// ```
/// 
/// # 注意事项
/// 
/// * 加密方案必须与服务器项目配置一致，否则通信会失败
/// * 生产环境推荐使用 `aes-256-gcm`
/// * 密钥格式必须与加密方案匹配（通常为64位十六进制字符串）
/// 
/// # 线程安全性
/// 
/// 此函数是线程安全的，但返回的客户端实例不是线程安全的。
#[no_mangle]
pub extern "C" fn nextkey_client_new_with_scheme(
    server_url: *const c_char,
    project_uuid: *const c_char,
    aes_key: *const c_char,
    scheme: *const c_char,
) -> *mut NextKeyClient {
    clear_last_error();

    let server_url = unsafe {
        match c_str_to_string(server_url) {
            Some(s) => s,
            None => {
                set_last_error("server_url参数无效".to_string());
                return ptr::null_mut();
            }
        }
    };

    let project_uuid = unsafe {
        match c_str_to_string(project_uuid) {
            Some(s) => s,
            None => {
                set_last_error("project_uuid参数无效".to_string());
                return ptr::null_mut();
            }
        }
    };

    let aes_key = unsafe {
        match c_str_to_string(aes_key) {
            Some(s) => s,
            None => {
                set_last_error("aes_key参数无效".to_string());
                return ptr::null_mut();
            }
        }
    };

    let scheme_str = unsafe {
        match c_str_to_string(scheme) {
            Some(s) => s,
            None => {
                set_last_error("scheme参数无效".to_string());
                return ptr::null_mut();
            }
        }
    };

    // 解析加密方案
    use crate::crypto::EncryptionScheme;
    let encryption_scheme = match scheme_str.to_lowercase().as_str() {
        "aes-256-gcm" => EncryptionScheme::AES256GCM,
        "chacha20-poly1305" => EncryptionScheme::ChaCha20Poly1305,
        "rc4" => EncryptionScheme::RC4,
        "xor" => EncryptionScheme::XOR,
        "custom-base64" => EncryptionScheme::CustomBase64,
        _ => {
            set_last_error(format!("不支持的加密方案: {}", scheme_str));
            return ptr::null_mut();
        }
    };

    match NextKeyClient::new_with_scheme(&server_url, &project_uuid, &aes_key, encryption_scheme) {
        Ok(client) => Box::into_raw(Box::new(client)),
        Err(e) => {
            set_last_error(format!("创建客户端失败: {}", e));
            ptr::null_mut()
        }
    }
}

/// 释放客户端资源
/// 
/// # 功能
/// 
/// 释放由 `nextkey_client_new` 或 `nextkey_client_new_with_scheme` 创建的客户端实例。
/// 此函数会清理所有相关资源，包括网络连接和内存。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针，可以为NULL（为NULL时函数不执行任何操作）
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new(url, uuid, key);
/// 
/// // 使用客户端...
/// 
/// nextkey_client_free(client);
/// client = NULL; // 推荐：避免悬垂指针
/// ```
/// 
/// # 注意事项
/// 
/// * 同一个客户端指针只能释放一次，重复释放会导致未定义行为
/// * 释放后不应再使用该指针
/// * 传入NULL是安全的，函数会忽略
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的，确保在释放客户端时没有其他线程正在使用它。
#[no_mangle]
pub extern "C" fn nextkey_client_free(client: *mut NextKeyClient) {
    if !client.is_null() {
        unsafe {
            let _ = Box::from_raw(client);
        }
    }
}

/// 卡密信息结构体
/// 
/// # 功能
/// 
/// 存储卡密的详细信息，由 `nextkey_login` 函数返回。
/// 此结构体使用C ABI布局，可在C/C++代码中直接使用。
/// 
/// # 字段
/// 
/// * `id` - 卡密ID，数据库中的唯一标识
/// * `card_key` - 卡密字符串，由SDK分配的C字符串指针
/// * `activated` - 激活状态（0=未激活，1=已激活）
/// * `duration` - 有效时长（秒），负数表示永久有效
/// * `custom_data` - 卡密专属信息，由SDK分配的C字符串指针，可为空字符串
/// 
/// # 内存管理
/// 
/// * `card_key` 和 `custom_data` 字段指向的内存由SDK分配
/// * 使用完毕后必须调用 `nextkey_free_card_info` 释放这些字符串
/// * 不要手动释放 `card_key` 和 `custom_data`，否则会导致双重释放
/// * 结构体本身可以在栈上分配，只需释放内部的字符串指针
/// 
/// # 示例
/// 
/// ```c
/// NextKeyCardInfo card_info;
/// char* token = NULL;
/// char* expire_at = NULL;
/// 
/// int ret = nextkey_login(client, card_key, hwid, ip, &token, &expire_at, &card_info);
/// if (ret == NEXTKEY_OK) {
///     printf("卡密ID: %llu\n", card_info.id);
///     printf("激活状态: %s\n", card_info.activated ? "已激活" : "未激活");
///     printf("有效时长: %lld秒\n", card_info.duration);
///     printf("专属信息: %s\n", card_info.custom_data);
///     
///     // 释放字符串字段
///     nextkey_free_card_info(&card_info);
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
/// }
/// ```
/// 
/// # 注意事项
/// 
/// * 结构体初始化后，在传入 `nextkey_login` 前不需要设置字段
/// * 调用 `nextkey_free_card_info` 后，不要再访问 `card_key` 和 `custom_data`
#[repr(C)]
pub struct NextKeyCardInfo {
    pub id: u64,
    pub card_key: *mut c_char,
    pub activated: i32,
    pub duration: i64,
    pub custom_data: *mut c_char,
}

/// 卡密登录验证
/// 
/// # 功能
/// 
/// 使用卡密、硬件ID和IP地址进行登录验证。
/// 登录成功后会返回访问令牌(Token)和卡密信息，Token会自动保存到客户端实例中。
/// 
/// # 参数
/// 
/// **输入参数:**
/// * `client` - 客户端实例指针，不能为NULL
/// * `card_key` - 卡密字符串，不能为NULL
/// * `hwid` - 硬件ID，用于设备绑定，可以为空字符串但不能为NULL
/// * `ip` - 客户端IP地址，可以为空字符串但不能为NULL
/// 
/// **输出参数:**
/// * `token_out` - 输出访问令牌的指针，不能为NULL，需要调用 `nextkey_free_string` 释放
/// * `expire_at_out` - 输出令牌过期时间的指针（ISO 8601格式），不能为NULL，需要调用 `nextkey_free_string` 释放
/// * `card_info_out` - 输出卡密信息的指针，不能为NULL，需要调用 `nextkey_free_card_info` 释放
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 登录成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（卡密无效、已过期或HWID不匹配）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * `NEXTKEY_ERR_DECRYPT` (-3) - 解密失败
/// * 其他错误码 - 服务器返回的业务错误码
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new(url, uuid, key);
/// 
/// char* token = NULL;
/// char* expire_at = NULL;
/// NextKeyCardInfo card_info;
/// 
/// int ret = nextkey_login(client, "ABCD-1234-EFGH", "hwid123", "192.168.1.1",
///                         &token, &expire_at, &card_info);
/// 
/// if (ret == NEXTKEY_OK) {
///     printf("登录成功\n");
///     printf("Token: %s\n", token);
///     printf("过期时间: %s\n", expire_at);
///     printf("卡密ID: %llu\n", card_info.id);
///     
///     // 释放内存
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
///     nextkey_free_card_info(&card_info);
/// } else {
///     printf("登录失败 (错误码: %d): %s\n", ret, nextkey_get_last_error());
/// }
/// 
/// nextkey_client_free(client);
/// ```
/// 
/// # 注意事项
/// 
/// * 登录成功后Token会自动保存在客户端实例中，后续API调用会自动使用
/// * 必须在调用其他需要认证的API前先登录
/// * HWID和IP可以为空字符串，但建议提供以增强安全性
/// * 卡密验证失败时，不要频繁重试，避免触发限流
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的，不要在多个线程中同时对同一个客户端实例调用。
#[no_mangle]
pub extern "C" fn nextkey_login(
    client: *mut NextKeyClient,
    card_key: *const c_char,
    hwid: *const c_char,
    ip: *const c_char,
    token_out: *mut *mut c_char,
    expire_at_out: *mut *mut c_char,
    card_info_out: *mut NextKeyCardInfo,
) -> i32 {
    clear_last_error();

    if client.is_null() || card_key.is_null() || token_out.is_null() 
        || expire_at_out.is_null() || card_info_out.is_null() {
        set_last_error("参数无效".to_string());
        return NEXTKEY_ERR_INVALID_PARAM;
    }

    let client = unsafe { &mut *client };

    let card_key = unsafe { c_str_to_string(card_key).unwrap_or_default() };
    let hwid = unsafe { c_str_to_string(hwid).unwrap_or_default() };
    let ip = unsafe { c_str_to_string(ip).unwrap_or_default() };

    match client.login(&card_key, &hwid, &ip) {
        Ok(response) => {
            if response.code != 0 {
                set_last_error(format!("登录失败: {}", response.message));
                return response.code;
            }

            if let Some(data) = response.data {
                unsafe {
                    *token_out = string_to_c_char(data.token);
                    *expire_at_out = string_to_c_char(data.expire_at);
                    
                    // 填充CardInfo
                    (*card_info_out).id = data.card.id;
                    (*card_info_out).card_key = string_to_c_char(data.card.card_key);
                    (*card_info_out).activated = if data.card.activated { 1 } else { 0 };
                    (*card_info_out).duration = data.card.duration;
                    (*card_info_out).custom_data = string_to_c_char(data.card.custom_data);
                }
                NEXTKEY_OK
            } else {
                set_last_error("响应数据为空".to_string());
                NEXTKEY_ERR_UNKNOWN
            }
        }
        Err(e) => {
            set_last_error(format!("登录异常: {}", e));
            NEXTKEY_ERR_NETWORK
        }
    }
}

/// 发送心跳保持会话活跃
/// 
/// # 功能
/// 
/// 向服务器发送心跳请求，验证Token有效性并保持会话活跃。
/// 服务器会检查卡密状态，如果卡密被禁用或过期，心跳会失败。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针，必须已登录，不能为NULL
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 心跳成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（未登录或Token失效）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他错误码 - 服务器返回的业务错误码
/// 
/// # 示例
/// 
/// ```c
/// // 登录后定期发送心跳
/// while (running) {
///     int ret = nextkey_heartbeat(client);
///     if (ret == NEXTKEY_OK) {
///         printf("心跳成功\n");
///     } else if (ret == NEXTKEY_ERR_AUTH) {
///         printf("会话失效，需要重新登录\n");
///         break;
///     } else {
///         printf("心跳失败: %s\n", nextkey_get_last_error());
///     }
///     
///     sleep(300); // 每5分钟发送一次
/// }
/// ```
/// 
/// # 使用建议
/// 
/// * 建议每5-10分钟发送一次心跳
/// * 在长时间运行的程序中使用心跳检测Token有效性
/// * 心跳失败时应考虑重新登录
/// * 不要过于频繁发送心跳，避免浪费服务器资源
/// 
/// # 前置条件
/// 
/// * 必须先调用 `nextkey_login` 登录成功
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的，不要在多个线程中同时对同一个客户端实例调用。
#[no_mangle]
pub extern "C" fn nextkey_heartbeat(client: *mut NextKeyClient) -> i32 {
    clear_last_error();

    if client.is_null() {
        set_last_error("client参数无效".to_string());
        return NEXTKEY_ERR_INVALID_PARAM;
    }

    let client = unsafe { &*client };

    match client.heartbeat() {
        Ok(response) => {
            if response.code != 0 {
                set_last_error(format!("心跳失败: {}", response.message));
            }
            response.code
        }
        Err(e) => {
            set_last_error(format!("心跳异常: {}", e));
            NEXTKEY_ERR_NETWORK
        }
    }
}

/// 获取云变量值
/// 
/// # 功能
/// 
/// 根据键名获取项目级别的云变量值。
/// 云变量是存储在服务器上的键值对，可用于配置管理、开关控制等场景。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针，必须已登录，不能为NULL
/// * `key` - 云变量键名，不能为NULL
/// * `value_out` - 输出云变量值的指针，不能为NULL，需要调用 `nextkey_free_string` 释放
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 获取成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（未登录或Token失效）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他错误码 - 服务器返回的业务错误码（如404表示键不存在）
/// 
/// # 示例
/// 
/// ```c
/// char* value = NULL;
/// int ret = nextkey_get_cloud_var(client, "app_version", &value);
/// 
/// if (ret == NEXTKEY_OK) {
///     printf("当前版本: %s\n", value);
///     nextkey_free_string(value);
/// } else if (ret == 404) {
///     printf("云变量不存在\n");
/// } else {
///     printf("获取失败: %s\n", nextkey_get_last_error());
/// }
/// ```
/// 
/// # 使用场景
/// 
/// * 远程配置管理 - 版本号、功能开关等
/// * 动态内容分发 - 公告、更新地址等
/// * 全局开关控制 - 维护模式、功能启用等
/// 
/// # 注意事项
/// 
/// * 云变量在管理后台配置，属于项目级别
/// * 键名不存在时会返回错误码404
/// * 建议对返回的值做本地缓存，减少网络请求
/// 
/// # 前置条件
/// 
/// * 必须先调用 `nextkey_login` 登录成功
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的。
#[no_mangle]
pub extern "C" fn nextkey_get_cloud_var(
    client: *mut NextKeyClient,
    key: *const c_char,
    value_out: *mut *mut c_char,
) -> i32 {
    clear_last_error();

    if client.is_null() || key.is_null() || value_out.is_null() {
        set_last_error("参数无效".to_string());
        return NEXTKEY_ERR_INVALID_PARAM;
    }

    let client = unsafe { &*client };
    let key = unsafe { c_str_to_string(key).unwrap_or_default() };

    match client.get_cloud_var(&key) {
        Ok(response) => {
            if response.code != 0 {
                set_last_error(format!("获取云变量失败: {}", response.message));
                return response.code;
            }

            if let Some(data) = response.data {
                unsafe {
                    *value_out = string_to_c_char(data.value);
                }
                NEXTKEY_OK
            } else {
                set_last_error("响应数据为空".to_string());
                NEXTKEY_ERR_UNKNOWN
            }
        }
        Err(e) => {
            set_last_error(format!("获取云变量异常: {}", e));
            NEXTKEY_ERR_NETWORK
        }
    }
}

/// 更新卡密专属信息
/// 
/// # 功能
/// 
/// 更新当前卡密的专属信息（custom_data字段）。
/// 专属信息是与卡密绑定的自定义数据，可用于存储用户备注、配置等。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针，必须已登录，不能为NULL
/// * `custom_data` - 要设置的专属信息字符串，不能为NULL（可以为空字符串表示清空）
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 更新成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（未登录或Token失效）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他错误码 - 服务器返回的业务错误码
/// 
/// # 示例
/// 
/// ```c
/// // 设置专属信息
/// int ret = nextkey_update_custom_data(client, "user_level:5|coins:1000");
/// if (ret == NEXTKEY_OK) {
///     printf("专属信息更新成功\n");
/// } else {
///     printf("更新失败: %s\n", nextkey_get_last_error());
/// }
/// 
/// // 清空专属信息
/// nextkey_update_custom_data(client, "");
/// ```
/// 
/// # 使用场景
/// 
/// * 存储用户等级、积分等游戏数据
/// * 保存用户偏好设置
/// * 记录备注信息
/// * 标记特殊状态
/// 
/// # 注意事项
/// 
/// * 专属信息长度有限制（通常不超过1024字符），超长会被截断
/// * 更新会覆盖原有内容，不是追加
/// * 专属信息可在管理后台查看
/// * 传入空字符串会清空专属信息
/// 
/// # 前置条件
/// 
/// * 必须先调用 `nextkey_login` 登录成功
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的。
#[no_mangle]
pub extern "C" fn nextkey_update_custom_data(
    client: *mut NextKeyClient,
    custom_data: *const c_char,
) -> i32 {
    clear_last_error();

    if client.is_null() || custom_data.is_null() {
        set_last_error("参数无效".to_string());
        return NEXTKEY_ERR_INVALID_PARAM;
    }

    let client = unsafe { &*client };
    let custom_data = unsafe { c_str_to_string(custom_data).unwrap_or_default() };

    match client.update_custom_data(&custom_data) {
        Ok(response) => {
            if response.code != 0 {
                set_last_error(format!("更新专属信息失败: {}", response.message));
            }
            response.code
        }
        Err(e) => {
            set_last_error(format!("更新专属信息异常: {}", e));
            NEXTKEY_ERR_NETWORK
        }
    }
}

/// 获取项目信息
/// 
/// # 功能
/// 
/// 获取当前项目的基本信息，包括UUID、名称、版本号和更新地址。
/// 可用于版本检测、更新提示等场景。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针，必须已登录，不能为NULL
/// * `uuid_out` - 输出项目UUID的指针，不能为NULL，需要调用 `nextkey_free_string` 释放
/// * `name_out` - 输出项目名称的指针，不能为NULL，需要调用 `nextkey_free_string` 释放
/// * `version_out` - 输出项目版本号的指针，不能为NULL，需要调用 `nextkey_free_string` 释放
/// * `update_url_out` - 输出更新地址的指针，不能为NULL，需要调用 `nextkey_free_string` 释放
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 获取成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（未登录或Token失效）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他错误码 - 服务器返回的业务错误码
/// 
/// # 示例
/// 
/// ```c
/// char* uuid = NULL;
/// char* name = NULL;
/// char* version = NULL;
/// char* update_url = NULL;
/// 
/// int ret = nextkey_get_project_info(client, &uuid, &name, &version, &update_url);
/// 
/// if (ret == NEXTKEY_OK) {
///     printf("项目名称: %s\n", name);
///     printf("当前版本: %s\n", version);
///     
///     // 版本检测逻辑
///     if (strcmp(version, LOCAL_VERSION) > 0) {
///         printf("发现新版本，下载地址: %s\n", update_url);
///     }
///     
///     // 释放内存
///     nextkey_free_string(uuid);
///     nextkey_free_string(name);
///     nextkey_free_string(version);
///     nextkey_free_string(update_url);
/// } else {
///     printf("获取项目信息失败: %s\n", nextkey_get_last_error());
/// }
/// ```
/// 
/// # 使用场景
/// 
/// * 启动时检测版本更新
/// * 显示项目名称和版本号
/// * 自动更新功能
/// * 版本兼容性检查
/// 
/// # 注意事项
/// 
/// * 所有输出字符串都需要手动释放
/// * update_url可能为空字符串（如果未配置更新地址）
/// * 版本号格式由项目配置决定，建议使用语义化版本
/// 
/// # 前置条件
/// 
/// * 必须先调用 `nextkey_login` 登录成功
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的。
#[no_mangle]
pub extern "C" fn nextkey_get_project_info(
    client: *mut NextKeyClient,
    uuid_out: *mut *mut c_char,
    name_out: *mut *mut c_char,
    version_out: *mut *mut c_char,
    update_url_out: *mut *mut c_char,
) -> i32 {
    clear_last_error();

    if client.is_null() || uuid_out.is_null() || name_out.is_null() 
        || version_out.is_null() || update_url_out.is_null() {
        set_last_error("参数无效".to_string());
        return NEXTKEY_ERR_INVALID_PARAM;
    }

    let client = unsafe { &*client };

    match client.get_project_info() {
        Ok(response) => {
            if response.code != 0 {
                set_last_error(format!("获取项目信息失败: {}", response.message));
                return response.code;
            }

            if let Some(data) = response.data {
                unsafe {
                    *uuid_out = string_to_c_char(data.uuid);
                    *name_out = string_to_c_char(data.name);
                    *version_out = string_to_c_char(data.version);
                    *update_url_out = string_to_c_char(data.update_url);
                }
                NEXTKEY_OK
            } else {
                set_last_error("响应数据为空".to_string());
                NEXTKEY_ERR_UNKNOWN
            }
        }
        Err(e) => {
            set_last_error(format!("获取项目信息异常: {}", e));
            NEXTKEY_ERR_NETWORK
        }
    }
}

/// 解绑硬件ID
/// 
/// # 功能
/// 
/// 解除卡密与指定硬件ID的绑定关系。
/// 解绑后，该卡密可以在新设备上重新激活使用。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针，必须已登录，不能为NULL
/// * `card_key` - 要解绑的卡密字符串，不能为NULL
/// * `hwid` - 要解绑的硬件ID，不能为NULL
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 解绑成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（未登录、Token失效或无权限）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他错误码 - 服务器返回的业务错误码（如卡密不存在、HWID不匹配等）
/// 
/// # 示例
/// 
/// ```c
/// int ret = nextkey_unbind_hwid(client, "ABCD-1234-EFGH", "old_hwid");
/// 
/// if (ret == NEXTKEY_OK) {
///     printf("HWID解绑成功，卡密可在新设备使用\n");
/// } else if (ret == NEXTKEY_ERR_AUTH) {
///     printf("无权限解绑，请检查Token是否有效\n");
/// } else {
///     printf("解绑失败: %s\n", nextkey_get_last_error());
/// }
/// ```
/// 
/// # 使用场景
/// 
/// * 用户更换设备
/// * 设备损坏需要迁移卡密
/// * 管理员协助用户解绑
/// * 解除错误绑定
/// 
/// # 注意事项
/// 
/// * 解绑操作可能需要特定权限，具体取决于服务器配置
/// * 解绑后卡密状态会重置，下次登录时会重新绑定新设备
/// * 频繁解绑可能触发风控限制
/// * 卡密和HWID必须完全匹配才能解绑
/// 
/// # 安全考虑
/// 
/// * 建议对解绑操作进行日志记录
/// * 可能需要额外的身份验证（如验证码）
/// * 限制解绑次数防止滥用
/// 
/// # 前置条件
/// 
/// * 必须先调用 `nextkey_login` 登录成功
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的。
#[no_mangle]
pub extern "C" fn nextkey_unbind_hwid(
    client: *mut NextKeyClient,
    card_key: *const c_char,
    hwid: *const c_char,
) -> i32 {
    clear_last_error();

    if client.is_null() || card_key.is_null() || hwid.is_null() {
        set_last_error("参数无效".to_string());
        return NEXTKEY_ERR_INVALID_PARAM;
    }

    let card_key = unsafe {
        match c_str_to_string(card_key) {
            Some(s) => s,
            None => {
                set_last_error("card_key参数无效".to_string());
                return NEXTKEY_ERR_INVALID_PARAM;
            }
        }
    };

    let hwid = unsafe {
        match c_str_to_string(hwid) {
            Some(s) => s,
            None => {
                set_last_error("hwid参数无效".to_string());
                return NEXTKEY_ERR_INVALID_PARAM;
            }
        }
    };

    let client = unsafe { &*client };

    match client.unbind_hwid(&card_key, &hwid) {
        Ok(response) => {
            if response.code != 0 {
                set_last_error(format!("解绑HWID失败: {}", response.message));
                return response.code;
            }
            NEXTKEY_OK
        }
        Err(e) => {
            set_last_error(format!("解绑HWID异常: {}", e));
            NEXTKEY_ERR_NETWORK
        }
    }
}

/// 获取最后的错误消息
/// 
/// # 功能
/// 
/// 获取最近一次FFI函数调用失败时的详细错误信息。
/// 每次调用FFI函数前会清除上次的错误消息，因此只能获取最近一次的错误。
/// 
/// # 返回值
/// 
/// * 有错误时: 返回指向错误消息的C字符串指针（UTF-8编码）
/// * 无错误时: 返回NULL
/// 
/// # 内存管理
/// 
/// * 返回的指针由SDK内部管理，**不需要**调用 `nextkey_free_string` 释放
/// * 指针的生命周期：在下次调用任何FFI函数前有效
/// * 调用者应立即复制错误消息内容，不要长期持有该指针
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new(url, uuid, invalid_key);
/// 
/// if (client == NULL) {
///     const char* error = nextkey_get_last_error();
///     if (error != NULL) {
///         // 立即使用或复制错误消息
///         printf("创建客户端失败: %s\n", error);
///         
///         // 如果需要长期保存，应复制字符串
///         char error_copy[256];
///         strncpy(error_copy, error, sizeof(error_copy) - 1);
///     }
/// }
/// 
/// // 下次调用FFI函数会清除错误消息
/// client = nextkey_client_new(url, uuid, valid_key);
/// ```
/// 
/// # 使用建议
/// 
/// * 在任何返回错误码或NULL的FFI调用后立即检查
/// * 用于调试和日志记录
/// * 不要依赖错误消息的具体格式，它们可能会改变
/// 
/// # 注意事项
/// 
/// * 错误消息为中文，适合显示给中文用户
/// * 每次FFI调用都会清除之前的错误消息
/// * 不要在多线程环境中使用，错误消息是全局共享的
/// * 返回的指针不应被修改或释放
/// 
/// # 线程安全性
/// 
/// 此函数不是线程安全的，因为错误消息存储在全局静态变量中。
/// 多线程环境下可能出现错误消息混乱的情况。
#[no_mangle]
pub extern "C" fn nextkey_get_last_error() -> *const c_char {
    unsafe {
        match &LAST_ERROR {
            Some(err) => {
                // 注意：这里返回的指针生命周期受限，调用者应立即复制
                err.as_ptr() as *const c_char
            }
            None => ptr::null(),
        }
    }
}

/// 释放由SDK分配的C字符串
/// 
/// # 功能
/// 
/// 释放由SDK分配并通过输出参数返回的C字符串内存。
/// 所有返回 `char**` 参数的函数（如 `nextkey_login` 的token_out）返回的字符串
/// 都必须使用此函数释放。
/// 
/// # 参数
/// 
/// * `s` - 需要释放的字符串指针，可以为NULL（为NULL时函数不执行任何操作）
/// 
/// # 示例
/// 
/// ```c
/// char* token = NULL;
/// char* expire_at = NULL;
/// NextKeyCardInfo card_info;
/// 
/// int ret = nextkey_login(client, card_key, hwid, ip, &token, &expire_at, &card_info);
/// 
/// if (ret == NEXTKEY_OK) {
///     printf("Token: %s\n", token);
///     
///     // 使用完毕后释放
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
///     
///     // 注意：card_info内部的字符串用nextkey_free_card_info释放
///     nextkey_free_card_info(&card_info);
/// }
/// ```
/// 
/// # 需要释放的字符串
/// 
/// 以下函数返回的字符串需要使用此函数释放：
/// * `nextkey_login` - token_out, expire_at_out
/// * `nextkey_get_cloud_var` - value_out
/// * `nextkey_get_project_info` - uuid_out, name_out, version_out, update_url_out
/// 
/// 以下情况**不应**使用此函数：
/// * `nextkey_get_last_error` 返回的指针（由SDK内部管理）
/// * `NextKeyCardInfo` 结构体（应使用 `nextkey_free_card_info`）
/// 
/// # 注意事项
/// 
/// * 同一个指针只能释放一次，重复释放会导致未定义行为
/// * 释放后应将指针设置为NULL，避免悬垂指针
/// * 传入NULL是安全的，函数会忽略
/// * 不要用标准库的 `free()` 释放SDK返回的字符串
/// 
/// # 线程安全性
/// 
/// 此函数是线程安全的，可以在不同线程中释放不同的字符串。
#[no_mangle]
pub extern "C" fn nextkey_free_string(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            let _ = CString::from_raw(s);
        }
    }
}

/// 释放卡密信息结构体内部的字符串
/// 
/// # 功能
/// 
/// 释放 `NextKeyCardInfo` 结构体内部的 `card_key` 和 `custom_data` 字段指向的内存。
/// 注意：此函数只释放结构体内部的字符串指针，不释放结构体本身。
/// 
/// # 参数
/// 
/// * `card_info` - 指向 `NextKeyCardInfo` 结构体的指针，可以为NULL（为NULL时函数不执行任何操作）
/// 
/// # 示例
/// 
/// ```c
/// char* token = NULL;
/// char* expire_at = NULL;
/// NextKeyCardInfo card_info;  // 在栈上分配
/// 
/// int ret = nextkey_login(client, card_key, hwid, ip, &token, &expire_at, &card_info);
/// 
/// if (ret == NEXTKEY_OK) {
///     printf("卡密: %s\n", card_info.card_key);
///     printf("专属信息: %s\n", card_info.custom_data);
///     
///     // 释放结构体内部的字符串
///     nextkey_free_card_info(&card_info);
///     
///     // 释放其他字符串
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
///     
///     // 结构体本身在栈上，会自动释放
/// }
/// ```
/// 
/// # 内存管理详解
/// 
/// `NextKeyCardInfo` 结构体通常在栈上分配：
/// ```c
/// NextKeyCardInfo card_info;  // 栈分配，无需手动释放
/// ```
/// 
/// 但结构体内的 `card_key` 和 `custom_data` 字段是由SDK在堆上分配的：
/// - `card_key` - 需要释放
/// - `custom_data` - 需要释放
/// - 其他字段 (`id`, `activated`, `duration`) - 是值类型，无需释放
/// 
/// # 释放行为
/// 
/// * 释放 `card_key` 指向的内存（如果不为NULL）
/// * 释放 `custom_data` 指向的内存（如果不为NULL）
/// * 将这两个指针设置为NULL，防止悬垂指针
/// * **不**释放结构体本身的内存
/// 
/// # 注意事项
/// 
/// * 必须在使用完 `NextKeyCardInfo` 后调用此函数
/// * 同一个结构体只能释放一次
/// * 释放后不要再访问 `card_key` 和 `custom_data` 字段
/// * 传入NULL是安全的，函数会忽略
/// * 不要对未初始化的结构体调用此函数
/// 
/// # 线程安全性
/// 
/// 此函数是线程安全的，可以在不同线程中释放不同的结构体。
#[no_mangle]
pub extern "C" fn nextkey_free_card_info(card_info: *mut NextKeyCardInfo) {
    if !card_info.is_null() {
        unsafe {
            let info = &mut *card_info;
            if !info.card_key.is_null() {
                let _ = CString::from_raw(info.card_key);
                info.card_key = ptr::null_mut();
            }
            if !info.custom_data.is_null() {
                let _ = CString::from_raw(info.custom_data);
                info.custom_data = ptr::null_mut();
            }
        }
    }
}

