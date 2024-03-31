import socket


class HTTP_Server:
    def __init__(self, host='127.0.0.1', port=3000):
        self.host = host
        self.port = port

    def init(self):
        print('Initializing server...')
        server_socket = socket.create_server(
            (self.host, self.port), reuse_port=True)
        conn, addr = server_socket.accept()
        while conn:
            print(f'Process: {addr[1]} - Connecting to {addr[0]}...')
            while True:
                data = conn.recv(1024)
                rx_data = data.decode().split("\n")[0]
                http_method = rx_data.split()[0]
                path = rx_data.split()[1]
                http_version = rx_data.split()[2]
                if not data:
                    break
                print(
                    f'Data packet received: {data} \n HTTP Method: {http_method},\n Path: {path},\n HTTP Version: {http_version}')

                if path == '/':
                    ack_msg = bytes('HTTP/1.1 200 OK\r\n\r\n',
                                    encoding='utf-8')
                    conn.sendall(ack_msg)
                else:
                    ack_msg = bytes(
                        'HTTP/1.1 400 Not Found\r\n\r\n', encoding='utf-8')
                    conn.sendall(ack_msg)
