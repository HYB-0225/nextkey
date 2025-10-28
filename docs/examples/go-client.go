package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ServerURL   = "http://localhost:8080"
	ProjectUUID = "your-project-uuid"
	AESKey      = "your-aes-key-32-bytes-long!!" // 32字节密钥
)

type Client struct {
	httpClient *http.Client
	aesKey     []byte
	token      string
}

type EncryptedRequest struct {
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Data      string `json:"data"`
}

type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		aesKey:     []byte(AESKey),
	}
}

func (c *Client) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Client) generateNonce() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (c *Client) Login(cardKey, hwid string) error {
	loginReq := map[string]string{
		"project_uuid": ProjectUUID,
		"card_key":     cardKey,
		"hwid":         hwid,
	}

	jsonData, _ := json.Marshal(loginReq)
	encryptedData, err := c.encrypt(string(jsonData))
	if err != nil {
		return err
	}

	req := EncryptedRequest{
		Timestamp: time.Now().Unix(),
		Nonce:     c.generateNonce(),
		Data:      encryptedData,
	}

	reqBody, _ := json.Marshal(req)
	resp, err := c.httpClient.Post(ServerURL+"/api/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Code != 0 {
		return fmt.Errorf("登录失败: %s", result.Message)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(result.Data, &loginResp)
	c.token = loginResp.Token

	fmt.Println("登录成功! Token:", c.token)
	return nil
}

func (c *Client) Heartbeat() error {
	req, _ := http.NewRequest("POST", ServerURL+"/api/heartbeat", nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Code != 0 {
		return fmt.Errorf("心跳失败: %s", result.Message)
	}

	fmt.Println("心跳成功")
	return nil
}

func (c *Client) GetCloudVar(key string) (string, error) {
	req, _ := http.NewRequest("GET", ServerURL+"/api/cloud-var/"+key, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("获取失败: %s", result.Message)
	}

	var cloudVar struct {
		Value string `json:"value"`
	}
	json.Unmarshal(result.Data, &cloudVar)

	return cloudVar.Value, nil
}

func main() {
	client := NewClient()

	if err := client.Login("your-card-key", "your-hwid"); err != nil {
		fmt.Println("登录失败:", err)
		return
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			client.Heartbeat()
		}
	}()

	value, err := client.GetCloudVar("test_key")
	if err != nil {
		fmt.Println("获取云变量失败:", err)
	} else {
		fmt.Println("云变量值:", value)
	}

	select {}
}
