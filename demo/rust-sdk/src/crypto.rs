use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    Aes256Gcm,
};
use aes_gcm::aead::generic_array::GenericArray;
use anyhow::{Context, Result};
use base64::{engine::general_purpose, Engine as _};
use rand::RngCore;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum EncryptionScheme {
    AES256GCM,
    ChaCha20Poly1305,
    RC4,
    CustomBase64,
    XOR,
}

pub struct Crypto {
    scheme: EncryptionScheme,
    key_data: Vec<u8>,
}

impl Crypto {
    /// 创建加密实例（默认AES-256-GCM）
    pub fn new(aes_key: &str) -> Result<Self> {
        Self::new_with_scheme(aes_key, EncryptionScheme::AES256GCM)
    }

    /// 创建指定加密方案的实例
    pub fn new_with_scheme(key: &str, scheme: EncryptionScheme) -> Result<Self> {
        let key_data = match scheme {
            EncryptionScheme::AES256GCM => Self::prepare_aes_key(key)?.to_vec(),
            EncryptionScheme::ChaCha20Poly1305 => Self::prepare_chacha20_key(key)?.to_vec(),
            EncryptionScheme::RC4 | EncryptionScheme::XOR => {
                // RC4和XOR使用hex或直接字节
                hex::decode(key).unwrap_or_else(|_| key.as_bytes().to_vec())
            }
            EncryptionScheme::CustomBase64 => {
                // 自定义Base64需要64字符的映射表
                if key.len() != 64 {
                    anyhow::bail!("自定义Base64密钥必须是64字符");
                }
                key.as_bytes().to_vec()
            }
        };

        Ok(Self { scheme, key_data })
    }

    /// 准备AES密钥，匹配Go后端的decodeKey逻辑
    fn prepare_aes_key(key_str: &str) -> Result<[u8; 32]> {
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

    /// 准备ChaCha20密钥，匹配Go后端的decodeKey逻辑
    fn prepare_chacha20_key(key_str: &str) -> Result<[u8; 32]> {
        // 尝试hex解码
        if let Ok(decoded) = hex::decode(key_str) {
            if decoded.len() == 32 {
                let mut key = [0u8; 32];
                key.copy_from_slice(&decoded);
                return Ok(key);
            }
        }

        // 尝试base64解码
        if let Ok(decoded) = general_purpose::STANDARD.decode(key_str) {
            if decoded.len() == 32 {
                let mut key = [0u8; 32];
                key.copy_from_slice(&decoded);
                return Ok(key);
            }
        }

        // 其他情况直接编码
        let key_bytes = key_str.as_bytes();
        if key_bytes.len() != 32 {
            anyhow::bail!(
                "ChaCha20密钥长度错误，应为32字节，实际: {}",
                key_bytes.len()
            );
        }

        let mut key = [0u8; 32];
        key.copy_from_slice(key_bytes);
        Ok(key)
    }

    /// 加密
    pub fn encrypt(&self, plaintext: &str) -> Result<String> {
        match self.scheme {
            EncryptionScheme::AES256GCM => self.encrypt_aes_gcm(plaintext),
            EncryptionScheme::ChaCha20Poly1305 => self.encrypt_chacha20(plaintext),
            EncryptionScheme::RC4 => self.encrypt_rc4(plaintext),
            EncryptionScheme::XOR => self.encrypt_xor(plaintext),
            EncryptionScheme::CustomBase64 => self.encrypt_custom_base64(plaintext),
        }
    }

    /// AES-GCM加密
    fn encrypt_aes_gcm(&self, plaintext: &str) -> Result<String> {
        let key: [u8; 32] = self.key_data[..32].try_into()?;
        let cipher = Aes256Gcm::new(&key.into());

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

    /// ChaCha20-Poly1305加密
    fn encrypt_chacha20(&self, plaintext: &str) -> Result<String> {
        use chacha20poly1305::{ChaCha20Poly1305, KeyInit, AeadCore};
        use chacha20poly1305::aead::Aead;
        
        let key: [u8; 32] = self.key_data[..32].try_into()?;
        let cipher = ChaCha20Poly1305::new(&key.into());

        // 生成12字节nonce
        let nonce = ChaCha20Poly1305::generate_nonce(&mut OsRng);

        // 加密
        let ciphertext = cipher
            .encrypt(&nonce, plaintext.as_bytes())
            .map_err(|_| anyhow::anyhow!("ChaCha20加密失败"))?;

        // 格式: nonce + ciphertext (ciphertext已包含tag)
        let mut encrypted = Vec::with_capacity(12 + ciphertext.len());
        encrypted.extend_from_slice(&nonce);
        encrypted.extend_from_slice(&ciphertext);

        Ok(general_purpose::STANDARD.encode(&encrypted))
    }

    /// RC4加密
    fn encrypt_rc4(&self, plaintext: &str) -> Result<String> {
        use rc4::{KeyInit, StreamCipher, Rc4};
        use rc4::consts::U32;
        
        // RC4支持1-256字节的密钥，这里使用32字节
        let mut cipher: Rc4<U32> = Rc4::new_from_slice(&self.key_data)
            .map_err(|_| anyhow::anyhow!("RC4密钥长度无效"))?;
        let mut data = plaintext.as_bytes().to_vec();
        cipher.apply_keystream(&mut data);
        
        Ok(general_purpose::STANDARD.encode(&data))
    }

    /// XOR加密
    fn encrypt_xor(&self, plaintext: &str) -> Result<String> {
        let plaintext_bytes = plaintext.as_bytes();
        let mut ciphertext = Vec::with_capacity(plaintext_bytes.len());
        
        for (i, &byte) in plaintext_bytes.iter().enumerate() {
            ciphertext.push(byte ^ self.key_data[i % self.key_data.len()]);
        }
        
        Ok(general_purpose::STANDARD.encode(&ciphertext))
    }

    /// 自定义Base64加密
    fn encrypt_custom_base64(&self, plaintext: &str) -> Result<String> {
        use base64::alphabet::Alphabet;
        use base64::engine::{GeneralPurpose, general_purpose::PAD};
        
        let alphabet = Alphabet::new(std::str::from_utf8(&self.key_data)?)?;
        let engine = GeneralPurpose::new(&alphabet, PAD);
        
        let custom_encoded = engine.encode(plaintext.as_bytes());
        Ok(general_purpose::STANDARD.encode(custom_encoded.as_bytes()))
    }

    /// 解密
    pub fn decrypt(&self, ciphertext: &str) -> Result<String> {
        match self.scheme {
            EncryptionScheme::AES256GCM => self.decrypt_aes_gcm(ciphertext),
            EncryptionScheme::ChaCha20Poly1305 => self.decrypt_chacha20(ciphertext),
            EncryptionScheme::RC4 => self.decrypt_rc4(ciphertext),
            EncryptionScheme::XOR => self.decrypt_xor(ciphertext),
            EncryptionScheme::CustomBase64 => self.decrypt_custom_base64(ciphertext),
        }
    }

    /// AES-GCM解密
    fn decrypt_aes_gcm(&self, ciphertext: &str) -> Result<String> {
        let key: [u8; 32] = self.key_data[..32].try_into()?;
        let cipher = Aes256Gcm::new(&key.into());

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

    /// ChaCha20-Poly1305解密
    fn decrypt_chacha20(&self, ciphertext: &str) -> Result<String> {
        use chacha20poly1305::{ChaCha20Poly1305, KeyInit};
        use chacha20poly1305::aead::Aead;
        
        let key: [u8; 32] = self.key_data[..32].try_into()?;
        let cipher = ChaCha20Poly1305::new(&key.into());

        // Base64解码
        let data = general_purpose::STANDARD
            .decode(ciphertext)
            .context("Base64解码失败")?;

        if data.len() < 12 {
            anyhow::bail!("密文长度不足");
        }

        // 提取nonce和密文
        let nonce = GenericArray::from_slice(&data[..12]);
        let ciphertext_data = &data[12..];

        // 解密
        let plaintext = cipher
            .decrypt(nonce, ciphertext_data)
            .map_err(|_| anyhow::anyhow!("ChaCha20解密失败"))?;

        String::from_utf8(plaintext).context("解密后的数据不是有效的UTF-8")
    }

    /// RC4解密
    fn decrypt_rc4(&self, ciphertext: &str) -> Result<String> {
        use rc4::{KeyInit, StreamCipher, Rc4};
        use rc4::consts::U32;
        
        let data = general_purpose::STANDARD.decode(ciphertext)?;
        // RC4支持1-256字节的密钥，这里使用32字节
        let mut cipher: Rc4<U32> = Rc4::new_from_slice(&self.key_data)
            .map_err(|_| anyhow::anyhow!("RC4密钥长度无效"))?;
        let mut plaintext = data;
        cipher.apply_keystream(&mut plaintext);
        
        String::from_utf8(plaintext).context("解密后的数据不是有效的UTF-8")
    }

    /// XOR解密
    fn decrypt_xor(&self, ciphertext: &str) -> Result<String> {
        let data = general_purpose::STANDARD.decode(ciphertext)?;
        let mut plaintext = Vec::with_capacity(data.len());
        
        for (i, &byte) in data.iter().enumerate() {
            plaintext.push(byte ^ self.key_data[i % self.key_data.len()]);
        }
        
        String::from_utf8(plaintext).context("解密后的数据不是有效的UTF-8")
    }

    /// 自定义Base64解密
    fn decrypt_custom_base64(&self, ciphertext: &str) -> Result<String> {
        use base64::alphabet::Alphabet;
        use base64::engine::{GeneralPurpose, general_purpose::PAD};
        
        // 先解开外层标准Base64
        let custom_encoded = general_purpose::STANDARD.decode(ciphertext)?;
        
        let alphabet = Alphabet::new(std::str::from_utf8(&self.key_data)?)?;
        let engine = GeneralPurpose::new(&alphabet, PAD);
        
        let plaintext = engine.decode(&custom_encoded)?;
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
        assert_eq!(crypto1.key_data[..4], b"6320"[..]);

        // 32字节base64密钥
        let key2 = general_purpose::STANDARD.encode(&[0u8; 32]);
        let crypto2 = Crypto::new(&key2).unwrap();
        assert_eq!(crypto2.key_data, [0u8; 32]);
    }

    #[test]
    fn test_nonce_generation() {
        let nonce1 = Crypto::generate_nonce();
        let nonce2 = Crypto::generate_nonce();

        assert_ne!(nonce1, nonce2);
        assert!(nonce1.len() > 20);
    }
}

