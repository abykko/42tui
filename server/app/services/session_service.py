from app.utils.browser import run_playwright
from app.config import SESSION_FILE, PROFILE_URL, AUTH_URL
import os

def log_out() -> bool:
    if os.path.exists(SESSION_FILE):
        try:
            os.remove(SESSION_FILE)
        except Exception as e:
            return False
    return True


def log_in() -> bool:
    def task(page, context):
        try:
            if not SESSION_FILE: return False

            page.goto(PROFILE_URL)
            page.wait_for_url(AUTH_URL, timeout=0)
            page.wait_for_url(PROFILE_URL, timeout=0)
            context.storage_state(path=SESSION_FILE)
        except Exception:
            return (False)
        return (True)

    return (run_playwright(task, headless=False, auth=False))
