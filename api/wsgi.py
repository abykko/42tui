import os
from app import create_app

app = create_app(
    secret=os.getenv("API_SECRET", "foo"),
    api_security=os.getenv("API_SECURITY", "true").lower() == "false"
)