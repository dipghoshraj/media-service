# import socket
# import ssl

# bindsocket = socket.socket()
# bindsocket.bind(('0.0.0.0', 5050))
# bindsocket.listen(5)

# context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
# context.load_cert_chain("agni.local.internal.pem", "agni.local.internal-key.pem")
# print("Server listening on 0.0.0.0:5050")

# while True:
#     try:
#         newsocket, fromaddr = bindsocket.accept()
#         print("Accepted connection from", fromaddr)
#         try:
#             conn = context.wrap_socket(newsocket, server_side=True)
#             data = conn.recv(4096)
#             print("Raw TLS bytes received:", data[:100])
#             conn.send(b"HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello")
#             conn.close()
#         except ssl.SSLError as e:
#             print("SSL handshake failed:", e)
#             newsocket.close()
#     except Exception as e:
#         print("Error accepting connection:", e)


from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello():
    return "Hello HTTPS from Flask + Gunicorn + Nginx!"