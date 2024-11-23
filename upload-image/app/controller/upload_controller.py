from flask import Blueprint, request, Response, jsonify
from app.model.ImageMetadata import ImageMetadata
from app.concern.minio_upload import Uploader


upload_blueprint = Blueprint('upload_bp', __name__, url_prefix='/api/v1/image')

@upload_blueprint.route('/')
def index():
    return jsonify({}, 200)


@upload_blueprint.route('/upload', methods=['POST'])
def upload():
    file = request.files['file']
    image_id = file.filename
    user_id = request.POST.get('user_id')

    client = Uploader()
    upload_data = client.upload(image_id, user_id)
    return jsonify(upload_data, 201)