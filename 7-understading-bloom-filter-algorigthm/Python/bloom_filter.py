import math
import mmh3
from bitarray import bitarray


class BloomFilter:
    """
    class for BloomFilter, using murmur3 hash functions
    """

    def __init__(self, items_count: int, fp_prob: float) -> None:
        """
        items_count: int
            Number of items expected to be stored in the filter
        fp_prob: float
            False positive probability in decimal
        """

        # false positive probability in decimal
        self.fp_prob = fp_prob

        # size of bit array to use
        self.size = self.get_size(items_count, fp_prob)

        # number of hash functions to use
        self.hash_count = self.get_hash_count(self.size, items_count)

        # bit array of givven size
        self.bit_array = bitarray(self.size)

        # initialize all bits as 0 (default)
        self.bit_array.setall(0)

    def add(self, item: str) -> None:
        """
        add an item in the filter

        Args:
            item (str): item to add
        """
        digests = list()
        for i in range(self.hash_count):
            # create digest for given item.
            # i work as seed to mmh3.hash fn with different seed,
            # diget created is different
            digest = mmh3.hash(item, i) % self.size
            digests.append(digest)

            # set the bit true in bit array
            self.bit_array[digest] = True

    def check(self, item: str) -> bool:
        """
        check for existence of an item in the filter

        Args:
            item (str): item to check for existence in the filter

        Returns:
            int: boolean
        """
        for i in range(self.hash_count):
            digest = mmh3.hash(item, i) % self.size
            if not self.bit_array[digest]:
                # if any of bit it False then, its not present in filter
                # else there is probability that it exist
                return False
        return True

    @classmethod
    def get_size(self, n: int, p: float) -> int:
        """
        return the size of bit_arary(m) to be used using following formula
        m = -(n * log(p)) / (log(2)^2)

        Args:
            n (int): number of items expected to be stored in filter
            p (float): flase positive probabiltiy in deicial

        Returns:
            int: size
        """
        return int(-(n * math.log(p)) / (math.log(2) ** 2))

    @classmethod
    def get_hash_count(self, m: int, n: int) -> int:
        """
        return the hash function(k) to be used using following formula
         k = (m/n) * log(2)

        Args:
            m (int): size of bit array
            n (int): number of items expected to be stored in filter

        Returns:
            int: hash count
        """
        return int((m / n) * math.log(2))
