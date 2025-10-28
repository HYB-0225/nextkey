import base64
import json
import time
import secrets
from Crypto.Cipher import AES
import requests

SERVER_URL = "http://localhost:8080"
PROJECT_UUID = "your-project-uuid"
AES_KEY = b"your-aes-key-32-bytes-long!!"  # 32字节密钥


class NextKeyClient:
    def __init__(self):
        self.aes_key = AES_KEY
        self.token = None
        self.session = requests.Session()

    def encrypt(self, plaintext):
        cipher = AES.new(self.aes_key, AES.MODE_GCM)
        ciphertext, tag = cipher.encrypt_and_digest(plaintext.encode())
        return base64.b64encode(cipher.nonce + tag + ciphertext).decode()

    def generate_nonce(self):
        return secrets.token_urlsafe(24)

    def login(self, card_key, hwid=""):
        login_data = {
            "project_uuid": PROJECT_UUID,
            "card_key": card_key,
            "hwid": hwid
        }

        json_data = json.dumps(login_data)
        encrypted_data = self.encrypt(json_data)

        req_body = {
            "timestamp": int(time.time()),
            "nonce": self.generate_nonce(),
            "data": encrypted_data
        }

        resp = self.session.post(f"{SERVER_URL}/api/auth/login", json=req_body)
        result = resp.json()

        if result["code"] != 0:
            raise Exception(f"登录失败: {result['message']}")

        self.token = result["data"]["token"]
        print(f"登录成功! Token: {self.token}")
        return result["data"]

    def heartbeat(self):
        headers = {"Authorization": f"Bearer {self.token}"}
        resp = self.session.post(f"{SERVER_URL}/api/heartbeat", headers=headers)
        result = resp.json()

        if result["code"] != 0:
            raise Exception(f"心跳失败: {result['message']}")

        print("心跳成功")
        return True

    def get_cloud_var(self, key):
        headers = {"Authorization": f"Bearer {self.token}"}
        resp = self.session.get(f"{SERVER_URL}/api/cloud-var/{key}", headers=headers)
        result = resp.json()

        if result["code"] != 0:
            raise Exception(f"获取失败: {result['message']}")

        return result["data"]["value"]

    def update_custom_data(self, custom_data):
        headers = {"Authorization": f"Bearer {self.token}"}
        data = {"custom_data": json.dumps(custom_data)}
        resp = self.session.post(f"{SERVER_URL}/api/card/custom-data", json=data, headers=headers)
        result = resp.json()

        if result["code"] != 0:
            raise Exception(f"更新失败: {result['message']}")

        print("更新专属信息成功")
        return True


if __name__ == "__main__":
    client = NextKeyClient()

    try:
        client.login("your-card-key", "your-hwid")

        import threading
        def heartbeat_loop():
            while True:
                time.sleep(30)
                try:
                    client.heartbeat()
                except Exception as e:
                    print(f"心跳错误: {e}")

        heartbeat_thread = threading.Thread(target=heartbeat_loop, daemon=True)
        heartbeat_thread.start()

        value = client.get_cloud_var("test_key")
        print(f"云变量值: {value}")

        client.update_custom_data({"user_level": 5, "points": 1000})

    except Exception as e:
        print(f"错误: {e}")

