from flask import Blueprint, request, Response, jsonify
from app.model.ImageMetadata import ImageMetadata

upload_blueprint = Blueprint('upload_bp', __name__, url_prefix='/api/v1/image')

@upload_blueprint.route('/')
def index():
    return jsonify({}, 200)