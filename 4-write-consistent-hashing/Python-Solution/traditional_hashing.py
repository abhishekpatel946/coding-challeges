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


def hash_fn(key):
    '''
        The function sums the bytes present in the `key` and then take a mod with 5 (total storage_nodes).
        This hash function thus generates output in the rage [0, 4]
    '''
    return sum(bytearray(key.encode('utf-8'))) % len(storage_nodes)


def upload(path):
    # we use the hash fn to get the index of the storage node that would hold the file
    index = hash_fn(path)

    # we get the StorageNode instance
    node = storage_nodes[index]

    # we put the file on the node and return
    return node.put_file(path)


def fetch(path):
    # we use the hash fn to get the index of the storage node that would hold the file
    index = hash_fn(path)

    # we get the StorageNode instance
    node = storage_nodes[index]

    # we fetch the file from the node and return
    return node.get_file(path)


if __name__ == '__main__':
    # Now we find where the 5 files "f1.txt", "f2.txt", "f3.txt", "f4.txt" and "f5.txt" are located on the storage nodes
    files = ['f1.txt', 'f2.txt', 'f3.txt', 'f4.txt',
             'f5.txt', 'f6.txt', 'f7.txt', 'f8.txt', 'f9.txt']
    node_instances = dict()

    for file in files:
        if not storage_nodes[hash_fn(file)].name in node_instances:
            node_instances[storage_nodes[hash_fn(file)].name] = 1
        else:
            node_instances[storage_nodes[hash_fn(file)].name] += 1

        print(
            f'file {file} resides on node {storage_nodes[hash_fn(file)].name}')

    print(f'\n {node_instances}')

    for key, value in node_instances.items():
        print(f'\n Number of files {value} stored on node {key}')

    # Now we find where the 5 files "f1.txt", "f2.txt", "f3.txt", "f4.txt" and "f5.txt" are located on the storage nodes
    files = ['f1.txt', 'f2.txt', 'f3.txt', 'f4.txt', 'f5.txt']
    node_instances = dict()

    for file in files:
        if not storage_nodes[hash_fn(file)].name in node_instances:
            node_instances[storage_nodes[hash_fn(file)].name] = 1
        else:
            node_instances[storage_nodes[hash_fn(file)].name] += 1

        print(
            f'file {file} resides on node {storage_nodes[hash_fn(file)].name}')

    print(f'\n {node_instances}')

    for key, value in node_instances.items():
        print(f'\n Number of files {value} stored on node {key}')
