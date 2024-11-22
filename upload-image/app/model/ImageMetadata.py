import datetime
from app.model import db


class ImageMetadata:
    id = db.Column(db.Integer, primary_key=True)
    key = db.Column(db.String)
    file_path = db.Column(db.String, nullable=False)
    upload_time = db.Column(db.DateTime, default=datetime.datetime.now())
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'), nullable=False)  # Foreign key to User

    def __repr__(self):
        return f"<ImageMetadate {self.key}>"