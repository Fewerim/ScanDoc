import argparse
from datetime import datetime
from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/health", methods=["GET"])
def health():
    return jsonify({"status": "ok"}), 200

@app.route("/recognize", methods=["POST"])
def recognize():
    # файл приходит как поле "file" в multipart/form-data
    uploaded = request.files.get("file")
    if uploaded is None:
        return jsonify({"error": "no file"}), 400

    # здесь можно ничего не распознавать, просто вернуть фиктивные данные
    result = {
        "filename": uploaded.filename,
        "doc_type": "stub",
        "processed_at": datetime.utcnow().isoformat() + "Z",
        "fields": {
            "example": "value"
        }
    }
    return jsonify(result), 200

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--port", type=int, required=True)
    args = parser.parse_args()

    # host=0.0.0.0 чтобы слушать все интерфейсы
    app.run(host="0.0.0.0", port=args.port)

if __name__ == "__main__":
    main()
