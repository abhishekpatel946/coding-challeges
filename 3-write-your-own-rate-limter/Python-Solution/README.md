
# Rate Limiter

Rate limiters are a key part of building an API or large scale distributed system, they help when we wish to throttle traffic based on the user. They allow you to ensure that one or more bad actors can’t accidentally or deliberately overload the service.

### A rate limiting strategy can make your API more reliable, when

- A user is responsible for a spike in traffic, and you need to stay up for everyone else.
- A user is accidentally sending you a lot of requests.
- A bad actor is trying to overwhelm your servers.
- A user is sending you a lot of lower-priority requests, and you want to make sure that it doesn’t affect your high-priority traffic.
Your service is degraded, and as a result you can’t handle your regular traffic load and need to drop low-priority requests.

## There are 6 common approaches to rate limiting

1. Token bucket - tokens are added to a ‘bucket’ at a fixed rate. The bucket has a fixed capacity. When a request is made it will only be accepted if there are enough tokens in the bucket. Tokens are removed from the bucket when a request is accepted.
2. Leaky bucket (as a meter) - This is equivalent to the token bucket, but implemented in a different way - a mirror image.
3. Leaky bucket (as a queue) - The bucket behaves like a FIFO queue with a limited capacity, if there is space in the bucket the request is accepted.
4. Fixed window counter - record the number of requests from a sender occurring in the rate limit’s fixed time interval, if it exceeds the limit the request is rejected.
5. Sliding window log - Store a timestamp for each request in a sorted set, when the size of the set trimmed to the window exceeds the limit, requests are rejected.
6. Sliding window counter - similar to the fixed window, but each request is timestamped and the window slides.

## rate-limiter configuration

sample_configuration.yml

```yml
name: Rate Limiter Configuration
version: 1.0

ratelimiter:
    - name: sample-rate-limiter-configuration
    algorithm:
        token_bucket: 
            active: enable,
            capacity: 10,        # no of requests per seconds
            refresh_interval: 60 # in seconds
        leaky_bucket_mirror:
            active: disable,
        leaky_bucket_queue:
            active: disable,
        sliding_window_log:
            active: disable,
        sliding_window_counter:
            active: disable,

redis:
    - name: sample-redis-configuration
    environment:
        hostname: localhost
        port: 6379
        interval: 1000 * 60
        label: development
    
```

## Installation

Install the repository

```bash
git clone https://github.com/abhishekpatel946/coding-challeges
```

Go to the following directory

```bash
cd /3-write-your-own-rate-limter/Python-Solution
```

Install the respective packages

```bash
pip3 install -r requirements.txt
```

Launch the redis server for local testing, if redis-server not available please install it first [download & installation](https://redis.io/docs/install/install-redis/)

```bash
redis-server
```

Launch the server

```bash
uvicorn main:app --reload
```

Open the swagger documentation on locally to check health information

```bash
http://127.0.0.1:8000/
```

Open the rate limiting endpoint

```bash
http://127.0.0.1:8000/authApp/docs
```

## Report

<h4>Performance Report</h4>
<object data="https://github.com/abhishekpatel946/coding-challeges/blob/main/3-write-your-own-rate-limter/Python-Solution/output/Rate-Limiter-performance-report.pdf" type="application/pdf" width="700px" height="700px">
    <embed src="https://github.com/abhishekpatel946/coding-challeges/blob/main/3-write-your-own-rate-limter/Python-Solution/output/Rate-Limiter-performance-report.pdf">
        <p>This browser does not support PDFs. Please download the PDF to view it: <a href="https://github.com/abhishekpatel946/coding-challeges/blob/main/3-write-your-own-rate-limter/Python-Solution/output/Rate-Limiter-performance-report.pdf">Download PDF</a>.</p>
    </embed>
</object>

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
