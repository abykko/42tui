# app/__init__.py

from flask import Flask
from app.routes.status import status_bp
from app.routes.session import session_bp
from app.routes.users import users_bp
from app.routes.slots import slots_bp
from app.middleware.auth import register_security

def create_app(api_security: bool, secret: str):

    app = Flask(__name__)

    # Middleware / seguridad
    register_security(app, api_security, secret)

    # Blueprints
    app.register_blueprint(status_bp)
    app.register_blueprint(session_bp)
    app.register_blueprint(users_bp)
    app.register_blueprint(slots_bp)

    return app
