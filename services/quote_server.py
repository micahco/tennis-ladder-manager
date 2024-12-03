import zmq
import json
import random
import signal

ADDR = "tcp://*:9876"
server_running = True

with open('quotes.json', 'r') as f:
    quotes = json.load(f)


def get_random_sports_quote(filters=None):
    if filters:
        filtered_quotes = []

        for quote in quotes:
            is_match = True

            for key, value in filters.items():
                if quote.get(key) != value:
                    is_match = False
                    break
            if is_match:
                filtered_quotes.append(quote)

        if filtered_quotes:
            return random.choice(filtered_quotes)
        else:
            return None

    return random.choice(quotes)


context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind(ADDR)

print("Quote server running...")

signal.signal(signal.SIGINT, signal.SIG_DFL)

while server_running:

    message = socket.recv_json()
    print("Server Received:", message)

    filters = message.get('filters', {})

    quote = get_random_sports_quote(filters)

    if quote:
        socket.send_json({
            "code": 200,
            "data": quote,
            "message": "Quote successfully returned."
        })
    else:
        socket.send_json({
            "code": 404,
            "data": "",
            "message": "No quote found for the specified conditions."
        })
