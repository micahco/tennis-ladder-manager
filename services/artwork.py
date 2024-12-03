import zmq
import signal

ADDR = "tcp://*:9875"
context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind(ADDR)

print("Artwork server running...")

signal.signal(signal.SIGINT, signal.SIG_DFL)

ART = '''
   ,odOO"bo,
 ,dOOOP'dOOOb,
,O3OP'dOO3OO33,
P",ad33O333O3Ob
?833O338333P",d
`88383838P,d38'
 `Y8888P,d88P'
   `"?8,8P"'
'''

colors = {
    "red": 31,
    "green": 32,
    "yellow": 33,
    "blue": 34,
    "purple": 35,
    "cyan": 36,
    "white": 37
}

while True:
    message = socket.recv_json()
    print("Server Received:", message)

    data = ART
    color = colors[message]
    if color is not None:
        data = f"\033[{color}m" + ART + "\033[0m"

    socket.send_string(data)
    