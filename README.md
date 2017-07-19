# Cocobolo

Cocobolo is a proxy service to make outbound requests. It uses gRPC to form
a bi-directional stream, all requests and responses are sent on this stream.

Some key features of this service are:

1. Queue multiple outbound request.
2. Ability to exponentially backoff if the external server is unavailable.
3. Reply back on gRPC with the response from the external request.


## Payload

    /request


- request_id
- URL
- method
- request_body
- request_type
- backoff_time
- headers


```
/response

```


- request_id
- URL
- method
- status_code
- response_body
- backoff_time
- response_headers
- request_headers

