use crate::client::NextKeyClient;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::ptr;

// 错误码定义
pub const NEXTKEY_OK: i32 = 0;
pub const NEXTKEY_ERR_INVALID_PARAM: i32 = -1;
pub const NEXTKEY_ERR_NETWORK: i32 = -2;
pub const NEXTKEY_ERR_AUTH: i32 = -401;
pub const NEXTKEY_ERR_DECRYPT: i32 = -3;
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

/// 创建NextKey客户端
#[no_mangle]
pub extern "C" fn nextkey_client_new(
    server_url: *const c_char,
    project_uuid: *const c_char,
    aes_key: *const c_char,
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

    match NextKeyClient::new(&server_url, &project_uuid, &aes_key) {
        Ok(client) => Box::into_raw(Box::new(client)),
        Err(e) => {
            set_last_error(format!("创建客户端失败: {}", e));
            ptr::null_mut()
        }
    }
}

/// 释放客户端
#[no_mangle]
pub extern "C" fn nextkey_client_free(client: *mut NextKeyClient) {
    if !client.is_null() {
        unsafe {
            let _ = Box::from_raw(client);
        }
    }
}

/// 卡密信息结构体（C ABI）
#[repr(C)]
pub struct NextKeyCardInfo {
    pub id: u64,
    pub card_key: *mut c_char,
    pub activated: i32,
    pub duration: i64,
    pub custom_data: *mut c_char,
}

/// 登录
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

/// 心跳
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

/// 获取云变量
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

/// 更新专属信息
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

/// 获取最后的错误消息
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

/// 释放由库分配的C字符串
#[no_mangle]
pub extern "C" fn nextkey_free_string(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            let _ = CString::from_raw(s);
        }
    }
}

/// 释放卡密信息结构体
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

