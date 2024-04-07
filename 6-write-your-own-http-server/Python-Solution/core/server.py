import socket
import time


class HTTP_Server:
    def __init__(self, host='127.0.0.1', port=3000):
        self.host = host
        self.port = port
        self._socket = None
        self.retry_timeout = 25

    def configure_socket(self):
        try:
            # configure the _socket
            self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self._socket.bind((self.host, self.port))
        except socket.error as e:
            print(f"\n Couldn't bind socket to {self.host}: {e}")
            self._socket.close()
            print(
                f'\n Retrying to configure socket with {self.host}: {self.port} after {self.retry_timeout} seconds')
            time.sleep(self.retry_timeout)
            self.configure_socket()

    def init(self):
        try:
            self.configure_socket()
            self._socket.listen()
            print('\n\nWaiting for connection...')

            client_socket, client_address = self._socket.accept()  # wait for connection/client

            with client_socket:
                print(
                    f'Process: {client_address[1]} - Connecting to {client_address[0]}...')
                while True:
                    data = client_socket.recv(1024)
                    rx_data = data.decode().split("\n")[0]
                    http_method = rx_data.split()[0]
                    path = rx_data.split()[1]
                    http_version = rx_data.split()[2]
                    print(
                        f'Data packet received: {data},/n HTTP Method: {http_method},/n Path: {path},/n HTTP Version: {http_version}')

                    if path == '/':
                        ack_msg = bytes('HTTP/1.1 200 OK\r\n\r\n',
                                        encoding='utf-8')
                        client_socket.sendall(ack_msg)
                    else:
                        ack_msg = bytes(
                            'HTTP/1.1 400 Not Found\r\n\r\n', encoding='utf-8')
                        client_socket.sendall(ack_msg)

        except Exception as e:
            print('Some error occurred', e)
