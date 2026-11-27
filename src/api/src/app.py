import json
import logging
import os
import psycopg2
from flask import Flask, Response

app = Flask(__name__)

logging.basicConfig(level=logging.INFO)

DB_CONFIG = {
    "host": os.getenv("DB_HOST", "db"),
    "port": int(os.getenv("DB_PORT", "5432")),
    "dbname": os.getenv("DB_NAME", "words_db"),
    "user": os.getenv("DB_USER", "words_app"),
    "password": os.getenv("DB_PASSWORD","12345")
}

def random_word(table):
    conn = psycopg2.connect(**DB_CONFIG)
    try:
        with conn.cursor() as cur:
            # Table names are fixed by the application, not user input.
            cur.execute(f"SELECT word FROM {table} ORDER BY random() LIMIT 1")
            row = cur.fetchone()

            if row is None:
                raise LookupError(table)

            return row[0]
    finally:
        conn.close()


def make_handler(table):
    def handler():
        word = random_word(table)

        response = {"word": word}
        logging.info(json.dumps(response))

        return Response(
            json.dumps(response),
            status=200,
            mimetype="application/json",
            headers={
                "Cache-Control": "private, no-cache, no-store, must-revalidate, max-age=0",
                "Pragma": "no-cache",
            },
        )

    return handler


app.add_url_rule("/noun", "noun", make_handler("nouns"))
app.add_url_rule("/verb", "verb", make_handler("verbs"))
app.add_url_rule("/adjective", "adjective", make_handler("adjectives"))

@app.get("/healthz")
def healthz(): return "ok", 200

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
