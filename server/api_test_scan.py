import requests
import time
import json

# 配置
BASE_URL = "http://localhost:12345/api"
ADMIN_USERNAME = "admin"
ADMIN_PASSWORD = "123456"

# 测试用户配置
TEST_USER_COUNT = 10
TEST_USER_PREFIX = "testuser"
TEST_USER_PASSWORD = "password123"
TEST_USER_EMAIL_TEMPLATE = "test{}@example.com"

def log(msg, type="INFO"):
    print(f"[{type}] {msg}")

def create_test_users():
    log("开始创建测试用户...")
    for i in range(1, TEST_USER_COUNT + 1):
        username = f"{TEST_USER_PREFIX}{i:03d}" # testuser001 ... testuser010
        email = TEST_USER_EMAIL_TEMPLATE.format(i)
        
        try:
            reg_payload = {
                "username": username,
                "password": TEST_USER_PASSWORD,
                "email": email
            }
            res = requests.post(f"{BASE_URL}/auth/register", json=reg_payload)
            if res.status_code == 200 and res.json().get("code") == 1000:
                log(f"用户注册成功: {username}")
            else:
                log(f"用户注册跳过 (可能已存在): {username} - {res.json().get('message')}", "WARN")
        except Exception as e:
            log(f"注册请求异常 {username}: {e}", "ERROR")
    log("测试用户创建流程结束。")

def run_task():
    # 0. 创建测试用户
    create_test_users()

    log("开始执行管理员扫描任务...")

    # 1. 管理员登录获取 Token
    token = ""
    try:
        login_payload = {
            "username": ADMIN_USERNAME,
            "password": ADMIN_PASSWORD
        }
        res = requests.post(f"{BASE_URL}/auth/login", json=login_payload)
        data = res.json()
        if data.get("code") == 1000:
            token = data["data"]["token"]
            log(f"管理员登录成功, 获取到 Token")
        else:
            log(f"管理员登录失败: {data.get('message')}", "ERROR")
            return
    except Exception as e:
        log(f"登录请求异常: {e}", "ERROR")
        return

    # 2. 扫描音乐
    if token:
        log(f"开始扫描音乐...")
        try:
            headers = {"Authorization": f"Bearer {token}"}
            
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
    
    log("任务执行结束。")

if __name__ == "__main__":
    run_task()
