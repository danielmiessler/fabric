from flask import render_template
from . import bp  # Import the blueprint


@bp.route('/')
def index():
    return render_template('index.html')
