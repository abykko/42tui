# app/routes/users.py

from flask import Blueprint, jsonify, request
from app.services.users_service import profile

users_bp = Blueprint("users", __name__)

@users_bp.route("/users/profile")
def user():
    user = request.args.get("user")
    user_data = profile(username=user)
    if user_data is None:
        return jsonify({"error": "Please, specify a user."}), 400
    return profile(username=user), 200
