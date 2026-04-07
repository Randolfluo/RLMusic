import requests
import argparse


def _log(logger, msg, level="INFO"):
    if logger:
        logger(msg, level)
    else:
        print(f"[{level}] {msg}")


def admin_login(base_url, admin_username, admin_password, timeout=10, logger=None):
    try:
        login_payload = {"username": admin_username, "password": admin_password}
        res = requests.post(f"{base_url}/auth/login", json=login_payload, timeout=timeout)
        data = res.json()
        if data.get("code") == 1000:
            token = data["data"]["token"]
            _log(logger, "管理员登录成功，已获取 Token")
            return token
        _log(logger, f"管理员登录失败: {data.get('message')}", "ERROR")
    except Exception as e:
        _log(logger, f"管理员登录请求异常: {e}", "ERROR")
    return ""


def run_scan(base_url, token, timeout=180, logger=None):
    if not token:
        _log(logger, "无有效 Token，跳过扫描。", "ERROR")
        return False

    _log(logger, "开始执行管理员扫描任务...")
    try:
        headers = {"Authorization": f"Bearer {token}"}
        res = requests.post(f"{base_url}/song/scan", headers=headers, timeout=timeout)
        if res.status_code == 200:
            scan_data = res.json().get("data", {})
            _log(
                logger,
                f"扫描完成! 新增: {scan_data.get('added', 0)}, 更新: {scan_data.get('updated', 0)}",
            )
            return True
        _log(logger, f"扫描请求失败: {res.status_code} - {res.text}", "ERROR")
    except requests.exceptions.Timeout:
        _log(logger, "扫描请求超时(后台可能仍在运行)", "WARN")
    except Exception as e:
        _log(logger, f"扫描过程中发生错误: {e}", "ERROR")
    return False


def main():
    parser = argparse.ArgumentParser(description="管理员登录并执行歌曲扫描")
    parser.add_argument("--base-url", default="http://localhost:12345/api", help="后端 API 基础地址")
    parser.add_argument("--username", default="admin", help="管理员用户名")
    parser.add_argument("--password", default="123456", help="管理员密码")
    parser.add_argument("--login-timeout", type=int, default=10, help="登录请求超时时间(秒)")
    parser.add_argument("--scan-timeout", type=int, default=180, help="扫描请求超时时间(秒)")
    args = parser.parse_args()

    token = admin_login(
        base_url=args.base_url,
        admin_username=args.username,
        admin_password=args.password,
        timeout=args.login_timeout,
        logger=_log_print,
    )
    ok = run_scan(
        base_url=args.base_url,
        token=token,
        timeout=args.scan_timeout,
        logger=_log_print,
    )
    if ok:
        _log_print("扫描任务执行完成。")
    else:
        _log_print("扫描任务执行失败。", "ERROR")


def _log_print(msg, level="INFO"):
    print(f"[{level}] {msg}")


if __name__ == "__main__":
    main()
