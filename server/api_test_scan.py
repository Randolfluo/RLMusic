import requests
import time
import json

# 配置
BASE_URL = "http://localhost:3000/api"
USER_COUNT = 5
BASE_USERNAME = "testuser"
BASE_PASSWORD = "password123"
BASE_EMAIL = "test{}@example.com"

def log(msg, type="INFO"):
    print(f"[{type}] {msg}")

def run_task():
    log("开始执行批量用户创建与扫描任务...")

    for i in range(1, USER_COUNT + 1):
        username = f"{BASE_USERNAME}{i:03d}" # testuser001 ... testuser005
        email = BASE_EMAIL.format(i)
        
        log(f"--- 处理用户: {username} ---")

        # 1. 注册 (如果已存在则忽略错误)
        try:
            reg_payload = {
                "username": username,
                "password": BASE_PASSWORD,
                "email": email
            }
            res = requests.post(f"{BASE_URL}/auth/register", json=reg_payload)
            if res.status_code == 200 and res.json().get("code") == 1000:
                log(f"注册成功: {username}")
            else:
                log(f"注册跳过 (可能已存在): {res.json().get('message')}", "WARN")
        except Exception as e:
            log(f"注册请求异常: {e}", "ERROR")
            continue

        # 2. 登录获取 Token
        token = ""
        try:
            login_payload = {
                "username": username,
                "password": BASE_PASSWORD
            }
            res = requests.post(f"{BASE_URL}/auth/login", json=login_payload)
            data = res.json()
            if data.get("code") == 1000:
                token = data["data"]["token"]
                log(f"登录成功, 获取到 Token")
            else:
                log(f"登录失败: {data.get('message')}", "ERROR")
                continue
        except Exception as e:
            log(f"登录请求异常: {e}", "ERROR")
            continue

        # 3. 扫描音乐
        if token:
            log(f"开始扫描音乐...")
            try:
                headers = {"Authorization": f"Bearer {token}"}
                # 注意：这里扫描的是服务器配置的默认 music 文件夹
                # 如果代码逻辑是扫描 *用户目录*，则不同用户会有不同结果
                # 如果代码逻辑是扫描 *系统公共目录* 并关联到用户，则所有用户都会有一份数据
                
                # 这是一个长连接请求，可能需要较长时间
                res = requests.post(f"{BASE_URL}/song/scan", headers=headers, timeout=120) 
                
                if res.status_code == 200:
                    scan_data = res.json().get("data", {})
                    log(f"扫描完成! 新增: {scan_data.get('added', 0)}, 更新: {scan_data.get('updated', 0)}")
                else:
                    log(f"扫描请求失败: {res.status_code} - {res.text}", "ERROR")

            except requests.exceptions.Timeout:
                log("扫描请求超时 (后台可能仍在运行)", "WARN")
            except Exception as e:
                log(f"扫描过程中发生错误: {e}", "ERROR")
        
        log(f"用户 {username} 处理完毕\n")
        time.sleep(1) # 稍微间隔一下，避免请求太快

    log("所有任务执行结束。")

if __name__ == "__main__":
    run_task()
