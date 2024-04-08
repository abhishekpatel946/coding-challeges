import socket
import time

from utils import crlf, lf, http_protocol, http_version, http_status_ok, http_status_not_found


class HTTP_Server:
    def __init__(self, host='127.0.0.1', port=3000):
        self.host = host
        self.port = port
        self._socket = None
        self.retry_timeout = 25

    def retry(self, timeout=0):
        for remaining_time in range(timeout+1, 0, -1):
            print(
                f'Retrying to configure socket with {self.host}: {self.port} after {remaining_time} seconds', end="\r")
            time.sleep(1)

    def configure_socket(self):
        try:
            # configure the _socket
            self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self._socket.bind((self.host, self.port))
        except socket.error as e:
            print(f"\nCouldn't bind socket to {self.host}: {e} \n")
            self._socket.close()
            self.retry(self.retry_timeout)
            self.configure_socket()

    def init(self):
        try:
            self.configure_socket()
            self._socket.listen()
            print('\n\nWaiting for connection...')

            client_socket, client_address = self._socket.accept()  # wait for connection/client
            _hostname, _port = client_socket.getsockname()

            with client_socket:
                print(
                    f'Process: {client_address[1]} - Connecting to {_hostname}:{_port}...')
                while True:
                    data = client_socket.recv(1024)
                    rx_data = data.decode().split("\n")[0]
                    _http_method = rx_data.split()[0]
                    path = rx_data.split()[1]
                    _http_version = rx_data.split()[2]
                    print(
                        f'{lf} Data packet received: {data},{lf} HTTP Method: {_http_method},{lf} Path: {path},{lf} HTTP Version: {_http_version}')

                    splitted_path = path.split('/')
                    for _path in splitted_path:
                        # emtpy string in splitted_path
                        if _path == '':
                            pass
                        # check for root urls
                        elif path == '/':
                            ack_msg = bytes(f'{http_protocol}{http_version} {http_status_ok} {crlf}',
                                            encoding='utf-8')
                            client_socket.sendall(ack_msg)
                        elif _path == 'echo':
                            get_data_in_url = rx_data.split(
                                ' ')[1].split('/')[-1]
                            ack_msg = bytes(f'{http_protocol}{http_version}{lf}Content-Type: text/plain \nContent-Length: {len(get_data_in_url)} {lf}{get_data_in_url} {lf}',
                                            encoding='utf-8')
                            client_socket.sendall(ack_msg)
                        else:
                            ack_msg = bytes(
                                f'{http_protocol}{http_version} {http_status_not_found} {crlf}', encoding='utf-8')
                            client_socket.sendall(ack_msg)

        except Exception as e:
            print('Some error occurred', e)
