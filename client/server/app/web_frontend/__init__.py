from flask import Blueprint

bp = Blueprint('webFrontend', __name__)

if 1 == 1:
    from app.web_frontend import routes
