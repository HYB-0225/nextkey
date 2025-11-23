use clap::{Parser, Subcommand};
use nextkey_sdk::NextKeyClient;
use std::process;

#[derive(Parser)]
#[command(name = "nextkey-cli")]
#[command(about = "NextKey SDK Command Line Tool", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// 登录获取Token
    Login {
        /// 服务器URL
        #[arg(short, long)]
        server: String,

        /// 项目UUID
        #[arg(short = 'u', long)]
        uuid: String,

        /// AES密钥
        #[arg(short, long)]
        key: String,

        /// 卡密
        #[arg(short, long)]
        cardkey: String,

        /// 设备码（可选）
        #[arg(long, default_value = "")]
        hwid: String,

        /// IP地址（可选）
        #[arg(long, default_value = "")]
        ip: String,
    },

    /// 心跳验证
    Heartbeat {
        /// 服务器URL
        #[arg(short, long)]
        server: String,

        /// 项目UUID
        #[arg(short = 'u', long)]
        uuid: String,

        /// AES密钥
        #[arg(short, long)]
        key: String,

        /// Token
        #[arg(short, long)]
        token: String,
    },

    /// 获取云变量
    GetCloudVar {
        /// 服务器URL
        #[arg(short, long)]
        server: String,

        /// 项目UUID
        #[arg(short = 'u', long)]
        uuid: String,

        /// AES密钥
        #[arg(short, long)]
        key: String,

        /// Token
        #[arg(short, long)]
        token: String,

        /// 变量名
        #[arg(short = 'n', long)]
        varkey: String,
    },

    /// 更新专属信息
    UpdateCustomData {
        /// 服务器URL
        #[arg(short, long)]
        server: String,

        /// 项目UUID
        #[arg(short = 'u', long)]
        uuid: String,

        /// AES密钥
        #[arg(short, long)]
        key: String,

        /// Token
        #[arg(short, long)]
        token: String,

        /// 专属数据
        #[arg(short, long)]
        data: String,
    },

    /// 获取项目信息
    GetProjectInfo {
        /// 服务器URL
        #[arg(short, long)]
        server: String,

        /// 项目UUID
        #[arg(short = 'u', long)]
        uuid: String,

        /// AES密钥
        #[arg(short, long)]
        key: String,

        /// Token
        #[arg(short, long)]
        token: String,
    },
}

fn main() {
    let cli = Cli::parse();

    match cli.command {
        Commands::Login {
            server,
            uuid,
            key,
            cardkey,
            hwid,
            ip,
        } => {
            let mut client = match NextKeyClient::new(&server, &uuid, &key) {
                Ok(c) => c,
                Err(e) => {
                    eprintln!("创建客户端失败: {}", e);
                    process::exit(1);
                }
            };

            match client.login(&cardkey, &hwid, &ip) {
                Ok(response) => {
                    if response.code == 0 {
                        if let Some(data) = response.data {
                            println!("登录成功!");
                            println!("Token: {}", data.token);
                            println!("Token 过期时间: {}", data.expire_at);

                            if let Some(card) = data.card {
                                println!("卡密信息:");
                                println!("  ID: {}", card.id);
                                println!("  ProjectID: {}", card.project_id);
                                println!("  卡密: {}", card.card_key);
                                println!("  已激活: {}", card.activated);
                                if let Some(activated_at) = card.activated_at {
                                    println!("  激活时间: {}", activated_at);
                                }
                                println!("  已冻结: {}", card.frozen);
                                println!("  时长(秒): {}", card.duration);
                                if let Some(expire_at) = card.expire_at {
                                    println!("  卡密到期时间: {}", expire_at);
                                }
                                println!("  类型: {}", card.card_type);
                                if !card.note.is_empty() {
                                    println!("  备注: {}", card.note);
                                }
                                if !card.custom_data.is_empty() {
                                    println!("  专属信息: {}", card.custom_data);
                                }
                                if !card.hwid_list.is_empty() {
                                    println!("  HWID列表: {:?}", card.hwid_list);
                                }
                                if !card.ip_list.is_empty() {
                                    println!("  IP列表: {:?}", card.ip_list);
                                }
                                println!("  MaxHWID: {}", card.max_hwid);
                                println!("  MaxIP: {}", card.max_ip);
                                println!("  创建时间: {}", card.created_at);
                                println!("  更新时间: {}", card.updated_at);
                            } else {
                                println!("卡密信息: (free 模式，无卡密)");
                            }
                        }
                    } else {
                        eprintln!("登录失败: {}", response.message);
                        process::exit(1);
                    }
                }
                Err(e) => {
                    eprintln!("登录异常: {}", e);
                    process::exit(1);
                }
            }
        }

        Commands::Heartbeat {
            server,
            uuid,
            key,
            token,
        } => {
            let mut client = match NextKeyClient::new(&server, &uuid, &key) {
                Ok(c) => c,
                Err(e) => {
                    eprintln!("创建客户端失败: {}", e);
                    process::exit(1);
                }
            };

            client.set_token(token);

            match client.heartbeat() {
                Ok(response) => {
                    if response.code == 0 {
                        println!("心跳成功");
                    } else {
                        eprintln!("心跳失败: {}", response.message);
                        process::exit(1);
                    }
                }
                Err(e) => {
                    eprintln!("心跳异常: {}", e);
                    process::exit(1);
                }
            }
        }

        Commands::GetCloudVar {
            server,
            uuid,
            key,
            token,
            varkey,
        } => {
            let mut client = match NextKeyClient::new(&server, &uuid, &key) {
                Ok(c) => c,
                Err(e) => {
                    eprintln!("创建客户端失败: {}", e);
                    process::exit(1);
                }
            };

            client.set_token(token);

            match client.get_cloud_var(&varkey) {
                Ok(response) => {
                    if response.code == 0 {
                        if let Some(data) = response.data {
                            println!("变量名: {}", data.key);
                            println!("变量值: {}", data.value);
                        }
                    } else {
                        eprintln!("获取云变量失败: {}", response.message);
                        process::exit(1);
                    }
                }
                Err(e) => {
                    eprintln!("获取云变量异常: {}", e);
                    process::exit(1);
                }
            }
        }

        Commands::UpdateCustomData {
            server,
            uuid,
            key,
            token,
            data,
        } => {
            let mut client = match NextKeyClient::new(&server, &uuid, &key) {
                Ok(c) => c,
                Err(e) => {
                    eprintln!("创建客户端失败: {}", e);
                    process::exit(1);
                }
            };

            client.set_token(token);

            match client.update_custom_data(&data) {
                Ok(response) => {
                    if response.code == 0 {
                        println!("专属信息更新成功");
                    } else {
                        eprintln!("更新失败: {}", response.message);
                        process::exit(1);
                    }
                }
                Err(e) => {
                    eprintln!("更新异常: {}", e);
                    process::exit(1);
                }
            }
        }

        Commands::GetProjectInfo {
            server,
            uuid,
            key,
            token,
        } => {
            let mut client = match NextKeyClient::new(&server, &uuid, &key) {
                Ok(c) => c,
                Err(e) => {
                    eprintln!("创建客户端失败: {}", e);
                    process::exit(1);
                }
            };

            client.set_token(token);

            match client.get_project_info() {
                Ok(response) => {
                    if response.code == 0 {
                        if let Some(data) = response.data {
                            println!("项目UUID: {}", data.uuid);
                            println!("项目名称: {}", data.name);
                            println!("版本: {}", data.version);
                            println!("更新地址: {}", data.update_url);
                        }
                    } else {
                        eprintln!("获取项目信息失败: {}", response.message);
                        process::exit(1);
                    }
                }
                Err(e) => {
                    eprintln!("获取项目信息异常: {}", e);
                    process::exit(1);
                }
            }
        }
    }
}

