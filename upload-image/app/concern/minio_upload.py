from minio import Minio
from minio.error import S3Error
import uuid, sys
from app.config import Config
from app.model import db
from app.model.Image import Image
from app.model.ImageMetadata import ImageMetadata

class Uploader:

    def __init__(self):
        self.client = Minio("minio:9000",access_key="minioadmin",secret_key="minioadmin", secure=False)
        
    
    def upload_iamge(self, file, user_id: str, post_id: str = None):

        key =uuid.uuid4().hex
        file_name = key + file.filename
        extension = file.filename.split('.')[1]


        # if bucket is not present then create one
        found = self.client.bucket_exists(Config.BUCKET_NAME)
        if not found: self.client.make_bucket(Config.BUCKET_NAME)

        try:
            meta_data = ImageMetadata(key=file_name, file_bucket=Config.BUCKET_NAME, size= sys.getsizeof(file), extension=extension)
            db.session.add(meta_data)
            db.session.flush()


            image = Image(name= file_name, url ="" , user_id= user_id, image_matadata_id=meta_data.id)
            db.session.add(image)
            db.session.flush()

            pobj = self.client.put_object(bucket_name=Config.BUCKET_NAME, object_name=file_name, data=file, length=sys.getsizeof(file))
            
            meta_data.image_id =image.id
            meta_data.etag = pobj.etag
            meta_data.version = pobj.version_id

            db.session.commit()
            return {"image": file_name}
        except Exception as e:
            self.client.remove_object(Config.BUCKET_NAME, file_name)
            db.session.rollback()
            raise Exception(f"An error occurred of image add {e.__repr__()}")
        finally:
            # Ensure the session is closed even if there's an error
            db.session.remove()

