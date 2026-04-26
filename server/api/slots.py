from utils.browser import run_playwright

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
            "https://profile.intra.42.fr/slots",
            wait_until="commit",
            timeout=10_000
        )

        csrf_token = page.locator('meta[name="csrf-token"]').get_attribute('content')
        if not csrf_token:
            return

        response = context.request.post(
            "https://profile.intra.42.fr/slots.json",
            headers={
                "X-CSRF-Token": csrf_token,
                "X-Requested-With": "XMLHttpRequest",
                "Referer": "https://profile.intra.42.fr/slots",
            },
            form=data,
        )

    result = run_playwright(task)

    return result
