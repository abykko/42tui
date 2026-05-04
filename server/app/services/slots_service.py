from app.utils.browser import run_playwright
from app.config import SLOTS_URL

def create_slot(begin_at, end_at):

    if not begin_at or not end_at:
        return None

    data = {
        "slot[begin_at]": begin_at,
        "slot[end_at]": end_at,
    }

    def task(page, context):

        page = context.new_page()

        page.goto(
            SLOTS_URL,
            wait_until="commit",
            timeout=10_000
        )

        csrf_token = page.locator('meta[name="csrf-token"]').get_attribute('content')
        if not csrf_token:
            return

        response = context.request.post(
            f"{SLOTS_URL}.json",
            headers={
                "X-CSRF-Token": csrf_token,
                "X-Requested-With": "XMLHttpRequest",
                "Referer": SLOTS_URL,
            },
            form=data,
        )

    result = run_playwright(task)

    return result
