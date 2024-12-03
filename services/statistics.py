import zmq
import signal

ADDR = "tcp://*:9877"
context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind(ADDR)

print("Statistics server running...")

signal.signal(signal.SIGINT, signal.SIG_DFL)

while True:
    message = socket.recv_json()
    print("Server Received:", message)

    if isinstance(message, list) and len(message) == 2 and all(isinstance(x, int) for x in message):
        wins, lose = message
        total = wins + lose
        per = int((wins / total) * 100)
        msg = f"Wins:\t{wins}\nLosses:\t{lose}\nRatio:\t{per}%"
        socket.send_string(msg)

    else: # Invalid request
        socket.send_string("ERROR: invalid request")
    