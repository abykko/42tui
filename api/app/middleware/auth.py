# app/middleware/auth.py

import time
import hmac
import hashlib
from flask import request, jsonify

def validate_sign(secret, timestamp, sign):
    now = int(time.time())

    if abs(now - int(timestamp)) > 5:
        return False

    msg = timestamp.encode()

    legit_sign = hmac.new(
        secret.encode(),
        msg,
        hashlib.sha256
    ).hexdigest()

    return hmac.compare_digest(sign, legit_sign)

def register_security(app, api_security, secret):

    @app.before_request
    def before_request():
        if not api_security:
            return

        timestamp = request.headers.get("X-Timestamp", "").strip()
        sign = request.headers.get("X-Signature", "").strip()

        if not timestamp or not sign:
            return jsonify({"error": "Unauthorized"}), 401

        if not validate_sign(secret, timestamp, sign):
            return jsonify({"error": "Invalid signature"}), 401
