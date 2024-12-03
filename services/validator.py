import zmq
import signal

ADDR = "tcp://*:9878"
context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind(ADDR)

print("Validator server running...")

signal.signal(signal.SIGINT, signal.SIG_DFL)

def is_valid_set(score1, score2):
    if score1 < 0 or score2 < 0:
        return False

    if score1 < 6 and score2 < 6:
        return False

    if abs(score1 - score2) >= 2 and (score1 == 6 or score2 == 6):
        return True

    if (score1 == 7 and score2 == 6) or (score1 == 6 and score2 == 7):
        return True

    return False

while True:
    message = socket.recv_json()
    print("Server Received:", message)

    if isinstance(message, list) and len(message) == 2 and all(isinstance(x, int) for x in message):
        score1, score2 = message
        is_valid = is_valid_set(score1, score2)
        socket.send_json(is_valid)

    else:
        socket.send_json(False)
    