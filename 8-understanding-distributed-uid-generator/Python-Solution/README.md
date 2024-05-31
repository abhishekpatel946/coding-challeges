## Distributed Unique Id Generator

Generates a unique identifier is 64 bits long. [Universally Unique Identifier](https://en.wikipedia.org/wiki/Universally_unique_identifier) contains the similiar functionalities in our implementation as well, will used this 64 bit uid in distributed databases and across the regions where the identifier is unique and ability to generate over 10,000 unique identifiers per second.

### Features or Requirements

1. IDs must be unique.
2. IDs are numerical values only.
3. IDs fit into 64-bit.
4. IDs are ordered by date and time.
5. Ability to generate over 10,000 unique identifiers per second.

### Approaches

Multiple options can be used to generate unique IDs in distributed systems. For example:

1. [Multi-master replication](https://arpitbhayani.me/blogs/multi-master-replication/)

> This approach uses the databases `auto_increment` feature to generate unique IDs. However,this strategy has some major drawbacks.

- Hard to scale with multiple data centers.
- IDs do not do up with time across distributed systems (multiple servers).
- It does not scale well when a server is added or removed.

2. Universally Unique Identifier [(UUID)](https://en.wikipedia.org/wiki/Universally_unique_identifier)

> Pros:

- Generating UUID is simple. No coordination between servers is needed so there will not be any  synchronization issues.
- The system is easy to scale because each server is responsible for generating IDs they consume. ID generator can easily scale with servers.

> Cons:

- IDs are 128 bits long, but our requirement is 64 bits.
- IDs do not go up with time.
- IDs could be non-numeric.

3. [Ticket Server](https://code.flickr.net/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/)

> Pros:

- Numeric IDs.
- It is easy to implement, and it workks for small to medium scale applications.

> Cons:

- Single point of failure (SPOF), single ticket server means if the ticket server goes down, all the system that depends on the ticket server will fail. To avoid this, we can set up multiple ticket servers. However, this will introduce new challenges such as data synchronization issues.

4. [Twitter Snowflake](https://blog.x.com/engineering/en_us/a/2010/announcing-snowflake)

> Twitter's unique ID generation system called "Twitter Snowflake" is inspiring and can satisfy our requirements. Divide and conquer is our friend. Instead of generating an ID directly we divide the and ID into different sections.

``` shell
| 1 bit |    4 bits   |     5 bits     |   5 bits    |   12 bits         |
|------------------------------------------------------------------------|
|   0   |  timestamp  |  datacenterID  |  machineID  |  sequence number  |
|------------------------------------------------------------------------|
                            64 bits ID   
```

- 1 bits (Sign Bit): It will always be 0. This is reserved for future uses. It can potentially be used to distinguish between signed and unsigned numbers.
- 41 bits (Timestamp): Milliseconds since the epoch or custom epoch.
- 5 bits (Datacenter ID): which gives us 2^5 = 32 datacenters.
- 5 bits  (Machine ID): which gives us 2^5 = 32 machines per datacenter.
- 12 bits (Sequence Number): For every ID generated that machine/processes, the sequence number is incremented by 1. The number is reset to 0 every millisecond.

### Installation

Install the repository

```bash
git clone https://github.com/abhishekpatel946/coding-challeges
```

Go to the following directory

```bash
cd /8-understanding-distributed-uid-generator/Python-Solution
```

Install the respective packages

```bash
pip3 install -r requirements.txt
```

Launch the server

```bash
python3 main.py
```

Open the swagger documentation on locally to check health information

```bash
http://127.0.0.1:8000/
```

Open the rate limiting endpoint

```bash
http://127.0.0.1:8000/authApp/docs
```

### Benchmarking

1. System Configuration

```bash
    machdep.cpu.cores_per_package: 8
    machdep.cpu.core_count: 8
    machdep.cpu.logical_per_package: 8
    machdep.cpu.thread_count: 8
    machdep.cpu.brand_string: Apple M1
```

2. Generate a single UID using "/generate" endpoint (kind of brute-force approach)

```bash
docker run --rm draftdev/rt http://127.0.0.1:8000/authApp/generate
```

```bash
          final_url:  http://127.0.0.1:8000/authApp/generate
      response_code:  000s
    time_namelookup:  0.002360s
       time_connect:  0.000000s
    time_appconnect:  0.000000s
   time_pretransfer:  0.000000s
      time_redirect:  0.000000s
 time_starttransfer:  0.000000s
                    ----------
         time_total:  0.004291s
```

3. Generate a single UID using "/optimized-generate" endpoint (kind of more efficient approach)

```bash
docker run --rm draftdev/rt http://127.0.0.1:8000/authApp/optimized-generate
```

```bash
          final_url:  http://127.0.0.1:8000/authApp/optimized-generate
      response_code:  000s
    time_namelookup:  0.002429s
       time_connect:  0.000000s
    time_appconnect:  0.000000s
   time_pretransfer:  0.000000s
      time_redirect:  0.000000s
 time_starttransfer:  0.000000s
                    ----------
         time_total:  0.004139s
```

4. Generate a single UID using "/optimized-generate?total_num_ids=10000" endpoint (generate the 10,000 UID internally)

```bash
docker run --rm draftdev/rt http://127.0.0.1:8000/authApp/test/optimized-generate?total_num_ids=10000
```

```bash
          final_url:  http://127.0.0.1:8000/authApp/test/optimized-generate?total_num_ids=10000
      response_code:  000s
    time_namelookup:  0.002531s
       time_connect:  0.000000s
    time_appconnect:  0.000000s
   time_pretransfer:  0.000000s
      time_redirect:  0.000000s
 time_starttransfer:  0.000000s
                    ----------
         time_total:  0.009991s
```

5. Test with unittest to call the "test/optimized-generate" endpoint with below parameters:

```bash
<!-- test_optimized_api -->

    run-unittest...
    total no of request 10000 with concurrency 1
    total duplicates are: 0 out of 10000
    Successfully generated unique IDs, Generated 10000 IDs in 43.305 seconds

```

### Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

### License

[MIT](https://choosealicense.com/licenses/mit/)
