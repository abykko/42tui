from utils.browser import run_playwright

def json_bundler(objs: list[tuple[str, object]]):
    data = {}

    for url, obj in objs:
        key = url.split("?")[0].strip("/").split("/")[-1]

        if key in data and isinstance(data[key], dict) and isinstance(obj, dict):
            data[key].update(obj)
        else:
            data[key] = obj

    return data


def get(username: str = "iamrani-"):

    resp_objs = []
    user_data = None

    def task(page, context):
        nonlocal user_data

        def response_filter(resp):
            nonlocal resp_objs
            try:
                if (
                    f"/api/v1/users/{username}" in resp.url
                    and resp.request.method == "GET"
                    and resp.status == 200
                ):
                    resp_objs.append((resp.url, resp.json()))
            except Exception as e:
                print(f"[response_filter] Error: {e}")

        page.on("response", response_filter)

        url = f"https://profile-v3.intra.42.fr/users/{username}"
        page.goto(url)

        while len(resp_objs) < 8: page.wait_for_timeout(20)

        user_data = json_bundler(resp_objs)

        return user_data

    user_data = run_playwright(task)

    return user_data