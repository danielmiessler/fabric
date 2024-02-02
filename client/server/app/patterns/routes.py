from flask import request, jsonify
from app.patterns import bp
import os


@bp.route("/get", methods=['GET'])
def get():
    try:
        config = os.path.expanduser("~/.config/fabric/patterns")
        patterns = os.listdir(config)
        pattern_list = {}
        for pattern in patterns:
            with open(f"{config}/{pattern}/system.md", "r") as f:
                pattern_list[pattern] = f.read()
        return jsonify(pattern_list)
    except Exception as e:
        return jsonify(message="Patterns not found"), 404
