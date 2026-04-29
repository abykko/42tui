# app/routes/session.py

from flask import Blueprint, jsonify
from app.services.session_service import log_in, log_out
from app.utils.browser import is_expired

session_bp = Blueprint("session", __name__)

@session_bp.route("/session/login")
def login():
    return jsonify({"login": log_in()}), 200

@session_bp.route("/session/logout")
def logout():
    return jsonify({"logout": log_out()}), 200

@session_bp.route("/session/expired")
def expired():
    return jsonify({"expired": is_expired()}), 200