# run.py

import os
from app import create_app

app = create_app(
    secret="foo",
    api_security=True
)

if __name__ == "__main__":
    app.run(
        host="127.0.0.1",
        port=6742,
        debug=True,
    )
