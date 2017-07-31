# Cocobolo

Cocobolo is a proxy service to make outbound requests. It uses gRPC to form
a bi-directional stream, all requests and responses are sent on this stream.

Some key features of this service are:

1. Queue multiple outbound request.
2. Ability to exponentially backoff if the external server is unavailable.
3. Reply back on gRPC with the response from the external request.


## Payload

```
    /request
```


- request_id
- endpoint
- method
- body
- headers
- backoff_time


```
/response

```


- request_id
- status_code
- response
- headers
- backoff_time
