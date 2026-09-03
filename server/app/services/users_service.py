import time
from flask import jsonify
from app.utils.browser import run_playwright
from app.config import PROFILE_URL


EXPECTED_USER_ENDPOINTS = {
    "": "profile",
    "/summary": "summary",

    "/projects/marked": "projects",

    "/pace": "pace",

    # Currently is not possible to obtain /marked and /ongoing
    # on the same page because you only can see /marked in your /profile
    # page. /ongoing is only visible in your /me
    # "/projects/ongoing": "projects",

    # "/patroning": "patroning",
    # "/patroned": "patroned",
    # "/achievements": "achievements",
    # "/partnerships": "partnerships",
    # "/accreditations": "accreditations",
}


def is_complete(username: str, responses: list[tuple[str, object]]) -> bool:
    found = set()

    for path, _ in responses:
        found.add(path)

    return found == set(EXPECTED_USER_ENDPOINTS.keys())


def json_bundler(objs: list[tuple[str, object]]):
    data = {}

    for path, obj in objs:

        key = EXPECTED_USER_ENDPOINTS[path]

        if key in data and isinstance(data[key], dict) and isinstance(obj, dict):
            data[key].update(obj)
        else:
            data[key] = obj

    return data


def profile(username):
    if not username:
        return None

    resp_objs = []

    def task(page, context):

        def response_filter(resp):

            if resp.status != 200:
                return

            prefix1 = f"/api/v1/users/{username}"
            pace_system = "pace-system.42.fr"

            if prefix1 in resp.url:
                path = resp.url.split(prefix1, 1)[1].split("?", 1)[0]

            elif pace_system in resp.url:
                path = "/pace"
            else:
                return

            # Solo aceptar endpoints registrados
            if path not in EXPECTED_USER_ENDPOINTS:
                return

            try:
                # evitar duplicados por path
                if not any(saved_path == path for saved_path, _ in resp_objs):
                    resp_objs.append((path, resp.json()))

                    print("\nDetectadas:")
                    for p, _ in resp_objs:
                        print(p)

            except:
                pass

        page.on("response", response_filter)

        # user page
        page.goto(f"{PROFILE_URL}/users/{username}")

        start_time = time.perf_counter()
        timeout_limit = 20.0

        while not is_complete(username, resp_objs):

            duration = time.perf_counter() - start_time

            if duration >= timeout_limit:
                return {
                    "error": "timeout",
                    "duration": f"{duration:.2f}s",
                    "found_items": len(resp_objs),
                    "paths": [p for p, _ in resp_objs],
                }

            page.wait_for_timeout(100)

        return json_bundler(resp_objs)

    user_data = run_playwright(task)

    return jsonify(user_data)