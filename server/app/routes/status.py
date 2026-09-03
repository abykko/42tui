# app/routes/status.py

from flask import Blueprint, jsonify

status_bp = Blueprint("status", __name__)

@status_bp.route("/status")
def status():
    return jsonify({"status": "ok"}), 200
