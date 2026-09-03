# app/routes/session.py

from flask import Blueprint, jsonify, request
from app.services.session_service import log_in, log_out
from app.utils.browser import is_expired

session_bp = Blueprint("session", __name__)

@session_bp.route("/session/login", methods=["GET"])
def login():
    username = request.args.get("username")
    password = request.args.get("password")

    if not username or not password:
        return jsonify({
            "error": "No credentials",
            "usage": "/session/login?username=XXX&password=YYY"
        }), 400

    success = log_in(username, password)
    
    return jsonify({"login": success}), 200 if success else 401

@session_bp.route("/session/logout")
def logout():
    return jsonify({"logout": log_out()}), 200

@session_bp.route("/session/expired")
def expired():
    return jsonify({"expired": is_expired()}), 200