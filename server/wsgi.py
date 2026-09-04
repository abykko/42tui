import os
from app import create_app

app = create_app(
    secret=os.getenv("SECRET")
)