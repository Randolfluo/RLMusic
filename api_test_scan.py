import statistics
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

import requests

# -----------------------------
# 基础地址配置
# -----------------------------
# 后端 API 根地址（用于注册、登录、系统信息等接口）
BASE_URL = "http://localhost:12345/api"
# 前端 Web 地址（用于首页可达性/响应性能压测）
WEB_URL = "http://localhost:23456"

# -----------------------------
# 测试用户参数
# -----------------------------
# 批量创建用户数量。建议:
# - 日常快速验证: 10~30
# - 较大规模压测: 100+
TEST_USER_COUNT = 100
# 用户名前缀，最终用户名格式示例: testuser001
TEST_USER_PREFIX = "testuser"
# 压测登录场景使用的统一密码
TEST_USER_PASSWORD = "password123"
# 测试邮箱模板，{} 会被 1..N 自动替换
TEST_USER_EMAIL_TEMPLATE = "test{}@example.com"

# -----------------------------
# 压测参数
# -----------------------------
# 单次 HTTP 请求超时（秒）
REQUEST_TIMEOUT = 10
# 流式播放探测读取上限（字节），避免一次请求下载整首歌
# 说明:
# - 该值越大，越接近真实连续播放，但会增加服务端与网络负载
# - 该值越小，更接近“首段缓冲”探测，适合高并发稳定性测试
STREAM_READ_BYTES = 64 * 1024
# 分阶段升压配置:
# - name: 阶段名
# - concurrency: 并发线程数
# - total_requests: 该阶段总请求数（所有场景合计）
PRESSURE_STAGES = [
    {"name": "stage-1", "concurrency": 20, "total_requests": 300},
    {"name": "stage-2", "concurrency": 50, "total_requests": 1000},
    {"name": "stage-3", "concurrency": 100, "total_requests": 3000},
]


def log(msg, level="INFO"):
    print(f"[{level}] {msg}")


def percentile(values, p):
    """计算百分位数（p 取值 0~1）。"""
    if not values:
        return 0.0
    if len(values) == 1:
        return values[0]
    sorted_values = sorted(values)
    k = (len(sorted_values) - 1) * p
    f = int(k)
    c = min(f + 1, len(sorted_values) - 1)
    if f == c:
        return sorted_values[f]
    return sorted_values[f] + (sorted_values[c] - sorted_values[f]) * (k - f)


def _pick_id(record, keys):
    """从字典中按候选键提取第一个可用 ID。"""
    if not isinstance(record, dict):
        return None
    for key in keys:
        val = record.get(key)
        if val not in (None, ""):
            return val
    return None


def discover_sample_ids():
    """自动发现公开歌单及歌曲样本ID。"""
    playlist_id = None
    song_id = None
    artist_id = None
    album_id = None

    try:
        res = requests.get(
            f"{BASE_URL}/song/playlists/public",
            params={"page": 1, "limit": 1},
            timeout=REQUEST_TIMEOUT,
        )
        if res.status_code == 200:
            body = res.json()
            data = body.get("data", {}) if isinstance(body, dict) else {}
            plist = data.get("list") or data.get("playlists") or []
            if isinstance(plist, list) and plist:
                playlist_id = _pick_id(plist[0], ["id", "ID", "playlist_id"])
    except Exception as e:
        log(f"自动发现 playlist_id 失败: {e}", "WARN")

    if playlist_id is not None:
        try:
            res = requests.get(
                f"{BASE_URL}/song/playlist/public/{playlist_id}",
                params={"page": 1, "limit": 1},
                timeout=REQUEST_TIMEOUT,
            )
            if res.status_code == 200:
                body = res.json()
                data = body.get("data", {}) if isinstance(body, dict) else {}
                songs = data.get("songs") or data.get("list") or []
                if isinstance(songs, list) and songs:
                    first_song = songs[0]
                    song_id = _pick_id(first_song, ["id", "ID", "song_id"])
                    artist_id = _pick_id(first_song, ["artist_id", "artistId", "ArtistID"])
                    album_id = _pick_id(first_song, ["album_id", "albumId", "AlbumID"])
        except Exception as e:
            log(f"自动发现 song_id 失败: {e}", "WARN")

    return {
        "playlist_id": playlist_id,
        "song_id": song_id,
        "artist_id": artist_id,
        "album_id": album_id,
    }


def build_scenarios(test_users):
    """构建通用 API 压测场景（不包含流式播放场景）。"""
    scenarios = [
        {"name": "api_system_stats", "method": "GET", "url": f"{BASE_URL}/system/stats", "auth": False},
        {"name": "api_local_ips", "method": "GET", "url": f"{BASE_URL}/system/local-ips", "params": {"port": "23456"}, "auth": False},
        {"name": "api_public_playlists", "method": "GET", "url": f"{BASE_URL}/song/playlists/public", "auth": False},
        {"name": "api_search_song", "method": "GET", "url": f"{BASE_URL}/search/song", "params": {"keywords": "love", "limit": 20, "offset": 1}, "auth": False},
        {"name": "api_search_artist", "method": "GET", "url": f"{BASE_URL}/search/artist", "params": {"keywords": "jay", "limit": 20, "offset": 1}, "auth": False},
        {"name": "api_search_album", "method": "GET", "url": f"{BASE_URL}/search/album", "params": {"keywords": "greatest", "limit": 20, "offset": 1}, "auth": False},
        {"name": "api_search_playlist", "method": "GET", "url": f"{BASE_URL}/search/playlist", "params": {"keywords": "流行", "limit": 20, "offset": 1}, "auth": False},
        {"name": "api_user_login", "method": "POST", "url": f"{BASE_URL}/auth/login", "auth": False},
        {"name": "web_home", "method": "GET", "url": f"{WEB_URL}/", "auth": False},
    ]

    samples = discover_sample_ids()
    playlist_id = samples["playlist_id"]
    song_id = samples["song_id"]
    artist_id = samples["artist_id"]
    album_id = samples["album_id"]

    if playlist_id is not None:
        scenarios.append(
            {
                "name": "api_playlist_public_detail",
                "method": "GET",
                "url": f"{BASE_URL}/song/playlist/public/{playlist_id}",
                "params": {"page": 1, "limit": 20},
                "auth": False,
            }
        )
    if song_id is not None:
        scenarios.extend(
            [
                {"name": "api_song_detail", "method": "GET", "url": f"{BASE_URL}/song/detail/{song_id}", "auth": False},
                {"name": "api_song_lyric", "method": "GET", "url": f"{BASE_URL}/song/lyric/{song_id}", "auth": False},
                {"name": "api_song_cover", "method": "GET", "url": f"{BASE_URL}/song/cover/{song_id}", "auth": False},
            ]
        )
    else:
        log("未发现可用 song_id，已跳过歌曲详情/歌词/封面场景。", "WARN")

    if artist_id is not None:
        scenarios.append(
            {"name": "api_artist_detail", "method": "GET", "url": f"{BASE_URL}/song/artist/{artist_id}", "params": {"page": 1, "limit": 20}, "auth": False}
        )
    if album_id is not None:
        scenarios.append(
            {"name": "api_album_detail", "method": "GET", "url": f"{BASE_URL}/song/album/{album_id}", "auth": False}
        )

    log(f"本轮压测场景数: {len(scenarios)}")
    return scenarios


def run_song_stream_test(stage):
    """独立执行歌曲流式播放压测。

    测试方式:
    1) 请求 /song/stream/:id，并携带 Range 头只拉取前 N 字节。
    2) 使用 stream=True 按块读取，模拟播放器启动后的首段缓冲。

    指标说明:
    - QPS: 流媒体接口整体吞吐
    - P95/P99: 完成首段读取的耗时分位
    - TTFB(P95): 首包时间分位，越低表示“点播后更快出声”
    """
    stage_name = stage["name"]
    concurrency = stage["concurrency"]
    total_requests = stage["total_requests"]
    samples = discover_sample_ids()
    song_id = samples["song_id"]
    if song_id is None:
        log(f"阶段 {stage_name} 未发现可用 song_id，跳过流式播放压测。", "WARN")
        return None

    scenario = {
        "name": "api_song_stream_only",
        "method": "GET",
        "url": f"{BASE_URL}/song/stream/{song_id}",
        # 通过 Range 控制读取窗口，避免下载整首歌
        "headers": {"Range": f"bytes=0-{STREAM_READ_BYTES - 1}"},
    }
    log(f"开始流式播放阶段: {stage_name} (并发={concurrency}, 总请求={total_requests})")

    stats = {"total": 0, "ok": 0, "fail": 0, "latencies": [], "ttfb": [], "status_codes": {}}

    def request_once():
        start = time.perf_counter()
        try:
            resp = requests.request(
                method=scenario["method"],
                url=scenario["url"],
                headers=scenario["headers"],
                timeout=REQUEST_TIMEOUT,
                stream=True,
            )
            try:
                first_chunk_ms = (time.perf_counter() - start) * 1000
                bytes_read = 0
                for chunk in resp.iter_content(chunk_size=8192):
                    if not chunk:
                        continue
                    bytes_read += len(chunk)
                    if bytes_read >= STREAM_READ_BYTES:
                        break
                elapsed_ms = (time.perf_counter() - start) * 1000
                ok = 200 <= resp.status_code < 400 and bytes_read > 0
                return ok, resp.status_code, elapsed_ms, first_chunk_ms, ("" if ok else "stream returned no bytes")
            finally:
                resp.close()
        except Exception as e:
            elapsed_ms = (time.perf_counter() - start) * 1000
            return False, "EXC", elapsed_ms, elapsed_ms, str(e)

    start_all = time.perf_counter()
    with ThreadPoolExecutor(max_workers=concurrency) as executor:
        futures = [executor.submit(request_once) for _ in range(total_requests)]
        for future in as_completed(futures):
            ok, status_code, elapsed_ms, ttfb_ms, err = future.result()
            stats["total"] += 1
            stats["latencies"].append(elapsed_ms)
            stats["ttfb"].append(ttfb_ms)
            code = str(status_code)
            stats["status_codes"][code] = stats["status_codes"].get(code, 0) + 1
            if ok:
                stats["ok"] += 1
            else:
                stats["fail"] += 1
                if err:
                    log(f"api_song_stream_only 请求异常: {err}", "WARN")

    total_elapsed = time.perf_counter() - start_all
    total_qps = (total_requests / total_elapsed) if total_elapsed > 0 else 0.0
    avg_ms = statistics.mean(stats["latencies"]) if stats["latencies"] else 0.0
    p95_ms = percentile(stats["latencies"], 0.95)
    p99_ms = percentile(stats["latencies"], 0.99)
    p95_ttfb_ms = percentile(stats["ttfb"], 0.95)
    success_rate = (stats["ok"] / stats["total"] * 100) if stats["total"] else 0.0

    print(f"\n===== stream-{stage_name} 结果汇总 =====")
    print(f"总请求数: {total_requests}, 并发: {concurrency}, 阶段总QPS: {total_qps:.2f} req/s")
    print(f"请求数={stats['total']}, 成功={stats['ok']}, 失败={stats['fail']}, 成功率={success_rate:.2f}%")
    print(f"平均耗时={avg_ms:.2f}ms, P95={p95_ms:.2f}ms, P99={p99_ms:.2f}ms")
    # TTFB 重点衡量“播放器是否能快速开始播放”
    print(f"首包P95(TTFB)={p95_ttfb_ms:.2f}ms, 状态码分布={stats['status_codes']}")

    return {
        "stage_name": f"stream-{stage_name}",
        "concurrency": concurrency,
        "total_requests": total_requests,
        "elapsed_sec": total_elapsed,
        "total_qps": total_qps,
    }


def create_test_users():
    """批量创建测试用户，已存在用户会自动跳过。"""
    log("开始创建测试用户...")
    users = []
    for i in range(1, TEST_USER_COUNT + 1):
        username = f"{TEST_USER_PREFIX}{i:03d}"  # testuser001 ... testuser010
        email = TEST_USER_EMAIL_TEMPLATE.format(i)
        users.append(username)

        try:
            reg_payload = {
                "username": username,
                "password": TEST_USER_PASSWORD,
                "email": email,
            }
            res = requests.post(f"{BASE_URL}/auth/register", json=reg_payload, timeout=REQUEST_TIMEOUT)
            code = None
            message = ""
            try:
                body = res.json()
                code = body.get("code")
                message = body.get("message", "")
            except Exception:
                message = res.text

            if res.status_code == 200 and code == 1000:
                log(f"用户注册成功: {username}")
            else:
                log(f"用户注册跳过(可能已存在): {username} - {res.status_code} {message}", "WARN")
        except Exception as e:
            log(f"注册请求异常 {username}: {e}", "ERROR")
    log("测试用户创建流程结束。")
    return users


def run_pressure_test(test_users, stage):
    """执行单个阶段压测并输出各场景成功率、延迟与 QPS。"""
    stage_name = stage["name"]
    concurrency = stage["concurrency"]
    total_requests = stage["total_requests"]
    log(f"开始压测阶段: {stage_name} (并发={concurrency}, 总请求={total_requests})")
    if not test_users:
        log("没有可用测试用户，无法执行登录场景压测。", "ERROR")
        return None

    scenarios = build_scenarios(test_users)

    stats = {
        s["name"]: {"total": 0, "ok": 0, "fail": 0, "latencies": [], "ttfb": [], "status_codes": {}}
        for s in scenarios
    }

    def request_once(index):
        # 用轮询方式把请求均匀分配到各场景
        scenario = scenarios[index % len(scenarios)]
        sname = scenario["name"]
        payload = None
        params = scenario.get("params")
        headers = scenario.get("headers")
        stream_probe = bool(scenario.get("stream_probe"))
        if sname == "api_user_login":
            username = test_users[index % len(test_users)]
            payload = {"username": username, "password": TEST_USER_PASSWORD}

        start = time.perf_counter()
        try:
            # 对流媒体场景进行小块读取，模拟播放器首段缓冲行为
            if stream_probe:
                resp = requests.request(
                    method=scenario["method"],
                    url=scenario["url"],
                    params=params,
                    headers=headers,
                    timeout=REQUEST_TIMEOUT,
                    stream=True,
                )
                try:
                    first_chunk_ms = (time.perf_counter() - start) * 1000
                    bytes_read = 0
                    for chunk in resp.iter_content(chunk_size=8192):
                        if not chunk:
                            continue
                        bytes_read += len(chunk)
                        if bytes_read >= STREAM_READ_BYTES:
                            break
                    elapsed_ms = (time.perf_counter() - start) * 1000
                    ok = 200 <= resp.status_code < 400 and bytes_read > 0
                    return {
                        "name": sname,
                        "ok": ok,
                        "status_code": resp.status_code,
                        "latency_ms": elapsed_ms,
                        "error": "" if ok else "stream returned no bytes",
                        "ttfb_ms": first_chunk_ms,
                    }
                finally:
                    resp.close()

            res = requests.request(
                method=scenario["method"],
                url=scenario["url"],
                params=params,
                headers=headers,
                json=payload,
                timeout=REQUEST_TIMEOUT,
            )
            elapsed_ms = (time.perf_counter() - start) * 1000
            ok = 200 <= res.status_code < 400
            return {
                "name": sname,
                "ok": ok,
                "status_code": res.status_code,
                "latency_ms": elapsed_ms,
                "error": "",
                "ttfb_ms": elapsed_ms,
            }
        except Exception as e:
            elapsed_ms = (time.perf_counter() - start) * 1000
            return {
                "name": sname,
                "ok": False,
                "status_code": "EXC",
                "latency_ms": elapsed_ms,
                "error": str(e),
                "ttfb_ms": elapsed_ms,
            }

    start_all = time.perf_counter()
    with ThreadPoolExecutor(max_workers=concurrency) as executor:
        futures = [executor.submit(request_once, i) for i in range(total_requests)]
        for future in as_completed(futures):
            result = future.result()
            bucket = stats[result["name"]]
            bucket["total"] += 1
            bucket["latencies"].append(result["latency_ms"])
            bucket["ttfb"].append(result.get("ttfb_ms", result["latency_ms"]))
            code = str(result["status_code"])
            bucket["status_codes"][code] = bucket["status_codes"].get(code, 0) + 1
            if result["ok"]:
                bucket["ok"] += 1
            else:
                bucket["fail"] += 1
                if result["error"]:
                    log(f"{result['name']} 请求异常: {result['error']}", "WARN")

    total_elapsed = time.perf_counter() - start_all
    total_qps = (total_requests / total_elapsed) if total_elapsed > 0 else 0.0
    log(f"阶段 {stage_name} 完成，总耗时: {total_elapsed:.2f}s")
    print(f"\n===== {stage_name} 结果汇总 =====")
    print(
        f"总请求数: {total_requests}, 并发: {concurrency}, "
        f"阶段总QPS: {total_qps:.2f} req/s"
    )

    scenario_summary = {}
    for name, data in stats.items():
        latencies = data["latencies"]
        avg_ms = statistics.mean(latencies) if latencies else 0.0
        p95_ms = percentile(latencies, 0.95)
        p99_ms = percentile(latencies, 0.99)
        p95_ttfb_ms = percentile(data["ttfb"], 0.95)
        success_rate = (data["ok"] / data["total"] * 100) if data["total"] else 0.0
        scenario_qps = (data["total"] / total_elapsed) if total_elapsed > 0 else 0.0

        print(f"\n场景: {name}")
        print(
            f"  请求数={data['total']}, 成功={data['ok']}, 失败={data['fail']}, "
            f"成功率={success_rate:.2f}%"
        )
        print(f"  场景QPS={scenario_qps:.2f} req/s")
        print(f"  平均耗时={avg_ms:.2f}ms, P95={p95_ms:.2f}ms, P99={p99_ms:.2f}ms")
        print(f"  首包P95(TTFB)={p95_ttfb_ms:.2f}ms")
        print(f"  状态码分布={data['status_codes']}")

        scenario_summary[name] = {
            "requests": data["total"],
            "success_rate": success_rate,
            "qps": scenario_qps,
            "p95_ms": p95_ms,
            "p99_ms": p99_ms,
            "p95_ttfb_ms": p95_ttfb_ms,
        }

    return {
        "stage_name": stage_name,
        "concurrency": concurrency,
        "total_requests": total_requests,
        "elapsed_sec": total_elapsed,
        "total_qps": total_qps,
        "scenarios": scenario_summary,
    }


def run_task():
    """主流程：先注册测试用户，再执行分阶段自动升压压测。"""
    # 1) 创建测试用户
    test_users = create_test_users()

    # 2) 分阶段自动升压压测
    all_stage_results = []
    for stage in PRESSURE_STAGES:
        result = run_pressure_test(test_users, stage)
        if result is not None:
            all_stage_results.append(result)

    # 3) 独立执行歌曲流式播放压测
    stream_stage_results = []
    for stage in PRESSURE_STAGES:
        stream_result = run_song_stream_test(stage)
        if stream_result is not None:
            stream_stage_results.append(stream_result)

    # 4) 输出阶段总览对比
    if all_stage_results:
        print("\n===== 分阶段总览 =====")
        print("阶段 | 并发 | 总请求 | 耗时(s) | 总QPS")
        for item in all_stage_results:
            print(
                f"{item['stage_name']} | {item['concurrency']} | {item['total_requests']} | "
                f"{item['elapsed_sec']:.2f} | {item['total_qps']:.2f}"
            )

    if stream_stage_results:
        print("\n===== 流式播放分阶段总览 =====")
        print("阶段 | 并发 | 总请求 | 耗时(s) | 总QPS")
        for item in stream_stage_results:
            print(
                f"{item['stage_name']} | {item['concurrency']} | {item['total_requests']} | "
                f"{item['elapsed_sec']:.2f} | {item['total_qps']:.2f}"
            )

    log("全部任务执行结束。")


if __name__ == "__main__":
    run_task()
