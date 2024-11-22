from flask import Flask
from app.config import Config
from app.controller import upload_blueprint

def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)  # Load config settings

    # Register Blueprints
    app.register_blueprint(upload_blueprint)

    return app
