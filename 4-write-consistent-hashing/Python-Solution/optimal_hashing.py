import hashlib
from bisect import bisect, bisect_left, bisect_right

from plot import plot_hashring
from StorageNode import StorageNode, storage_nodes


def hash_fn(key: str, total_slots: int) -> int:
    '''
        hash_fn creates an integer equivalent of a SHA256 hash and
        takes a modulo with the total number of slots in hash ring.
    '''
    _hash = hashlib.sha256()

    # converting data into bytes and passing it to hash function
    _hash.update(bytes(key.encode('utf-8')))

    # converting the HEX digest into equivalent integer value
    return int(_hash.hexdigest(), 16) % total_slots


class ConsistentHash:
    '''
        ConsistentHash represents an array based implementation of
        consistent hashing algorithm
    '''

    def __init__(self):
        self._keys = list()     # indices taken up in the ring
        # nodes present in the ring, nodes[i] is present at index keys[i]
        self.nodes = list()
        self.total_slots = 50   # total slots in the ring

    def add_node(self, node: StorageNode) -> int:
        '''
            add_node function adds a new node in the system and returns
            the key from hash space where it was placed
        '''

        # handling error when hash ring is full
        if len(self._keys) == self.total_slots:
            raise Exception('Hash ring is full')

        key = hash_fn(node.host, self.total_slots)

        # find the index where the key should be inserted in the keys array
        # this will be the index where the StorageNode will be added in the
        # nodes array
        index = bisect(self._keys, key)

        # if we have already seen the key i.e. node already exists for the
        # given key,
        # we raise CollisionError
        if index > 0 and self._keys[index - 1] == key:
            raise Exception('Collision detected')

        # insert the node_id and the key at the same `index` position.
        # this insertion will keep nodes and keys sorted w.r.t. keys
        self.nodes.insert(index, node)
        self._keys.insert(index, key)

        return key

    def remove_node(self, node: StorageNode) -> int:
        '''
            remove_node removes the node and returns the key from the hash ring
            on which the node was placed.
        '''

        # handling error when space is empty
        if len(self._keys) == 0:
            raise Exception('Hash ring is empty')

        key = hash_fn(node.host, self.total_slots)

        # we find the index where the key would reside in the keys
        index = bisect_left(self._keys, key)

        # if key doesn't exist in the hash ring array then raise an exception
        if index >= len(self._keys) and self._keys[index] != key:
            raise Exception('Node does not exist in the hash ring')

        # now that all sanity checks have passed, we poping the keys and nodes
        # at the index and thus removing presennce of the node from the hash ring
        self._keys.pop(index)
        self._nodes.pop(index)

    def assign(self, item: str) -> str:
        '''
            given an item, the function returns the node it is associated with
        '''
        key = hash_fn(item, self.total_slots)

        # we find the first node to the right of this key if bisect_right
        # returns index which is out of bounds then we circle back to the
        # first in the array in circular fashion
        index = bisect_right(self._keys, key) % len(self._keys)

        # return the node present at the index
        return self.nodes[index]

    def plot(self, item: str = None, node: StorageNode = None) -> None:
        plot_hashring(
            self.total_slots,
            self._keys,
            self.nodes,
            item_key=hash_fn(item, self.total_slots) if item else None,
            node_key=hash_fn(node.host, self.total_slots) if node else None
        )


if __name__ == '__main__':
    ch = ConsistentHash()

    for node in storage_nodes:
        ch.add_node(node)
    ch.plot()

    files = ['f1.txt', 'f2.txt', 'f3.txt', 'f4.txt', 'f5.txt']
    for file in files:
        print(
            f"file {file} (shown in green) resides on node {ch.assign(file).name} (shown in red)")
        ch.plot(file, ch.assign(file))
