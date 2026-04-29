# app/routes/slots.py

from flask import Blueprint, jsonify, request
from app.services.slots_service import create_slot

slots_bp = Blueprint("slots", __name__)

@slots_bp.route("/slots/create")
def slots_create():
    begin_at = request.args.get("begin_at")
    end_at = request.args.get("end_at")

    create_slot(begin_at, end_at)
    return jsonify({"status": "created"}), 200
