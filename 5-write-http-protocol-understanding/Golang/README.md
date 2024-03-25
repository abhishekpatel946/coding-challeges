# An Overview of HTTP

HTTP is a protocol for fetching resources such as HTML documents. It is the foundation of any data exchange on the Web and it is a client-server protocol, which means requests are initiated by the recipient, usually the Web browser. A complete document is reconstructed from the different sub-documents fetched, for instance, text, layout description, images, videos, scripts, and more.

![Alt text](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview/fetching_a_page.png)

For more information, Please visit [Mozila Docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview).

## HTTP Flow

When a client wants to communicate with a server, either the final server or an intermediate proxy, it performs the following steps:

1. Open a TCP connection: The TCP connection is used to send a request, or several, and receive an answer. The client may open a new connection, reuse an existing connection, or open several TCP connections to the servers.

2. Send an HTTP message: HTTP messages (before HTTP/2) are human-readable. With HTTP/2, these simple messages are encapsulated in frames, making them impossible to read directly, but the principle remains the same. For example:

```http
GET / HTTP/1.1
Host: developer.mozilla.org
Accept-Language: fr
```

3. Read the response sent by the server, such as:

```http
HTTP/1.1 200 OK
Date: Sat, 09 Oct 2010 14:28:02 GMT
Server: Apache
Last-Modified: Tue, 01 Dec 2009 20:18:22 GMT
ETag: "51142bc1-7449-479b075b2891b"
Accept-Ranges: bytes
Content-Length: 29769
Content-Type: text/html

<!DOCTYPE html>… (here come the 29769 bytes of the requested web page)
```

4. Close or reuse the connection for further requests.

If HTTP pipelining is activated, several requests can be sent without waiting for the first response to be fully received. HTTP pipelining has proven difficult to implement in existing networks, where old pieces of software coexist with modern versions. HTTP pipelining has been superseded in HTTP/2 with more robust multiplexing requests within a frame.

## HTTP Messages

HTTP messages, as defined in HTTP/1.1 and earlier, are human-readable. In HTTP/2, these messages are embedded into a binary structure, a frame, allowing optimizations like compression of headers and multiplexing. Even if only part of the original HTTP message is sent in this version of HTTP, the semantics of each message is unchanged and the client reconstitutes (virtually) the original HTTP/1.1 request. It is therefore useful to comprehend HTTP/2 messages in the HTTP/1.1 format.

There are two types of HTTP messages, requests and responses, each with its own format.

### Requests

An example HTTP request:
![Alt text](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview/http_request.png)

Requests consist of the following elements:

- An HTTP method, usually a verb like GET, POST, or a noun like OPTIONS or HEAD that defines the operation the client wants to perform. Typically, a client wants to fetch a resource (using GET) or post the value of an HTML form (using POST), though more operations may be needed in other cases.
- The path of the resource to fetch; the URL of the resource stripped from elements that are obvious from the context, for example without the protocol (<http://>), the domain (here, developer.mozilla.org), or the TCP port (here, 80).
- The version of the HTTP protocol.
- Optional headers that convey additional information for the servers.
- A body, for some methods like POST, similar to those in responses, which contain the resource sent.

### Requests

An example response:
![Alt text](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview/http_response.png)

Responses consist of the following elements:

- The version of the HTTP protocol they follow.
- A status code, indicating if the request was successful or not, and why.
- A status message, a non-authoritative short description of the status code.
- HTTP headers, like those for requests.
- Optionally, a body containing the fetched resource.

## Few Examples with Raw HTTP Client using NetCat/NC

- ### Getting the data '/foo' get endpoint with proper headers
<style>
    table {
        border-collapse: collapse;
        width: 100%;
    }
    th, td {
        border: 1px solid black;
        padding: 8px;
        vertical-align: top;
    }
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        font-family: "Courier New", Courier, monospace;
        margin: 0;
    }
</style>

<table>
<tr>
<th style="width: 50%;">Ternimal 1</th>
<th style="width: 50%;">Ternimal 3</th>
</tr>
<tr>
<td>
<pre>

```sh
header -> Accept-Encoding value -> [gzip, deflate, br]
header -> Accept-Language value -> [en-GB,en;q=0.8]
header -> User-Agent value -> [Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36]
header -> Port value -> [1234]
header -> Accept value -> [*/*]
```

</pre>
</td>
<td>

```http
nc localhost 1234
GET /foo HTTP/1.1
Host: localhost
Port: 1234
Accept: */*   
Accept-Encoding: gzip, deflate, br
Accept-Language: en-GB,en;q=0.8
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36

HTTP/1.1 200 OK
Date: Mon, 25 Mar 2024 11:01:58 GMT
Content-Length: 3
Content-Type: text/plain; charset=utf-8

bar
```

</td>
</tr>
</table>

- ### Getting the data '/foo' get endpoint with invalid headers
<style>
    table {
        border-collapse: collapse;
        width: 100%;
    }
    th, td {
        border: 1px solid black;
        padding: 8px;
        vertical-align: top;
    }
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        font-family: "Courier New", Courier, monospace;
        margin: 0;
    }
</style>

<table>
<tr>
<th style="width: 50%;">Ternimal 1</th>
<th style="width: 50%;">Ternimal 3</th>
</tr>
<tr>
<td>
<pre>

```http
nc localhost 1234
GET /foo http/1.1
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
Connection: close
```

</pre>
</td>
<td>

```http
nc localhost 1234
GET /foo HTTP/2.0
Host: localhost

HTTP/1.1 505 HTTP Version Not Supported: unsupported protocol version
Content-Type: text/plain; charset=utf-8
Connection: close

505 HTTP Version Not Supported: unsupported protocol version% 
```

</td>
</tr>
</table>

- ### Writting the data '/login' get endpoint with valid headers but no body
<style>
    table {
        border-collapse: collapse;
        width: 100%;
    }
    th, td {
        border: 1px solid black;
        padding: 8px;
        vertical-align: top;
    }
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        font-family: "Courier New", Courier, monospace;
        margin: 0;
    }
</style>

<table>
<tr>
<th style="width: 50%;">Ternimal 1</th>
<th style="width: 50%;">Ternimal 3</th>
</tr>
<tr>
<td>
<pre>

```sh
header -> Port value -> [1234]
header -> Accept value -> [*/*]
header -> Accept-Encoding value -> [gzip, deflate, br]
header -> Accept-Language value -> [en-GB,en;q=0.8]
header -> User-Agent value -> [Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36]
body ==> 

```

</pre>
</td>
<td>

```http
nc localhost 1234
POST /login HTTP/1.1
Host: localhost
Port: 1234
Accept: */*   
Accept-Encoding: gzip, deflate, br
Accept-Language: en-GB,en;q=0.8
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36

HTTP/1.1 200 OK
Date: Mon, 25 Mar 2024 11:03:25 GMT
Content-Length: 22
Content-Type: text/plain; charset=utf-8

login successful...!!!
```

</td>
</tr>
</table>

- ### Writting the data '/login' get endpoint with valid headers and body
<style>
    table {
        border-collapse: collapse;
        width: 100%;
    }
    th, td {
        border: 1px solid black;
        padding: 8px;
        vertical-align: top;
    }
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        font-family: "Courier New", Courier, monospace;
        margin: 0;
    }
</style>

<table>
<tr>
<th style="width: 50%;">Ternimal 1</th>
<th style="width: 50%;">Ternimal 3</th>
</tr>
<tr>
<td>
<pre>

```sh
Server started at =: http://localhost:1234
header -> Content-Length value -> [28]
body ==> user=abhishek&password=pass
```

</pre>
</td>
<td>

```http
nc localhost 1234
POST /login HTTP/1.1
Host: localhost
Content-Length: 28

user=abhishek&password=pass
HTTP/1.1 200 OK
Date: Mon, 25 Mar 2024 11:21:34 GMT
Content-Length: 22
Content-Type: text/plain; charset=utf-8

login successful...!!!
```

</td>
</tr>
</table>

- ### Writting the data '/login' get endpoint with invalid headers

<style>
    table {
        border-collapse: collapse;
        width: 100%;
    }
    th, td {
        border: 1px solid black;
        padding: 8px;
        vertical-align: top;
    }
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        font-family: "Courier New", Courier, monospace;
        margin: 0;
    }
</style>

<table>
<tr>
<th style="width: 50%;">Ternimal 1</th>
<th style="width: 50%;">Ternimal 2</th>
</tr>
<tr>
<td>
<pre>

```http
nc localhost 1234
POST /loginn HTTP/1.1
Host: localhost

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Mon, 25 Mar 2024 10:57:45 GMT
Content-Length: 19

404 page not found
```

</pre>
</td>
<td>

```http
nc localhost 1234
POST /login HTTP/1.1

HTTP/1.1 400 Bad Request: missing required Host header
Content-Type: text/plain; charset=utf-8
Connection: close

400 Bad Request: missing required Host header%  
```

</td>
</tr>
</table>
