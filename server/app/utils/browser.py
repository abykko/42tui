from playwright.sync_api import sync_playwright
from app.config import SESSION_FILE, PROFILE_URL, AUTH_URL
import json
import os

def authenticate():
    try:
        with open(SESSION_FILE, "r", encoding="utf-8") as f:
            data = json.load(f)
            f.close()
        return data
    except Exception as e:
        return None


def is_expired():

    if not os.path.exists(SESSION_FILE): return True

    def task(page, context):
        try:
            page.goto(PROFILE_URL)
            page.wait_for_url(AUTH_URL, timeout=3000)
            return True
        except Exception:
            # If we reach the wait_for_url timeout
            return False
    
    return run_playwright(task)


def run_playwright(
    task,
    headless=True,
    auth=True,
):

    state = None
    if auth: state = authenticate()
    
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=headless)
        context = browser.new_context(storage_state=state)
        page = context.new_page()

        try:
            return task(page, context)
        finally:
            browser.close()
