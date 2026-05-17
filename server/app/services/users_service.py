import time
from flask import jsonify
from app.utils.browser import run_playwright
from app.config import PROFILE_URL

def json_bundler(objs: list[tuple[str, object]]):
    data = {}
    for key, obj in objs:

        # rename
        if "projects" in key:
            key = "projects"

        if "summary" in key:
            key = "summary"

        if key in data and isinstance(data[key], dict) and isinstance(obj, dict):
            data[key].update(obj)
        else:
            data[key] = obj
    return data

def profile(username):
    if not username: return None

    resp_objs = []

    def task(page, context):
        def response_filter(resp):
            if f"/api/v1/users/{username}" in resp.url:
                try:
                    resp_objs.append((resp.url, resp.json()))
                except:
                    pass

        page.on("response", response_filter)
        page.goto(f"{PROFILE_URL}/users/{username}")

        start_time = time.perf_counter()
        timeout_limit = 20.0

        while len(resp_objs) < 8:
            duration = time.perf_counter() - start_time
            
            if duration >= timeout_limit:
                return {
                    "error": "timeout",
                    "duration": f"{duration:.2f}s",
                    "found_items": len(resp_objs)
                }
            
            page.wait_for_timeout(100)

        return json_bundler(resp_objs)

    user_data = run_playwright(task)

    return jsonify(user_data)
