## Consitent Hashing algorithm

### Traditional Hashing algorithm

```sh
python traditional_hashing.py
```

```bash
$ python3 traditional_hashing.py

  file f1.txt resides on node E
  file f2.txt resides on node A
  file f3.txt resides on node B
  file f4.txt resides on node C
  file f5.txt resides on node D
  file f6.txt resides on node E
  file f7.txt resides on node A
  file f8.txt resides on node B
  file f9.txt resides on node C
  
  {'E': 2, 'A': 2, 'B': 2, 'C': 2, 'D': 1}
  
  file f1.txt resides on node E
  file f2.txt resides on node A
  file f3.txt resides on node B
  file f4.txt resides on node C
  file f5.txt resides on node D
  
  {'E': 1, 'A': 1, 'B': 1, 'C': 1, 'D': 1}
```

### Optimal Consistent Hashing algorithm

```sh
$ python3 optimal_hashing.py 

  file f1.txt (shown in green) resides on node E (shown in red)
  file f2.txt (shown in green) resides on node B (shown in red)
  file f3.txt (shown in green) resides on node B (shown in red)
  file f4.txt (shown in green) resides on node C (shown in red)
  file f5.txt (shown in green) resides on node E (shown in red)
```
