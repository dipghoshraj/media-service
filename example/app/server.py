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
#         print(f"[server.py] TCP connection accepted from {fromaddr}")
#         try:
#             conn = context.wrap_socket(newsocket, server_side=True)
#             print(f"[server.py] TLS handshake OK — version={conn.version()} cipher={conn.cipher()}")

#             cert = conn.getpeercert()
#             if cert:
#                 print(f"[server.py] Client cert: {cert}")
#             else:
#                 print("[server.py] No client certificate presented")

#             data = conn.recv(4096)
#             if data:
#                 try:
#                     print(f"[server.py] Decrypted HTTP request ({len(data)} bytes):\n{data.decode('utf-8', errors='replace')}")
#                 except Exception:
#                     print(f"[server.py] Received {len(data)} bytes (non-UTF8)")
#             else:
#                 print("[server.py] Received 0 bytes after TLS handshake")

#             conn.sendall(b"HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello")
#             print("[server.py] Response sent")
#             conn.close()
#         except ssl.SSLError as e:
#             print(f"[server.py] TLS handshake failed from {fromaddr}: {e}")
#             newsocket.close()
#     except Exception as e:
#         print(f"[server.py] Error: {e}")


from flask import Flask, jsonify

app = Flask(__name__)

@app.route("/")
def hello():
    return "Hello HTTPS from Flask + Gunicorn + Nginx!"

@app.route("/todo")
def todo():
    return jsonify({
        "userId": 1,
        "id": 1,
        "title": "delectus aut autem", 
        "completed": False
    })