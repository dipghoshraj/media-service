from minio import Minio
from minio.error import S3Error
import uuid
from config import Config
from app.model import db
from app.model.Image import Image
from app.model.ImageMetadata import ImageMetadata

class Uploader:

    def __init__(self):
        self.client = Minio("minio:9000",
                                access_key="minioadmin",
                                secret_key="minioadmin")
        
    
    def upload_iamge(self, file, user_id, post_id):

        key =uuid.uuid4().hex
        file_name = key + file.filename


        # if bucket is not present then create one
        found = self.client.bucket_exists(Config.BUCKET_NAME)
        if not found: self.client.make_bucket(Config.BUCKET_NAME)

        try:
            ""
            image = Image(name= file_name, url ="" , user_id= user_id)
            db.session.add(image)
            db.session.flush()

            meta_data = ImageMetadata(key=file_name, file_path=f"{Config.BUCKET_NAME}/{file_name}", image_id= image.id)
            db.session.add(meta_data)
            db.session.flush()

            image.image_matadata_id = meta_data.id

            self.client.put_object(Bucket=Config.BUCKET_NAME, Key=file_name, Body=file)
            db.session.commit()

            return {"image": file_name}
        except Exception as e:
            self.client.remove_object(Config.BUCKET_NAME, file_name)
            db.session.rollback()
            raise Exception(f"An error occurred of image add {e.__repr__()}")
        finally:
            # Ensure the session is closed even if there's an error
            db.session.remove()

