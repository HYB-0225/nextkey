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

/// 全局错误消息存储
/// 
/// 用于存储最近一次FFI调用的错误信息，供 `nextkey_get_last_error` 读取
static mut LAST_ERROR: Option<String> = None;

/// 设置全局错误消息
fn set_last_error(err: String) {
    unsafe {
        LAST_ERROR = Some(err);
    }
}

/// 清除全局错误消息
fn clear_last_error() {
    unsafe {
        LAST_ERROR = None;
    }
}

/// 从C字符串转换为Rust字符串
/// 
/// # Safety
/// 
/// `c_str` 必须是有效的C字符串指针或NULL
unsafe fn c_str_to_string(c_str: *const c_char) -> Option<String> {
    if c_str.is_null() {
        return None;
    }
    CStr::from_ptr(c_str).to_str().ok().map(|s| s.to_string())
}

/// 将Rust字符串转换为C字符串指针
/// 
/// 返回的指针需要调用者使用 `nextkey_free_string` 释放
fn string_to_c_char(s: String) -> *mut c_char {
    match CString::new(s) {
        Ok(c_str) => c_str.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

/// 创建NextKey客户端实例（使用默认AES-256-GCM加密）
/// 
/// 使用AES-256-GCM加密方案创建客户端，这是最安全的加密方案，推荐用于生产环境。
/// 
/// # 参数
/// 
/// * `server_url` - 服务器地址，例如 "http://localhost:8080"
/// * `project_uuid` - 项目UUID，从管理后台获取
/// * `aes_key` - AES密钥（64位十六进制字符串）
/// 
/// # 返回值
/// 
/// 成功返回客户端实例指针，失败返回NULL（通过 `nextkey_get_last_error` 获取错误信息）。
/// 返回的指针必须调用 `nextkey_client_free` 释放。
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new(
///     "http://localhost:8080",
///     "550e8400-e29b-41d4-a716-446655440000",
///     "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037"
/// );
/// if (!client) {
///     printf("失败: %s\n", nextkey_get_last_error());
///     return -1;
/// }
/// nextkey_client_free(client);
/// ```
/// 
/// # 线程安全性
/// 
/// 函数本身线程安全，但返回的客户端实例不可在多线程间共享。
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

/// 创建NextKey客户端实例（指定加密方案）
/// 
/// 支持自定义加密方案以满足不同的安全需求和性能要求。
/// **注意**: 加密方案必须与服务器项目配置一致。
/// 
/// # 参数
/// 
/// * `server_url` - 服务器地址
/// * `project_uuid` - 项目UUID
/// * `aes_key` - 加密密钥（格式根据scheme而定）
/// * `scheme` - 加密方案名称（不区分大小写）:
///   - `"aes-256-gcm"` - AES-256-GCM（推荐，高安全）
///   - `"chacha20-poly1305"` - ChaCha20-Poly1305（高安全，移动端友好）
///   - `"rc4"` - RC4（中等安全，快速）
///   - `"xor"` - XOR（低安全，仅测试）
///   - `"custom-base64"` - 自定义Base64（低安全，仅调试）
/// 
/// # 返回值
/// 
/// 成功返回客户端实例指针，失败返回NULL。
/// 返回的指针必须调用 `nextkey_client_free` 释放。
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new_with_scheme(
///     "http://localhost:8080", "550e8400-...", "632005a33e...", "rc4"
/// );
/// if (!client) return -1;
/// nextkey_client_free(client);
/// ```
/// 
/// # 线程安全性
/// 
/// 函数本身线程安全，但返回的客户端实例不可在多线程间共享。
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

/// 释放客户端实例资源
/// 
/// 释放由 `nextkey_client_new` 或 `nextkey_client_new_with_scheme` 创建的客户端实例。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针（可为NULL）
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new(url, uuid, key);
/// nextkey_client_free(client);
/// client = NULL; // 推荐：避免悬垂指针
/// ```
/// 
/// # 注意
/// 
/// 同一指针只能释放一次，重复释放导致未定义行为。传入NULL是安全的。
/// 
/// # 线程安全性
/// 
/// 不是线程安全的，确保释放时没有其他线程使用该客户端实例。
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
/// 存储卡密的详细信息，由 `nextkey_login` 返回。使用C ABI布局。
/// 
/// # 字段
/// 
/// * `id` - 卡密ID（数据库唯一标识）
/// * `card_key` - 卡密字符串（SDK分配，需用 `nextkey_free_card_info` 释放）
/// * `activated` - 激活状态（0=未激活，1=已激活）
/// * `duration` - 有效时长（秒），负数表示永久
/// * `custom_data` - 专属信息（SDK分配，需用 `nextkey_free_card_info` 释放）
/// 
/// # 内存管理
/// 
/// 结构体可在栈上分配，但内部的 `card_key` 和 `custom_data` 由SDK在堆上分配，
/// 使用完毕必须调用 `nextkey_free_card_info` 释放。
/// 
/// # 示例
/// 
/// ```c
/// NextKeyCardInfo card_info;
/// char *token, *expire_at;
/// if (nextkey_login(client, key, hwid, ip, &token, &expire_at, &card_info) == NEXTKEY_OK) {
///     printf("卡密ID: %llu\n", card_info.id);
///     nextkey_free_card_info(&card_info);
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
/// }
/// ```
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
/// 使用卡密进行登录验证，返回访问令牌和卡密信息。Token会自动保存到客户端实例中。
/// 
/// # 参数
/// 
/// **输入:**
/// * `client` - 客户端实例指针
/// * `card_key` - 卡密字符串
/// * `hwid` - 硬件ID（可为空字符串）
/// * `ip` - 客户端IP（可为空字符串）
/// 
/// **输出:**
/// * `token_out` - 访问令牌（需用 `nextkey_free_string` 释放）
/// * `expire_at_out` - 令牌过期时间，ISO 8601格式（需用 `nextkey_free_string` 释放）
/// * `card_info_out` - 卡密信息（需用 `nextkey_free_card_info` 释放）
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * `NEXTKEY_ERR_DECRYPT` (-3) - 解密失败
/// * 其他 - 服务器业务错误码
/// 
/// # 示例
/// 
/// ```c
/// char *token, *expire_at;
/// NextKeyCardInfo card_info;
/// 
/// int ret = nextkey_login(client, "ABCD-1234", "hwid123", "192.168.1.1",
///                         &token, &expire_at, &card_info);
/// if (ret == NEXTKEY_OK) {
///     printf("Token: %s\n", token);
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
///     nextkey_free_card_info(&card_info);
/// }
/// ```
/// 
/// # 注意
/// 
/// * Token自动保存，后续API调用会自动使用
/// * 必须先登录才能调用其他需要认证的API
/// * 避免频繁重试，防止触发限流
/// 
/// # 线程安全性
/// 
/// 不是线程安全的，不可在多线程中同时对同一客户端实例调用。
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
/// 向服务器发送心跳，验证Token有效性并保持会话。服务器会检查卡密状态。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针（必须已登录）
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败（未登录或Token失效）
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他 - 服务器业务错误码
/// 
/// # 示例
/// 
/// ```c
/// while (running) {
///     if (nextkey_heartbeat(client) == NEXTKEY_ERR_AUTH) {
///         printf("会话失效，需重新登录\n");
///         break;
///     }
///     sleep(300); // 每5分钟一次
/// }
/// ```
/// 
/// # 建议
/// 
/// 建议每5-10分钟发送一次。心跳失败时考虑重新登录。
/// 
/// # 线程安全性
/// 
/// 不是线程安全的。
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
/// 根据键名获取项目级别的云变量，用于配置管理、开关控制等场景。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针（必须已登录）
/// * `key` - 云变量键名
/// * `value_out` - 输出云变量值（需用 `nextkey_free_string` 释放）
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 404 - 键不存在
/// * 其他 - 服务器业务错误码
/// 
/// # 示例
/// 
/// ```c
/// char* value = NULL;
/// int ret = nextkey_get_cloud_var(client, "app_version", &value);
/// if (ret == NEXTKEY_OK) {
///     printf("版本: %s\n", value);
///     nextkey_free_string(value);
/// } else if (ret == 404) {
///     printf("键不存在\n");
/// }
/// ```
/// 
/// # 使用场景
/// 
/// 远程配置管理、动态内容分发、全局开关控制等。建议缓存结果减少请求。
/// 
/// # 线程安全性
/// 
/// 不是线程安全的。
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
/// 更新当前卡密的专属信息（custom_data字段），用于存储用户备注、配置等自定义数据。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针（必须已登录）
/// * `custom_data` - 专属信息字符串（可为空字符串表示清空）
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他 - 服务器业务错误码
/// 
/// # 示例
/// 
/// ```c
/// // 设置专属信息
/// nextkey_update_custom_data(client, "user_level:5|coins:1000");
/// 
/// // 清空专属信息
/// nextkey_update_custom_data(client, "");
/// ```
/// 
/// # 使用场景
/// 
/// 存储游戏数据、用户偏好、备注信息、特殊状态等。
/// 
/// # 注意
/// 
/// * 长度限制通常为1024字符，超长会被截断
/// * 更新会覆盖原有内容，不是追加
/// * 可在管理后台查看
/// 
/// # 线程安全性
/// 
/// 不是线程安全的。
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
/// 获取项目的UUID、名称、版本号和更新地址，用于版本检测、更新提示等场景。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针（必须已登录）
/// * `uuid_out` - 输出项目UUID（需用 `nextkey_free_string` 释放）
/// * `name_out` - 输出项目名称（需用 `nextkey_free_string` 释放）
/// * `version_out` - 输出版本号（需用 `nextkey_free_string` 释放）
/// * `update_url_out` - 输出更新地址（需用 `nextkey_free_string` 释放）
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他 - 服务器业务错误码
/// 
/// # 示例
/// 
/// ```c
/// char *uuid, *name, *version, *update_url;
/// 
/// if (nextkey_get_project_info(client, &uuid, &name, &version, &update_url) == NEXTKEY_OK) {
///     printf("项目: %s v%s\n", name, version);
///     if (strcmp(version, LOCAL_VERSION) > 0) {
///         printf("新版本: %s\n", update_url);
///     }
///     nextkey_free_string(uuid);
///     nextkey_free_string(name);
///     nextkey_free_string(version);
///     nextkey_free_string(update_url);
/// }
/// ```
/// 
/// # 使用场景
/// 
/// 版本更新检测、显示项目信息、自动更新功能等。
/// 
/// # 注意
/// 
/// 所有输出字符串都需释放。update_url可能为空字符串。
/// 
/// # 线程安全性
/// 
/// 不是线程安全的。
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
/// 解除卡密与硬件ID的绑定，解绑后可在新设备上重新激活。
/// 
/// # 参数
/// 
/// * `client` - 客户端实例指针（必须已登录）
/// * `card_key` - 要解绑的卡密字符串
/// * `hwid` - 要解绑的硬件ID
/// 
/// # 返回值
/// 
/// * `NEXTKEY_OK` (0) - 成功
/// * `NEXTKEY_ERR_INVALID_PARAM` (-1) - 参数无效
/// * `NEXTKEY_ERR_AUTH` (-401) - 认证失败或无权限
/// * `NEXTKEY_ERR_NETWORK` (-2) - 网络错误
/// * 其他 - 服务器业务错误码
/// 
/// # 示例
/// 
/// ```c
/// if (nextkey_unbind_hwid(client, "ABCD-1234", "old_hwid") == NEXTKEY_OK) {
///     printf("解绑成功\n");
/// } else {
///     printf("失败: %s\n", nextkey_get_last_error());
/// }
/// ```
/// 
/// # 使用场景
/// 
/// 用户更换设备、设备损坏、解除错误绑定等。
/// 
/// # 注意
/// 
/// * 可能需要特定权限
/// * 解绑后下次登录会重新绑定新设备
/// * 频繁解绑可能触发风控
/// * 卡密和HWID必须完全匹配
/// 
/// # 线程安全性
/// 
/// 不是线程安全的。
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
/// 获取最近一次FFI函数调用失败时的详细错误信息。
/// 
/// # 返回值
/// 
/// 有错误时返回错误消息指针（UTF-8编码），无错误时返回NULL。
/// 
/// # 内存管理
/// 
/// 返回的指针由SDK内部管理，**不需要**调用 `nextkey_free_string` 释放。
/// 指针在下次调用任何FFI函数前有效，需要长期保存请复制字符串内容。
/// 
/// # 示例
/// 
/// ```c
/// NextKeyClient* client = nextkey_client_new(url, uuid, invalid_key);
/// if (!client) {
///     const char* error = nextkey_get_last_error();
///     if (error) {
///         printf("失败: %s\n", error);
///         // 长期保存需复制: strncpy(buf, error, sizeof(buf) - 1);
///     }
/// }
/// ```
/// 
/// # 注意
/// 
/// * 错误消息为中文
/// * 每次FFI调用会清除之前的错误
/// * 不要修改或释放返回的指针
/// * 不是线程安全的（全局共享）
/// 
/// # 线程安全性
/// 
/// 不是线程安全的，多线程环境下可能混乱。
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
/// 释放通过输出参数返回的字符串内存（如 `nextkey_login` 的 token_out）。
/// 
/// # 参数
/// 
/// * `s` - 字符串指针（可为NULL）
/// 
/// # 示例
/// 
/// ```c
/// char *token, *expire_at;
/// NextKeyCardInfo card_info;
/// 
/// if (nextkey_login(client, key, hwid, ip, &token, &expire_at, &card_info) == NEXTKEY_OK) {
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
///     nextkey_free_card_info(&card_info); // card_info内部字符串用此函数
/// }
/// ```
/// 
/// # 需要释放的字符串
/// 
/// * `nextkey_login` - token_out, expire_at_out
/// * `nextkey_get_cloud_var` - value_out
/// * `nextkey_get_project_info` - uuid_out, name_out, version_out, update_url_out
/// 
/// **不应**使用此函数：
/// * `nextkey_get_last_error` 返回的指针（SDK内部管理）
/// * `NextKeyCardInfo` 结构体（应使用 `nextkey_free_card_info`）
/// 
/// # 注意
/// 
/// * 同一指针只能释放一次，重复释放导致未定义行为
/// * 释放后应设置为NULL
/// * 传入NULL是安全的
/// * 不要用标准库 `free()` 释放SDK字符串
/// 
/// # 线程安全性
/// 
/// 线程安全，可在不同线程中释放不同字符串。
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
/// 释放 `NextKeyCardInfo` 结构体内的 `card_key` 和 `custom_data` 字段的内存。
/// 只释放内部字符串指针，不释放结构体本身。
/// 
/// # 参数
/// 
/// * `card_info` - 指向 `NextKeyCardInfo` 的指针（可为NULL）
/// 
/// # 示例
/// 
/// ```c
/// char *token, *expire_at;
/// NextKeyCardInfo card_info;  // 栈上分配
/// 
/// if (nextkey_login(client, key, hwid, ip, &token, &expire_at, &card_info) == NEXTKEY_OK) {
///     printf("卡密: %s\n", card_info.card_key);
///     
///     nextkey_free_card_info(&card_info); // 释放内部字符串
///     nextkey_free_string(token);
///     nextkey_free_string(expire_at);
///     // card_info结构体本身在栈上，自动释放
/// }
/// ```
/// 
/// # 内存管理
/// 
/// 结构体通常在栈上分配，但内部字段由SDK在堆上分配：
/// - `card_key` - 需要释放
/// - `custom_data` - 需要释放
/// - 其他字段（`id`, `activated`, `duration`）- 值类型，无需释放
/// 
/// # 释放行为
/// 
/// 释放 `card_key` 和 `custom_data` 指向的内存并将指针设为NULL。
/// **不**释放结构体本身。
/// 
/// # 注意
/// 
/// * 使用完毕必须调用此函数
/// * 同一结构体只能释放一次
/// * 释放后不要再访问 `card_key` 和 `custom_data`
/// * 传入NULL是安全的
/// * 不要对未初始化的结构体调用
/// 
/// # 线程安全性
/// 
/// 线程安全，可在不同线程中释放不同结构体。
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

