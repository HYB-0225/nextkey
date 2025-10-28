package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/crypto"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type EncryptedResponse struct {
	Nonce string `json:"nonce"`
	Data  string `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func EncryptedSuccess(c *gin.Context, data interface{}) {
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		Error(c, 500, "序列化失败")
		return
	}

	encryptedData, err := crypto.Encrypt(string(jsonData))
	if err != nil {
		Error(c, 500, "加密失败")
		return
	}

	nonce, _ := c.Get("request_nonce")
	nonceStr, _ := nonce.(string)

	c.JSON(200, EncryptedResponse{
		Nonce: nonceStr,
		Data:  encryptedData,
	})
}

func EncryptedError(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		Error(c, 500, "序列化失败")
		return
	}

	encryptedData, err := crypto.Encrypt(string(jsonData))
	if err != nil {
		Error(c, 500, "加密失败")
		return
	}

	nonce, _ := c.Get("request_nonce")
	nonceStr, _ := nonce.(string)

	c.JSON(200, EncryptedResponse{
		Nonce: nonceStr,
		Data:  encryptedData,
	})
}
