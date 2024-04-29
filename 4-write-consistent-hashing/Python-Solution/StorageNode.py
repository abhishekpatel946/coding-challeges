import requests


class StorageNode:
    def __init__(self, name=None, host=None):
        self.name = name
        self.host = host

    def get_file(self, path):
        return requests.get(f'https://{self.host}:1231/{path}').text

    def put_file(self, path):
        with open(path, 'r') as file:
            content = file.read()
            return requests.post(f'https://{self.host}:1231/{path}', body=content).text


# storage_nodes holding instances of actual storage node objects
storage_nodes = [
    StorageNode(name='A', host='239.67.52.72'),
    StorageNode(name='B', host='137.70.131.229'),
    StorageNode(name='C', host='98.5.87.182'),
    StorageNode(name='D', host='11.225.158.95'),
    StorageNode(name='E', host='203.187.116.210'),
]
