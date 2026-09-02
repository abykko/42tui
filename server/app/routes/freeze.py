# app/routes/is_freeze.py

from flask import Blueprint, jsonify, request
from app.services.freeze_service import is_freeze

freeze_bp = Blueprint("freeze", __name__)

@freeze_bp.route("/freeze")
def freeze():
    is_user_freeze = is_freeze()
    print(is_user_freeze)
    if is_user_freeze is None:
        return jsonify({"error": "Unexpected error"}), 400
    return jsonify({"freeze": is_user_freeze}), 200
