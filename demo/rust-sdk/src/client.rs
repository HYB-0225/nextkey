use crate::crypto::Crypto;
use anyhow::{Context, Result};
use reqwest::blocking::Client as HttpClient;
use serde::{Deserialize, Serialize};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Debug, Serialize, Deserialize)]
pub struct EncryptedRequest {
    timestamp: u64,
    nonce: String,
    data: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EncryptedResponse {
    nonce: String,
    data: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct InternalRequest<T> {
    nonce: String,
    timestamp: u64,
    data: T,
}

#[derive(Debug, Serialize, Deserialize)]
struct InternalResponse<T> {
    nonce: String,
    timestamp: u64,
    data: T,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ApiResponse<T> {
    pub code: i32,
    pub message: String,
    pub data: Option<T>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LoginRequest {
    project_uuid: String,
    card_key: String,
    #[serde(skip_serializing_if = "String::is_empty")]
    hwid: String,
    #[serde(skip_serializing_if = "String::is_empty")]
    ip: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CardInfo {
    pub id: u64,
    pub card_key: String,
    pub activated: bool,
    pub duration: i64,
    pub custom_data: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LoginData {
    pub token: String,
    pub expire_at: String,
    pub card: CardInfo,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CloudVarData {
    pub id: u64,
    pub project_id: u64,
    pub key: String,
    pub value: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectInfo {
    pub uuid: String,
    pub name: String,
    pub version: String,
    pub update_url: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct CustomDataRequest {
    custom_data: String,
}

pub struct NextKeyClient {
    server_url: String,
    project_uuid: String,
    crypto: Crypto,
    token: Option<String>,
    http_client: HttpClient,
}

impl NextKeyClient {
    pub fn new(server_url: &str, project_uuid: &str, aes_key: &str) -> Result<Self> {
        let crypto = Crypto::new(aes_key)?;
        let http_client = HttpClient::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()
            .context("创建HTTP客户端失败")?;

        Ok(Self {
            server_url: server_url.trim_end_matches('/').to_string(),
            project_uuid: project_uuid.to_string(),
            crypto,
            token: None,
            http_client,
        })
    }

    pub fn new_with_scheme(server_url: &str, project_uuid: &str, aes_key: &str, scheme: crate::crypto::EncryptionScheme) -> Result<Self> {
        let crypto = Crypto::new_with_scheme(aes_key, scheme)?;
        let http_client = HttpClient::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()
            .context("创建HTTP客户端失败")?;

        Ok(Self {
            server_url: server_url.trim_end_matches('/').to_string(),
            project_uuid: project_uuid.to_string(),
            crypto,
            token: None,
            http_client,
        })
    }

    /// 获取当前token
    pub fn token(&self) -> Option<&str> {
        self.token.as_deref()
    }

    /// 设置token
    pub fn set_token(&mut self, token: String) {
        self.token = Some(token);
    }

    /// 发送加密请求并验证响应Nonce
    fn make_encrypted_request<T: Serialize, R: for<'de> Deserialize<'de>>(
        &self,
        endpoint: &str,
        data: &T,
        method: &str,
    ) -> Result<ApiResponse<R>> {
        // 生成并记住请求nonce
        let request_nonce = Crypto::generate_nonce();
        let request_timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)?
            .as_secs();

        // 包装内层数据，嵌入nonce和timestamp
        let internal_data = InternalRequest {
            nonce: request_nonce.clone(),
            timestamp: request_timestamp,
            data,
        };

        // 加密请求数据
        let json_data = serde_json::to_string(&internal_data)?;
        let encrypted_data = self.crypto.encrypt(&json_data)?;

        let req_body = EncryptedRequest {
            timestamp: request_timestamp,
            nonce: request_nonce.clone(),
            data: encrypted_data,
        };

        // 构建请求
        let url = format!("{}{}", self.server_url, endpoint);
        let mut request_builder = match method {
            "POST" => self.http_client.post(&url),
            "GET" => self.http_client.get(&url),
            _ => anyhow::bail!("不支持的HTTP方法: {}", method),
        };

        // 添加token（如果有）
        if let Some(token) = &self.token {
            request_builder = request_builder.header("Authorization", format!("Bearer {}", token));
        }

        // 发送请求
        let response = request_builder
            .json(&req_body)
            .send()
            .context("发送请求失败")?;

        if !response.status().is_success() {
            anyhow::bail!("HTTP错误: {}", response.status());
        }

        let resp_json: EncryptedResponse = response.json().context("解析响应失败")?;

        // 验证外层响应nonce
        if resp_json.nonce != request_nonce {
            anyhow::bail!("外层响应Nonce不匹配，可能遭受重放攻击！");
        }

        // 解密响应数据
        let decrypted = self.crypto.decrypt(&resp_json.data)?;
        let internal_response: InternalResponse<ApiResponse<R>> = serde_json::from_str(&decrypted)?;

        // 验证内层响应nonce（双重验证）
        if internal_response.nonce != request_nonce {
            anyhow::bail!("内层响应Nonce不匹配，响应数据可能被篡改！");
        }

        // 验证服务器时间戳
        let current_time = SystemTime::now()
            .duration_since(UNIX_EPOCH)?
            .as_secs();
        let time_diff = if current_time > internal_response.timestamp {
            current_time - internal_response.timestamp
        } else {
            internal_response.timestamp - current_time
        };

        if time_diff > 300 {
            anyhow::bail!("响应时间戳异常，可能遭受离线攻击！时间差: {}秒", time_diff);
        }

        Ok(internal_response.data)
    }

    /// 登录
    pub fn login(
        &mut self,
        card_key: &str,
        hwid: &str,
        ip: &str,
    ) -> Result<ApiResponse<LoginData>> {
        let login_data = LoginRequest {
            project_uuid: self.project_uuid.clone(),
            card_key: card_key.to_string(),
            hwid: hwid.to_string(),
            ip: ip.to_string(),
        };

        let result: ApiResponse<LoginData> = self.make_encrypted_request("/api/auth/login", &login_data, "POST")?;

        // 保存token
        if result.code == 0 {
            if let Some(ref data) = result.data {
                self.token = Some(data.token.clone());
            }
        }

        Ok(result)
    }

    /// 心跳验证
    pub fn heartbeat(&self) -> Result<ApiResponse<serde_json::Value>> {
        if self.token.is_none() {
            anyhow::bail!("未登录，请先调用login");
        }

        let empty = serde_json::json!({});
        self.make_encrypted_request("/api/heartbeat", &empty, "POST")
    }

    /// 获取云变量
    pub fn get_cloud_var(&self, key: &str) -> Result<ApiResponse<CloudVarData>> {
        if self.token.is_none() {
            anyhow::bail!("未登录，请先调用login");
        }

        let empty = serde_json::json!({});
        self.make_encrypted_request(&format!("/api/cloud-var/{}", key), &empty, "GET")
    }

    /// 更新专属信息
    pub fn update_custom_data(&self, custom_data: &str) -> Result<ApiResponse<serde_json::Value>> {
        if self.token.is_none() {
            anyhow::bail!("未登录，请先调用login");
        }

        let data = CustomDataRequest {
            custom_data: custom_data.to_string(),
        };

        self.make_encrypted_request("/api/card/custom-data", &data, "POST")
    }

    /// 获取项目信息
    pub fn get_project_info(&self) -> Result<ApiResponse<ProjectInfo>> {
        if self.token.is_none() {
            anyhow::bail!("未登录，请先调用login");
        }

        let empty = serde_json::json!({});
        self.make_encrypted_request("/api/project/info", &empty, "GET")
    }

    /// 解绑HWID
    pub fn unbind_hwid(&self, card_key: &str, hwid: &str) -> Result<ApiResponse<serde_json::Value>> {
        if self.token.is_none() {
            anyhow::bail!("未登录，请先调用login");
        }

        let data = serde_json::json!({
            "project_uuid": self.project_uuid,
            "card_key": card_key,
            "hwid": hwid,
        });

        self.make_encrypted_request("/api/card/unbind", &data, "POST")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_client_creation() {
        let key = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037";
        let client = NextKeyClient::new("http://localhost:8080", "test-uuid", key);
        assert!(client.is_ok());
    }

    #[test]
    fn test_token_management() {
        let key = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037";
        let mut client = NextKeyClient::new("http://localhost:8080", "test-uuid", key).unwrap();

        assert!(client.token().is_none());

        client.set_token("test-token".to_string());
        assert_eq!(client.token(), Some("test-token"));
    }
}

