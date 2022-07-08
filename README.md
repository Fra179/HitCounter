# HitCounter
HitCounter is a simple and minimal service useful to count hits on a website. 
You can see it as a personal remake of [brentvollebregt/hit-counter](https://github.com/brentvollebregt/hit-counter)'s service. 

This project is still a WIP and as such it doesn't have any kind of request's rate limiter and it's only equipped 
with a simple system to prevent users from constantly updating the page to increase the views count.  

## Endpoints
### `GET /count?url=<url>`
Increases the hits count for `<url>` and returns it as a number in the response body. 

### `GET /get?url=<url>`
Returns the hits count for `<url>` without increasing it. 

### `GET /status`
Simply returns `Ok` if the service is still alive. If it's not your request will probably fail :P

## Usage
### Getting the raw count

```javascript
fetch('https://hitcounter.francescodb.me/count?url=example.com', {
    credentials: 'include'
})
    .then(res => res.text())
    .then(count => console.log('Count: ' + count))
```

#### Using XMLHttpRequest

```javascript
let xmlHttp = new XMLHttpRequest();
xmlHttp.withCredentials = true;
xmlHttp.onload = function() {
    console.log('Count: ' + this.responseText);
};
xmlHttp.open('GET', 'https://hitcounter.francescodb.me/count?url=example.com', true);
xmlHttp.send(null);
```

#### Using Ajax

```javascript
let targetUrl = "example.com";
$.ajax('https://hitcounter.francescodb.me/count', {
    data: { url: targetUrl },
    xhrFields: { withCredentials: true }
}).then(count => console.log('Count: ' + count));
```

## Deployment 
To deploy this service you can simply use docker-compose
```shell
$ git clone https://github.com/Fra179/HitCounter.git
$ cd HitCounter
$ docker-compose up -d
```

## TODO:
- [x] Count hits of a URL
- [ ] Return the number of hits of a URL
- [ ] Return an SVG with the number of hits
- [ ] Remove the use of `<url>` as a query parameter by parsing the source from the `Origin` header
- [ ] Statistics
- [ ] More...