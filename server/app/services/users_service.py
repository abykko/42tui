from app.utils.browser import run_playwright
from app.config import PROFILE_URL

def json_bundler(objs: list[tuple[str, object]]):
    data = {}
    for key, obj in objs:
        if key in data and isinstance(data[key], dict) and isinstance(obj, dict):
            data[key].update(obj)
        else:
            data[key] = obj

    return data


def profile(username):

    if not username: return None

    resp_objs, user_data = [], None

    def task(page, context):

        def response_filter(resp):
            if f"/api/v1/users/{username}" in resp.url:
                print(resp.url)
                resp_objs.append((resp.url, resp.json()))

        page.on("response", response_filter)
        page.goto(f"{PROFILE_URL}/users/{username}")

        while len(resp_objs) < 8:
            page.wait_for_timeout(10)

        return json_bundler(resp_objs)

    user_data = run_playwright(task)

    return user_data
