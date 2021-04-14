# Client with a balancer to connect with  several servers at same time base on a custom RoundTripper.

### Test case to show that it works

1. TestServer_Run : to start three simple server with healthy check port
   
2. TestNewBalancerRoundTripper: use the client with balancer to do the request

```bash
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9882 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9882
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9881 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9881
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9881 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9881
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9881 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9881
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9880 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9880
2021/04/14 23:47:00 Response:  Here is localhost:9880
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9881 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9881
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9881 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9881
2021/04/14 23:47:00 balancer modfiy localhost:9880 => localhost:9882 : GET http://localhost:9880/name
2021/04/14 23:47:00 Response:  Here is localhost:9882
```


