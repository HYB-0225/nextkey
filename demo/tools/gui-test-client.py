#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
NextKey GUI 测试客户端
提供可视化界面用于测试 NextKey API 对接
"""

import tkinter as tk
from tkinter import ttk, scrolledtext, messagebox, filedialog
import json
import base64
import time
import secrets
import threading
from datetime import datetime
from Crypto.Cipher import AES, ARC4
import requests
import os
import yaml


class NextKeyClient:
    """NextKey API 客户端"""
    
    def __init__(self, server_url, project_uuid, aes_key, encryption_scheme="aes-256-gcm"):
        self.server_url = server_url.rstrip('/')
        self.project_uuid = project_uuid
        self.encryption_scheme = encryption_scheme.lower()
        self.aes_key = self._prepare_key(aes_key)
        self.token = None
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
    
    def _prepare_key(self, key_str):
        """准备加密密钥 - 根据加密方案处理不同格式"""
        if self.encryption_scheme == "aes-256-gcm":
            # AES-GCM需要32字节密钥
            # 尝试base64解码
            try:
                key_bytes = base64.b64decode(key_str)
                if len(key_bytes) == 32:
                    return key_bytes
            except:
                pass
            
            # 64字符时，取前32字符的UTF-8字节（匹配Go的[]byte(key)[:32]）
            if len(key_str) == 64:
                return key_str[:32].encode('utf-8')
            
            # 其他情况直接编码
            key_bytes = key_str.encode('utf-8')
            if len(key_bytes) != 32:
                raise ValueError(f"AES密钥长度错误，应为32字节，实际: {len(key_bytes)}")
            return key_bytes
        
        elif self.encryption_scheme in ["rc4", "xor"]:
            # RC4和XOR尝试hex解码，否则直接使用字节
            try:
                return bytes.fromhex(key_str)
            except ValueError:
                return key_str.encode('utf-8')
        
        elif self.encryption_scheme == "custom-base64":
            # 自定义Base64需要64字符的映射表
            if len(key_str) != 64:
                raise ValueError(f"自定义Base64密钥必须是64字符，实际: {len(key_str)}")
            return key_str.encode('utf-8')
        
        else:
            raise ValueError(f"不支持的加密方案: {self.encryption_scheme}")
    
    def encrypt(self, plaintext):
        """根据加密方案加密"""
        if self.encryption_scheme == "aes-256-gcm":
            return self._encrypt_aes_gcm(plaintext)
        elif self.encryption_scheme == "rc4":
            return self._encrypt_rc4(plaintext)
        elif self.encryption_scheme == "xor":
            return self._encrypt_xor(plaintext)
        elif self.encryption_scheme == "custom-base64":
            return self._encrypt_custom_base64(plaintext)
        else:
            raise ValueError(f"不支持的加密方案: {self.encryption_scheme}")
    
    def decrypt(self, ciphertext):
        """根据加密方案解密"""
        if self.encryption_scheme == "aes-256-gcm":
            return self._decrypt_aes_gcm(ciphertext)
        elif self.encryption_scheme == "rc4":
            return self._decrypt_rc4(ciphertext)
        elif self.encryption_scheme == "xor":
            return self._decrypt_xor(ciphertext)
        elif self.encryption_scheme == "custom-base64":
            return self._decrypt_custom_base64(ciphertext)
        else:
            raise ValueError(f"不支持的加密方案: {self.encryption_scheme}")
    
    def _encrypt_aes_gcm(self, plaintext):
        """AES-GCM加密"""
        # 生成12字节nonce (Go的gcm.NonceSize()返回12)
        nonce = secrets.token_bytes(12)
        cipher = AES.new(self.aes_key, AES.MODE_GCM, nonce=nonce)
        ciphertext, tag = cipher.encrypt_and_digest(plaintext.encode())
        # 格式: nonce + ciphertext + tag (与Go的gcm.Seal输出一致)
        encrypted = nonce + ciphertext + tag
        return base64.b64encode(encrypted).decode()
    
    def _decrypt_aes_gcm(self, ciphertext):
        """AES-GCM解密"""
        data = base64.b64decode(ciphertext)
        nonce = data[:12]
        tag = data[-16:]
        ciphertext = data[12:-16]
        cipher = AES.new(self.aes_key, AES.MODE_GCM, nonce=nonce)
        return cipher.decrypt_and_verify(ciphertext, tag).decode()
    
    def _encrypt_rc4(self, plaintext):
        """RC4加密"""
        cipher = ARC4.new(self.aes_key)
        ciphertext = cipher.encrypt(plaintext.encode())
        return base64.b64encode(ciphertext).decode()
    
    def _decrypt_rc4(self, ciphertext):
        """RC4解密"""
        data = base64.b64decode(ciphertext)
        cipher = ARC4.new(self.aes_key)
        plaintext = cipher.decrypt(data)
        return plaintext.decode()
    
    def _encrypt_xor(self, plaintext):
        """XOR加密"""
        plaintext_bytes = plaintext.encode()
        ciphertext = bytearray()
        key_len = len(self.aes_key)
        for i, byte in enumerate(plaintext_bytes):
            ciphertext.append(byte ^ self.aes_key[i % key_len])
        return base64.b64encode(bytes(ciphertext)).decode()
    
    def _decrypt_xor(self, ciphertext):
        """XOR解密"""
        data = base64.b64decode(ciphertext)
        plaintext = bytearray()
        key_len = len(self.aes_key)
        for i, byte in enumerate(data):
            plaintext.append(byte ^ self.aes_key[i % key_len])
        return bytes(plaintext).decode()
    
    def _encrypt_custom_base64(self, plaintext):
        """自定义Base64加密"""
        # 创建自定义字符表
        custom_alphabet = self.aes_key.decode('utf-8')
        standard_alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
        
        # 使用标准Base64编码
        standard_encoded = base64.b64encode(plaintext.encode()).decode()
        
        # 转换为自定义字符表
        trans_table = str.maketrans(standard_alphabet, custom_alphabet)
        custom_encoded = standard_encoded.translate(trans_table)
        
        # 外层再用标准Base64包裹
        return base64.b64encode(custom_encoded.encode()).decode()
    
    def _decrypt_custom_base64(self, ciphertext):
        """自定义Base64解密"""
        # 先解开外层标准Base64
        custom_encoded = base64.b64decode(ciphertext).decode()
        
        # 创建自定义字符表
        custom_alphabet = self.aes_key.decode('utf-8')
        standard_alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
        
        # 转换回标准字符表
        trans_table = str.maketrans(custom_alphabet, standard_alphabet)
        standard_encoded = custom_encoded.translate(trans_table)
        
        # 使用标准Base64解码
        return base64.b64decode(standard_encoded).decode()
    
    def generate_nonce(self):
        """生成随机nonce"""
        return secrets.token_urlsafe(24)
    
    def make_encrypted_request(self, endpoint, data, method="POST"):
        """发送加密请求并验证响应Nonce"""
        # 生成并记住请求nonce
        request_nonce = self.generate_nonce()
        request_timestamp = int(time.time())
        
        # 包装内层数据，嵌入nonce和timestamp
        internal_data = {
            "nonce": request_nonce,
            "timestamp": request_timestamp,
            "data": data
        }
        
        json_data = json.dumps(internal_data)
        encrypted_data = self.encrypt(json_data)
        
        req_body = {
            "timestamp": request_timestamp,
            "nonce": request_nonce,
            "data": encrypted_data
        }
        
        url = f"{self.server_url}{endpoint}"
        if method == "POST":
            response = self.session.post(url, json=req_body)
        else:
            response = self.session.get(url, json=req_body)
        
        resp_json = response.json()
        
        # 验证响应nonce
        if resp_json.get("nonce") != request_nonce:
            raise ValueError("响应Nonce不匹配，可能遭受重放攻击！")
        
        # 解密响应数据
        decrypted = self.decrypt(resp_json["data"])
        internal_response = json.loads(decrypted)
        
        # 验证服务器时间戳
        server_timestamp = internal_response.get("timestamp", 0)
        current_time = int(time.time())
        time_diff = abs(current_time - server_timestamp)
        if time_diff > 300:
            raise ValueError(f"响应时间戳异常，可能遭受离线攻击！时间差: {time_diff}秒")
        
        # 提取实际业务数据
        result = internal_response.get("data", {})
        
        # 返回解密后的结果、请求体和原始加密响应
        return result, req_body, resp_json
    
    def login(self, card_key, hwid="", ip=""):
        """登录"""
        login_data = {
            "project_uuid": self.project_uuid,
            "card_key": card_key,
        }
        if hwid:
            login_data["hwid"] = hwid
        if ip:
            login_data["ip"] = ip
        
        result, request, encrypted_response = self.make_encrypted_request("/api/auth/login", login_data)
        
        if result.get("code") == 0:
            self.token = result["data"]["token"]
            self.session.headers.update({"Authorization": f"Bearer {self.token}"})
        
        return result, request, encrypted_response
    
    def heartbeat(self):
        """心跳"""
        result, _, _ = self.make_encrypted_request("/api/heartbeat", {})
        return result
    
    def get_cloud_var(self, key):
        """获取云变量"""
        result, _, _ = self.make_encrypted_request(f"/api/cloud-var/{key}", {}, method="GET")
        return result
    
    def update_custom_data(self, custom_data):
        """更新专属信息 - 支持任意字符串"""
        # custom_data 可以是任意字符串，不限于 JSON
        data = {"custom_data": custom_data}
        result, _, _ = self.make_encrypted_request("/api/card/custom-data", data)
        return result
    
    def get_project_info(self):
        """获取项目信息"""
        result, _, _ = self.make_encrypted_request("/api/project/info", {}, method="GET")
        return result
    
    def unbind_hwid(self, card_key, hwid):
        """解绑HWID"""
        data = {
            "project_uuid": self.project_uuid,
            "card_key": card_key,
            "hwid": hwid
        }
        result, _, _ = self.make_encrypted_request("/api/card/unbind", data)
        return result


class NextKeyGUI:
    """GUI主窗口"""
    
    def __init__(self, root):
        self.root = root
        self.root.title("NextKey 测试客户端 v1.0")
        self.root.geometry("1000x750")
        
        self.client = None
        self.heartbeat_thread = None
        self.heartbeat_running = False
        self.config_file = "nextkey_client_config.json"
        
        self.setup_ui()
        self.load_config()
    
    def setup_ui(self):
        """设置UI"""
        # 创建笔记本标签页
        notebook = ttk.Notebook(self.root)
        notebook.pack(fill=tk.BOTH, expand=True, padx=5, pady=5)
        
        # 配置页
        config_frame = ttk.Frame(notebook)
        notebook.add(config_frame, text="配置")
        self.setup_config_tab(config_frame)
        
        # 登录页
        login_frame = ttk.Frame(notebook)
        notebook.add(login_frame, text="登录测试")
        self.setup_login_tab(login_frame)
        
        # API测试页
        api_frame = ttk.Frame(notebook)
        notebook.add(api_frame, text="API测试")
        self.setup_api_tab(api_frame)
        
        # 日志页
        log_frame = ttk.Frame(notebook)
        notebook.add(log_frame, text="日志")
        self.setup_log_tab(log_frame)
    
    def setup_config_tab(self, parent):
        """配置标签页"""
        frame = ttk.LabelFrame(parent, text="服务器配置", padding=10)
        frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=10)
        
        # 服务器URL
        ttk.Label(frame, text="服务器URL:").grid(row=0, column=0, sticky=tk.W, pady=5)
        self.server_url_var = tk.StringVar(value="http://localhost:8080")
        ttk.Entry(frame, textvariable=self.server_url_var, width=50).grid(row=0, column=1, pady=5, padx=5)
        
        # 项目UUID
        ttk.Label(frame, text="项目UUID:").grid(row=1, column=0, sticky=tk.W, pady=5)
        self.project_uuid_var = tk.StringVar()
        ttk.Entry(frame, textvariable=self.project_uuid_var, width=50).grid(row=1, column=1, pady=5, padx=5)
        
        # 加密方案
        ttk.Label(frame, text="加密方案:").grid(row=2, column=0, sticky=tk.W, pady=5)
        self.encryption_scheme_var = tk.StringVar(value="aes-256-gcm")
        scheme_combo = ttk.Combobox(frame, textvariable=self.encryption_scheme_var, width=47, state="readonly")
        scheme_combo['values'] = (
            'aes-256-gcm (推荐-安全)', 
            'rc4 (已弃用-不安全)', 
            'xor (已弃用-不安全)', 
            'custom-base64 (不安全)'
        )
        scheme_combo.grid(row=2, column=1, pady=5, padx=5)
        scheme_combo.bind('<<ComboboxSelected>>', self.on_scheme_changed)
        
        # AES密钥
        ttk.Label(frame, text="加密密钥:").grid(row=3, column=0, sticky=tk.W, pady=5)
        self.aes_key_var = tk.StringVar()
        self.key_entry = ttk.Entry(frame, textvariable=self.aes_key_var, width=50, show="*")
        self.key_entry.grid(row=3, column=1, pady=5, padx=5)
        
        # 按钮框
        btn_frame = ttk.Frame(frame)
        btn_frame.grid(row=4, column=0, columnspan=2, pady=10)
        
        ttk.Button(btn_frame, text="从config.yaml读取", command=self.load_from_yaml).pack(side=tk.LEFT, padx=5)
        ttk.Button(btn_frame, text="保存配置", command=self.save_config).pack(side=tk.LEFT, padx=5)
        ttk.Button(btn_frame, text="测试连接", command=self.test_connection).pack(side=tk.LEFT, padx=5)
        
        # 显示/隐藏密钥
        self.show_key_var = tk.BooleanVar()
        ttk.Checkbutton(frame, text="显示密钥", variable=self.show_key_var, 
                       command=self.toggle_key_visibility).grid(row=3, column=2, padx=5)
        
        # 状态信息
        status_frame = ttk.LabelFrame(parent, text="状态信息", padding=10)
        status_frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=10)
        
        self.status_text = scrolledtext.ScrolledText(status_frame, height=10, wrap=tk.WORD)
        self.status_text.pack(fill=tk.BOTH, expand=True)
    
    def setup_login_tab(self, parent):
        """登录测试标签页"""
        frame = ttk.LabelFrame(parent, text="登录信息", padding=10)
        frame.pack(fill=tk.X, padx=10, pady=10)
        
        # 卡密
        ttk.Label(frame, text="卡密:").grid(row=0, column=0, sticky=tk.W, pady=5)
        self.card_key_var = tk.StringVar()
        ttk.Entry(frame, textvariable=self.card_key_var, width=40).grid(row=0, column=1, pady=5, padx=5)
        
        # HWID
        ttk.Label(frame, text="设备码 (可选):").grid(row=1, column=0, sticky=tk.W, pady=5)
        self.hwid_var = tk.StringVar()
        ttk.Entry(frame, textvariable=self.hwid_var, width=40).grid(row=1, column=1, pady=5, padx=5)
        
        # IP
        ttk.Label(frame, text="IP地址 (可选):").grid(row=2, column=0, sticky=tk.W, pady=5)
        self.ip_var = tk.StringVar()
        ttk.Entry(frame, textvariable=self.ip_var, width=40).grid(row=2, column=1, pady=5, padx=5)
        
        # 登录按钮
        ttk.Button(frame, text="登录", command=self.do_login, width=20).grid(row=3, column=0, columnspan=2, pady=10)
        
        # Token显示
        token_frame = ttk.LabelFrame(parent, text="Token信息", padding=10)
        token_frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=10)
        
        self.token_text = scrolledtext.ScrolledText(token_frame, height=15, wrap=tk.WORD)
        self.token_text.pack(fill=tk.BOTH, expand=True)
    
    def setup_api_tab(self, parent):
        """API测试标签页"""
        # 心跳测试
        heartbeat_frame = ttk.LabelFrame(parent, text="心跳测试", padding=10)
        heartbeat_frame.pack(fill=tk.X, padx=10, pady=5)
        
        btn_frame = ttk.Frame(heartbeat_frame)
        btn_frame.pack()
        
        ttk.Button(btn_frame, text="手动心跳", command=self.do_heartbeat).pack(side=tk.LEFT, padx=5)
        ttk.Button(btn_frame, text="开始自动心跳", command=self.start_auto_heartbeat).pack(side=tk.LEFT, padx=5)
        ttk.Button(btn_frame, text="停止自动心跳", command=self.stop_auto_heartbeat).pack(side=tk.LEFT, padx=5)
        
        # 云变量
        cloudvar_frame = ttk.LabelFrame(parent, text="云变量查询", padding=10)
        cloudvar_frame.pack(fill=tk.X, padx=10, pady=5)
        
        ttk.Label(cloudvar_frame, text="变量Key:").pack(side=tk.LEFT, padx=5)
        self.cloudvar_key_var = tk.StringVar()
        ttk.Entry(cloudvar_frame, textvariable=self.cloudvar_key_var, width=30).pack(side=tk.LEFT, padx=5)
        ttk.Button(cloudvar_frame, text="查询", command=self.get_cloudvar).pack(side=tk.LEFT, padx=5)
        
        # 专属信息
        custom_frame = ttk.LabelFrame(parent, text="专属信息更新", padding=10)
        custom_frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=5)
        
        ttk.Label(custom_frame, text="专属数据（支持任意文本）:").pack(anchor=tk.W)
        self.custom_data_text = scrolledtext.ScrolledText(custom_frame, height=5, wrap=tk.WORD)
        self.custom_data_text.pack(fill=tk.BOTH, expand=True, pady=5)
        self.custom_data_text.insert(1.0, '{"user_level": 1, "points": 0}')
        
        ttk.Button(custom_frame, text="更新", command=self.update_custom).pack(pady=5)
        
        # 项目信息
        project_frame = ttk.LabelFrame(parent, text="项目信息", padding=10)
        project_frame.pack(fill=tk.X, padx=10, pady=5)
        
        ttk.Button(project_frame, text="获取项目信息", command=self.get_project_info).pack()
        
        # 解绑HWID
        unbind_frame = ttk.LabelFrame(parent, text="解绑HWID测试", padding=10)
        unbind_frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=5)
        
        # 卡密输入
        row1 = ttk.Frame(unbind_frame)
        row1.pack(fill=tk.X, pady=5)
        ttk.Label(row1, text="卡密:", width=12).pack(side=tk.LEFT, padx=5)
        self.unbind_card_key_var = tk.StringVar()
        ttk.Entry(row1, textvariable=self.unbind_card_key_var, width=30).pack(side=tk.LEFT, padx=5)
        ttk.Button(row1, text="使用登录卡密", command=self.use_login_card_key).pack(side=tk.LEFT, padx=5)
        
        # 当前HWID显示
        row2 = ttk.Frame(unbind_frame)
        row2.pack(fill=tk.X, pady=5)
        ttk.Label(row2, text="当前HWID:", width=12).pack(side=tk.LEFT, padx=5)
        self.current_hwid_display = tk.StringVar(value="未设置")
        ttk.Label(row2, textvariable=self.current_hwid_display, foreground="blue", width=30).pack(side=tk.LEFT, padx=5)
        
        # HWID输入
        row3 = ttk.Frame(unbind_frame)
        row3.pack(fill=tk.X, pady=5)
        ttk.Label(row3, text="HWID:", width=12).pack(side=tk.LEFT, padx=5)
        self.unbind_hwid_var = tk.StringVar()
        ttk.Entry(row3, textvariable=self.unbind_hwid_var, width=30).pack(side=tk.LEFT, padx=5)
        ttk.Button(row3, text="使用登录HWID", command=self.use_login_hwid).pack(side=tk.LEFT, padx=5)
        
        # 快捷操作按钮
        row4 = ttk.Frame(unbind_frame)
        row4.pack(fill=tk.X, pady=10)
        ttk.Button(row4, text="更改HWID", command=self.change_hwid, width=15).pack(side=tk.LEFT, padx=5)
        ttk.Button(row4, text="恢复原HWID", command=self.restore_hwid, width=15).pack(side=tk.LEFT, padx=5)
        ttk.Button(row4, text="解绑HWID", command=self.do_unbind, width=15).pack(side=tk.LEFT, padx=5)
        
        # 解绑历史记录
        ttk.Label(unbind_frame, text="解绑操作历史:").pack(anchor=tk.W, pady=(10,5))
        self.unbind_history_text = scrolledtext.ScrolledText(unbind_frame, height=6, wrap=tk.WORD)
        self.unbind_history_text.pack(fill=tk.BOTH, expand=True)
        
        # 保存原始HWID
        self.original_hwid = ""
    
    def setup_log_tab(self, parent):
        """日志标签页"""
        frame = ttk.Frame(parent)
        frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=10)
        
        # 工具栏
        toolbar = ttk.Frame(frame)
        toolbar.pack(fill=tk.X, pady=5)
        
        ttk.Button(toolbar, text="清空日志", command=self.clear_log).pack(side=tk.LEFT, padx=5)
        ttk.Button(toolbar, text="导出日志", command=self.export_log).pack(side=tk.LEFT, padx=5)
        
        # 日志文本
        self.log_text = scrolledtext.ScrolledText(frame, wrap=tk.WORD)
        self.log_text.pack(fill=tk.BOTH, expand=True)
        
        # 配置日志颜色标签
        self.log_text.tag_config("success", foreground="green")
        self.log_text.tag_config("error", foreground="red")
        self.log_text.tag_config("info", foreground="blue")
    
    def log(self, message, level="info"):
        """记录日志"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        log_msg = f"[{timestamp}] {message}\n"
        
        self.log_text.insert(tk.END, log_msg, level)
        self.log_text.see(tk.END)
        self.log_text.update()
    
    def clear_log(self):
        """清空日志"""
        self.log_text.delete(1.0, tk.END)
    
    def export_log(self):
        """导出日志"""
        filename = filedialog.asksaveasfilename(
            defaultextension=".txt",
            filetypes=[("Text files", "*.txt"), ("All files", "*.*")]
        )
        if filename:
            with open(filename, 'w', encoding='utf-8') as f:
                f.write(self.log_text.get(1.0, tk.END))
            messagebox.showinfo("成功", f"日志已导出到: {filename}")
    
    def on_scheme_changed(self, event=None):
        """加密方案变更时的回调"""
        scheme_display = self.encryption_scheme_var.get()
        scheme = scheme_display.split(' ')[0]  # 提取实际方案名
        
        # 显示安全警告
        if scheme in ['rc4', 'xor', 'custom-base64']:
            messagebox.showwarning(
                "安全警告", 
                f"警告：{scheme} 是不安全的加密方案！\n\n"
                "仅用于测试和兼容性需求。\n"
                "生产环境请使用 aes-256-gcm。"
            )
        
        # 更新密钥格式提示
        self.update_key_hint(scheme)
    
    def update_key_hint(self, scheme):
        """更新密钥格式提示"""
        hints = {
            'aes-256-gcm': "32字节密钥 (64字符hex或base64)",
            'rc4': "hex编码的密钥或任意字符串",
            'xor': "hex编码的密钥或任意字符串",
            'custom-base64': "64个不重复字符的映射表"
        }
        hint = hints.get(scheme, "")
        self.status_text.delete(1.0, tk.END)
        self.status_text.insert(tk.END, f"当前加密方案: {scheme}\n密钥格式: {hint}\n")
    
    def toggle_key_visibility(self):
        """切换密钥显示"""
        if self.show_key_var.get():
            self.key_entry.config(show="")
        else:
            self.key_entry.config(show="*")
    
    def load_from_yaml(self):
        """从config.yaml读取配置（已弃用）"""
        messagebox.showinfo(
            "提示", 
            "加密密钥已改为项目级配置！\n\n"
            "请按以下步骤获取密钥：\n"
            "1. 登录管理后台\n"
            "2. 进入项目管理页面\n"
            "3. 查看项目详情，复制加密密钥\n"
            "4. 将密钥粘贴到AES密钥输入框"
        )
        self.log("提示: 请从管理后台获取项目加密密钥", "info")
    
    def save_config(self):
        """保存配置"""
        scheme_display = self.encryption_scheme_var.get()
        scheme = scheme_display.split(' ')[0]  # 提取实际方案名
        
        config = {
            "server_url": self.server_url_var.get(),
            "project_uuid": self.project_uuid_var.get(),
            "aes_key": self.aes_key_var.get(),
            "encryption_scheme": scheme
        }
        
        try:
            with open(self.config_file, 'w', encoding='utf-8') as f:
                json.dump(config, f, indent=2)
            messagebox.showinfo("成功", "配置已保存")
            self.log("配置已保存", "success")
        except Exception as e:
            messagebox.showerror("错误", f"保存配置失败: {e}")
            self.log(f"保存配置失败: {e}", "error")
    
    def load_config(self):
        """加载配置"""
        if os.path.exists(self.config_file):
            try:
                with open(self.config_file, 'r', encoding='utf-8') as f:
                    config = json.load(f)
                
                self.server_url_var.set(config.get("server_url", "http://localhost:8080"))
                self.project_uuid_var.set(config.get("project_uuid", ""))
                self.aes_key_var.set(config.get("aes_key", ""))
                
                # 加载加密方案
                scheme = config.get("encryption_scheme", "aes-256-gcm")
                scheme_map = {
                    'aes-256-gcm': 'aes-256-gcm (推荐-安全)',
                    'rc4': 'rc4 (已弃用-不安全)',
                    'xor': 'xor (已弃用-不安全)',
                    'custom-base64': 'custom-base64 (不安全)'
                }
                self.encryption_scheme_var.set(scheme_map.get(scheme, 'aes-256-gcm (推荐-安全)'))
                
                self.log("配置已加载", "success")
            except Exception as e:
                self.log(f"加载配置失败: {e}", "error")
    
    def test_connection(self):
        """测试连接"""
        try:
            scheme_display = self.encryption_scheme_var.get()
            scheme = scheme_display.split(' ')[0]  # 提取实际方案名
            
            self.client = NextKeyClient(
                self.server_url_var.get(),
                self.project_uuid_var.get(),
                self.aes_key_var.get(),
                scheme
            )
            
            url = f"{self.server_url_var.get()}/api/heartbeat"
            response = requests.get(url, timeout=5)
            
            self.status_text.delete(1.0, tk.END)
            self.status_text.insert(tk.END, "✓ 服务器连接成功\n")
            self.status_text.insert(tk.END, f"✓ 服务器URL: {self.server_url_var.get()}\n")
            self.status_text.insert(tk.END, f"✓ 加密方案: {scheme}\n")
            self.status_text.insert(tk.END, f"✓ 响应状态: {response.status_code}\n")
            self.status_text.insert(tk.END, f"✓ 密钥前8字节(hex): {self.client.aes_key[:8].hex()}\n")
            
            self.log(f"服务器连接测试成功 (方案: {scheme})", "success")
            messagebox.showinfo("成功", "服务器连接正常")
        except Exception as e:
            self.status_text.delete(1.0, tk.END)
            self.status_text.insert(tk.END, f"✗ 连接失败: {e}\n")
            self.log(f"连接测试失败: {e}", "error")
            messagebox.showerror("错误", f"连接失败: {e}")
    
    def do_login(self):
        """执行登录"""
        if not self.client:
            try:
                scheme_display = self.encryption_scheme_var.get()
                scheme = scheme_display.split(' ')[0]  # 提取实际方案名
                
                self.client = NextKeyClient(
                    self.server_url_var.get(),
                    self.project_uuid_var.get(),
                    self.aes_key_var.get(),
                    scheme
                )
            except Exception as e:
                messagebox.showerror("错误", f"初始化客户端失败: {e}")
                self.log(f"初始化客户端失败: {e}", "error")
                return
        
        card_key = self.card_key_var.get()
        if not card_key:
            messagebox.showwarning("警告", "请输入卡密")
            return
        
        try:
            self.log(f"开始登录，卡密: {card_key}", "info")
            result, request, encrypted_response = self.client.login(
                card_key,
                self.hwid_var.get(),
                self.ip_var.get()
            )
            
            # 显示结果
            self.token_text.delete(1.0, tk.END)
            
            # 请求信息
            self.token_text.insert(tk.END, "=== 请求信息 ===\n", "info")
            self.token_text.insert(tk.END, f"Timestamp: {request['timestamp']}\n")
            self.token_text.insert(tk.END, f"Nonce: {request['nonce']}\n")
            self.token_text.insert(tk.END, f"加密数据: {request['data'][:50]}...\n\n")
            
            # 加密响应信息
            self.token_text.insert(tk.END, "=== 加密响应 ===\n", "info")
            self.token_text.insert(tk.END, f"响应Nonce: {encrypted_response.get('nonce', 'N/A')}\n")
            self.token_text.insert(tk.END, f"加密数据: {encrypted_response.get('data', 'N/A')[:50]}...\n\n")
            
            # 解密后的响应信息
            self.token_text.insert(tk.END, "=== 解密后的响应 ===\n", "info")
            self.token_text.insert(tk.END, json.dumps(result, indent=2, ensure_ascii=False))
            
            if result.get("code") == 0:
                self.log("登录成功", "success")
                messagebox.showinfo("成功", "登录成功！")
            else:
                self.log(f"登录失败: {result.get('message')}", "error")
                messagebox.showerror("失败", f"登录失败: {result.get('message')}")
        
        except Exception as e:
            self.log(f"登录异常: {e}", "error")
            messagebox.showerror("错误", f"登录异常: {e}")
    
    def do_heartbeat(self):
        """执行心跳"""
        if not self.client or not self.client.token:
            messagebox.showwarning("警告", "请先登录")
            return
        
        try:
            result = self.client.heartbeat()
            
            if result.get("code") == 0:
                self.log("心跳成功", "success")
            else:
                self.log(f"心跳失败: {result.get('message')}", "error")
                messagebox.showerror("失败", f"心跳失败: {result.get('message')}")
        
        except Exception as e:
            self.log(f"心跳异常: {e}", "error")
            messagebox.showerror("错误", f"心跳异常: {e}")
    
    def start_auto_heartbeat(self):
        """开始自动心跳"""
        if not self.client or not self.client.token:
            messagebox.showwarning("警告", "请先登录")
            return
        
        if self.heartbeat_running:
            messagebox.showinfo("提示", "自动心跳已在运行")
            return
        
        self.heartbeat_running = True
        self.heartbeat_thread = threading.Thread(target=self.auto_heartbeat_loop, daemon=True)
        self.heartbeat_thread.start()
        self.log("已启动自动心跳 (间隔30秒)", "success")
    
    def stop_auto_heartbeat(self):
        """停止自动心跳"""
        self.heartbeat_running = False
        self.log("已停止自动心跳", "info")
    
    def auto_heartbeat_loop(self):
        """自动心跳循环"""
        while self.heartbeat_running:
            time.sleep(30)
            if self.heartbeat_running:
                try:
                    result = self.client.heartbeat()
                    if result.get("code") == 0:
                        self.log("自动心跳成功", "success")
                    else:
                        self.log(f"自动心跳失败: {result.get('message')}", "error")
                except Exception as e:
                    self.log(f"自动心跳异常: {e}", "error")
    
    def get_cloudvar(self):
        """获取云变量"""
        if not self.client or not self.client.token:
            messagebox.showwarning("警告", "请先登录")
            return
        
        key = self.cloudvar_key_var.get()
        if not key:
            messagebox.showwarning("警告", "请输入变量Key")
            return
        
        try:
            result = self.client.get_cloud_var(key)
            
            if result.get("code") == 0:
                value = result['data']['value']
                self.log(f"云变量 [{key}] = {value}", "success")
                messagebox.showinfo("成功", f"变量值: {value}")
            else:
                self.log(f"获取云变量失败: {result.get('message')}", "error")
                messagebox.showerror("失败", f"获取失败: {result.get('message')}")
        
        except Exception as e:
            self.log(f"获取云变量异常: {e}", "error")
            messagebox.showerror("错误", f"获取异常: {e}")
    
    def update_custom(self):
        """更新专属信息"""
        if not self.client or not self.client.token:
            messagebox.showwarning("警告", "请先登录")
            return
        
        try:
            custom_data_str = self.custom_data_text.get(1.0, tk.END).strip()
            
            # 直接使用字符串，不强制要求 JSON 格式
            result = self.client.update_custom_data(custom_data_str)
            
            if result.get("code") == 0:
                self.log(f"专属信息更新成功: {custom_data_str}", "success")
                messagebox.showinfo("成功", "专属信息更新成功")
            else:
                self.log(f"更新专属信息失败: {result.get('message')}", "error")
                messagebox.showerror("失败", f"更新失败: {result.get('message')}")
        
        except Exception as e:
            self.log(f"更新专属信息异常: {e}", "error")
            messagebox.showerror("错误", f"更新异常: {e}")
    
    def get_project_info(self):
        """获取项目信息"""
        if not self.client or not self.client.token:
            messagebox.showwarning("警告", "请先登录")
            return
        
        try:
            result = self.client.get_project_info()
            
            if result.get("code") == 0:
                info = json.dumps(result['data'], indent=2, ensure_ascii=False)
                self.log(f"项目信息: {info}", "success")
                messagebox.showinfo("项目信息", info)
            else:
                self.log(f"获取项目信息失败: {result.get('message')}", "error")
                messagebox.showerror("失败", f"获取失败: {result.get('message')}")
        
        except Exception as e:
            self.log(f"获取项目信息异常: {e}", "error")
            messagebox.showerror("错误", f"获取异常: {e}")
    
    def use_login_card_key(self):
        """使用登录页面的卡密"""
        card_key = self.card_key_var.get()
        if card_key:
            self.unbind_card_key_var.set(card_key)
            self.log(f"已复制登录卡密: {card_key}", "info")
        else:
            messagebox.showwarning("警告", "登录页面未设置卡密")
    
    def use_login_hwid(self):
        """使用登录页面的HWID"""
        hwid = self.hwid_var.get()
        if hwid:
            self.unbind_hwid_var.set(hwid)
            self.current_hwid_display.set(hwid)
            if not self.original_hwid:
                self.original_hwid = hwid
            self.log(f"已复制登录HWID: {hwid}", "info")
        else:
            messagebox.showwarning("警告", "登录页面未设置HWID")
    
    def change_hwid(self):
        """更改HWID（生成新的测试HWID）"""
        current_hwid = self.unbind_hwid_var.get()
        if not current_hwid:
            messagebox.showwarning("警告", "请先设置当前HWID")
            return
        
        # 保存原始HWID
        if not self.original_hwid:
            self.original_hwid = current_hwid
        
        # 生成新HWID
        timestamp = datetime.now().strftime("%H%M%S")
        new_hwid = f"{current_hwid}-changed-{timestamp}"
        
        self.unbind_hwid_var.set(new_hwid)
        self.current_hwid_display.set(new_hwid)
        
        self.log(f"HWID已更改: {current_hwid} -> {new_hwid}", "info")
        self.unbind_history_text.insert(tk.END, f"[{datetime.now().strftime('%H:%M:%S')}] HWID更改: {current_hwid} -> {new_hwid}\n")
        self.unbind_history_text.see(tk.END)
        
        messagebox.showinfo("成功", f"HWID已更改为:\n{new_hwid}\n\n提示：如需测试解绑，请先用新HWID登录，然后解绑旧HWID")
    
    def restore_hwid(self):
        """恢复原始HWID"""
        if not self.original_hwid:
            messagebox.showwarning("警告", "未保存原始HWID")
            return
        
        current = self.unbind_hwid_var.get()
        self.unbind_hwid_var.set(self.original_hwid)
        self.current_hwid_display.set(self.original_hwid)
        
        self.log(f"HWID已恢复: {current} -> {self.original_hwid}", "info")
        self.unbind_history_text.insert(tk.END, f"[{datetime.now().strftime('%H:%M:%S')}] HWID恢复: {current} -> {self.original_hwid}\n")
        self.unbind_history_text.see(tk.END)
        
        messagebox.showinfo("成功", f"HWID已恢复为:\n{self.original_hwid}")
    
    def do_unbind(self):
        """执行解绑"""
        if not self.client or not self.client.token:
            messagebox.showwarning("警告", "请先登录")
            return
        
        card_key = self.unbind_card_key_var.get()
        hwid = self.unbind_hwid_var.get()
        
        if not card_key:
            messagebox.showwarning("警告", "请输入卡密")
            return
        
        if not hwid:
            messagebox.showwarning("警告", "请输入要解绑的HWID")
            return
        
        # 确认对话框
        confirm = messagebox.askyesno(
            "确认解绑", 
            f"确定要解绑以下HWID吗？\n\n卡密: {card_key}\nHWID: {hwid}\n\n解绑后该HWID将从卡密绑定列表中移除。"
        )
        
        if not confirm:
            return
        
        try:
            self.log(f"开始解绑，卡密: {card_key}, HWID: {hwid}", "info")
            result = self.client.unbind_hwid(card_key, hwid)
            
            timestamp = datetime.now().strftime('%H:%M:%S')
            
            if result.get("code") == 0:
                self.log(f"解绑成功: {hwid}", "success")
                self.unbind_history_text.insert(tk.END, f"[{timestamp}] ✓ 解绑成功: {hwid}\n", "success")
                self.unbind_history_text.see(tk.END)
                messagebox.showinfo("成功", f"HWID解绑成功！\n\n已解绑: {hwid}\n\n提示：现在可以用这个HWID重新登录")
            else:
                error_msg = result.get('message', '未知错误')
                self.log(f"解绑失败: {error_msg}", "error")
                self.unbind_history_text.insert(tk.END, f"[{timestamp}] ✗ 解绑失败: {error_msg}\n", "error")
                self.unbind_history_text.see(tk.END)
                messagebox.showerror("失败", f"解绑失败:\n{error_msg}")
        
        except Exception as e:
            timestamp = datetime.now().strftime('%H:%M:%S')
            self.log(f"解绑异常: {e}", "error")
            self.unbind_history_text.insert(tk.END, f"[{timestamp}] ✗ 解绑异常: {e}\n", "error")
            self.unbind_history_text.see(tk.END)
            messagebox.showerror("错误", f"解绑异常:\n{e}")
        
        # 配置历史记录文本的颜色标签
        self.unbind_history_text.tag_config("success", foreground="green")
        self.unbind_history_text.tag_config("error", foreground="red")


def main():
    root = tk.Tk()
    app = NextKeyGUI(root)
    root.mainloop()


if __name__ == "__main__":
    main()

