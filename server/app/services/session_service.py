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

def log_in(username, password) -> bool:
    def task(page, context):
        try:
            page.goto(PROFILE_URL)

            page.wait_for_url(AUTH_URL)

            page.wait_for_selector('#username', timeout=10000)

            page.wait_for_selector('#password', timeout=10000)

            page.fill('#username', username)
            page.fill('#password', password)

            page.wait_for_selector('#kc-login', timeout=10000)

            page.click('#kc-login')

            page.wait_for_url(PROFILE_URL, timeout=10000)

            context.storage_state(path=SESSION_FILE)

            return True
            
        except Exception as e:
            return False

    return run_playwright(task, headless=True, auth=False)
