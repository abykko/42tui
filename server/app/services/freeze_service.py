import time
from app.utils.browser import run_playwright


def is_freeze() -> bool:
    result = None
    timeout = time.time() + 10 # 10 seconds

    def task(page, context):
        def response_filter(resp):
            nonlocal result
            try:
                data = resp.json()

                if isinstance(data, dict) and "is_freeze" in data:
                    result = data["is_freeze"]

            except Exception:
                pass

        page.on("response", response_filter)

        page.goto("https://freeze.42.fr/")

        while time.time() < timeout:
            if result != None:
                break
            page.wait_for_timeout(100)

    run_playwright(task)

    return result