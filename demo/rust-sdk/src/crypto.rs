use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    Aes256Gcm,
};
use aes_gcm::aead::generic_array::GenericArray;
use anyhow::{Context, Result};
use base64::{engine::general_purpose, Engine as _};
use rand::RngCore;

pub struct Crypto {
    key: [u8; 32],
}

impl Crypto {
    /// 创建加密实例，匹配服务端密钥处理逻辑
    pub fn new(aes_key: &str) -> Result<Self> {
        let key = Self::prepare_key(aes_key)?;
        Ok(Self { key })
    }

    /// 准备AES密钥，匹配Go后端的decodeKey逻辑
    fn prepare_key(key_str: &str) -> Result<[u8; 32]> {
        // 尝试base64解码
        if let Ok(decoded) = general_purpose::STANDARD.decode(key_str) {
            if decoded.len() == 32 {
                let mut key = [0u8; 32];
                key.copy_from_slice(&decoded);
                return Ok(key);
            }
        }

        // 64字符时，取前32字符的UTF-8字节（匹配Go的[]byte(key)[:32]）
        if key_str.len() == 64 {
            let mut key = [0u8; 32];
            key.copy_from_slice(&key_str.as_bytes()[..32]);
            return Ok(key);
        }

        // 其他情况直接编码
        let key_bytes = key_str.as_bytes();
        if key_bytes.len() != 32 {
            anyhow::bail!(
                "AES密钥长度错误，应为32字节，实际: {}",
                key_bytes.len()
            );
        }

        let mut key = [0u8; 32];
        key.copy_from_slice(key_bytes);
        Ok(key)
    }

    /// AES-GCM加密
    pub fn encrypt(&self, plaintext: &str) -> Result<String> {
        let cipher = Aes256Gcm::new(&self.key.into());

        // 生成12字节nonce (GCM标准)
        let mut nonce_bytes = [0u8; 12];
        OsRng.fill_bytes(&mut nonce_bytes);
        let nonce = GenericArray::from_slice(&nonce_bytes);

        // 加密
        let ciphertext = cipher
            .encrypt(nonce, plaintext.as_bytes())
            .map_err(|_| anyhow::anyhow!("加密失败"))?;

        // 格式: nonce + ciphertext (ciphertext已包含tag)
        let mut encrypted = Vec::with_capacity(12 + ciphertext.len());
        encrypted.extend_from_slice(&nonce_bytes);
        encrypted.extend_from_slice(&ciphertext);

        Ok(general_purpose::STANDARD.encode(&encrypted))
    }

    /// AES-GCM解密
    pub fn decrypt(&self, ciphertext: &str) -> Result<String> {
        let cipher = Aes256Gcm::new(&self.key.into());

        // Base64解码
        let data = general_purpose::STANDARD
            .decode(ciphertext)
            .context("Base64解码失败")?;

        if data.len() < 12 {
            anyhow::bail!("密文长度不足");
        }

        // 提取nonce和密文
        let nonce = GenericArray::from_slice(&data[..12]);
        let ciphertext = &data[12..];

        // 解密
        let plaintext = cipher
            .decrypt(nonce, ciphertext)
            .map_err(|_| anyhow::anyhow!("解密失败"))?;

        String::from_utf8(plaintext).context("解密后的数据不是有效的UTF-8")
    }

    /// 生成随机nonce字符串（用于请求）
    pub fn generate_nonce() -> String {
        let mut bytes = [0u8; 24];
        OsRng.fill_bytes(&mut bytes);
        general_purpose::URL_SAFE_NO_PAD.encode(&bytes)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_encrypt_decrypt() {
        let key = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037";
        let crypto = Crypto::new(key).unwrap();

        let plaintext = r#"{"test": "data"}"#;
        let encrypted = crypto.encrypt(plaintext).unwrap();
        let decrypted = crypto.decrypt(&encrypted).unwrap();

        assert_eq!(plaintext, decrypted);
    }

    #[test]
    fn test_key_preparation() {
        // 64字符hex密钥
        let key1 = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037";
        let crypto1 = Crypto::new(key1).unwrap();
        assert_eq!(crypto1.key[..4], b"6320"[..]);

        // 32字节base64密钥
        let key2 = general_purpose::STANDARD.encode(&[0u8; 32]);
        let crypto2 = Crypto::new(&key2).unwrap();
        assert_eq!(crypto2.key, [0u8; 32]);
    }

    #[test]
    fn test_nonce_generation() {
        let nonce1 = Crypto::generate_nonce();
        let nonce2 = Crypto::generate_nonce();

        assert_ne!(nonce1, nonce2);
        assert!(nonce1.len() > 20);
    }
}

