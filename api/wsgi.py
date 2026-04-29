from app import create_app

app = create_app(
    secret="foo",
    api_security=True
)