from random import shuffle
from typing import List

from bloom_filter import BloomFilter

# no. of items to add
n = 20
# false positive probability
p = 0.05

bloom_filter = BloomFilter(n, p)
print(f"\nSize of bit array: {bloom_filter.size}")
print(f"False positive probability: {bloom_filter.fp_prob}")
print(f"Number of hash functions: {bloom_filter.hash_count}\n")

# words to be added
word_present: List[str] = [
    "abound",
    "abounds",
    "abundance",
    "abundant",
    "accessible",
    "bloom",
    "blossom",
    "bolster",
    "bonny",
    "bonus",
    "bonuses",
    "coherent",
    "cohesive",
    "colorful",
    "comely",
    "comfort",
    "gems",
    "generosity",
    "generous",
    "generously",
    "genial",
]

# word not added
word_absent: List[str] = [
    "bluff",
    "cheater",
    "hate",
    "war",
    "humanity",
    "racism",
    "hurt",
    "nuke",
    "gloomy",
    "facebook",
    "geeksforgeeks",
    "twitter",
]

# add all the words present in word_present array
for item in word_present:
    bloom_filter.add(item)

# randomly suffle words
shuffle(word_present)
shuffle(word_absent)

# add some random words into test_word array
test_words: List[str] = word_present[:10] + word_absent
shuffle(test_words)

print("word_present \n", word_present, "\n")
print("word_absent \n", word_absent, "\n")
print("test_words \n", test_words, "\n")

for word in test_words:
    if bloom_filter.check(word):
        if word in word_absent:
            print(f"{word} is false positive!!! ⁉️")
        else:
            print(f"{word} is probably present. 💯✅")
    else:
        print(f"{word} is definitely not present...!❌")
