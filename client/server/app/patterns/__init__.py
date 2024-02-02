from flask import Blueprint

bp = Blueprint('auth', __name__)
if 1 == 1:
    from app.patterns import routes
