from utils.browser import run_playwright
import os

SESSION_FILE = "security/session.json"

def log_out():
    if os.path.exists(SESSION_FILE):
        try:
            os.remove(SESSION_FILE)
        except Exception as e:
            pass

def log_in():
    def task(page, context):
        page.goto("https://profile-v3.intra.42.fr/")
        page.wait_for_url("https://auth.42.fr/**", timeout=0)
        page.wait_for_url("https://profile-v3.intra.42.fr/", timeout=0)
        context.storage_state(path=SESSION_FILE)

    run_playwright(task, headless=False)