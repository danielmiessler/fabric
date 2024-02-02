from flask import Flask, Blueprint
from instance.config import DevelopmentConfig
from flask_socketio import SocketIO

sockio = SocketIO()


def create_app(config_class=DevelopmentConfig):
    app = Flask(__name__)
    app.config.from_object(config_class)
    sockio.init_app(app)

    from app.chatgpt import bp as chatgpt_bp
    app.register_blueprint(chatgpt_bp, url_prefix='/chatgpt')

    from app.web_frontend import bp as web_frontend_bp
    app.register_blueprint(web_frontend_bp)

    from app.patterns import bp as patterns_bp
    app.register_blueprint(patterns_bp, url_prefix='/patterns')

    return app
