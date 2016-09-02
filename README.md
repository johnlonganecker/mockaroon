# Mockaroon
Mockaroon: A Simple HTTP Server for Web App development

Made to act and behave exactly like `python -m SimpleHTTPServer`

## Setup

```
brew tap johnlonganecker/mockaroon
brew install mockaroon
```

**Upgrade to latest release**
```
brew upgrade mockaroon
```

## How to Use

Serve static content
```
> mockaroon
Serving HTTP on 0.0.0.0 port 8000 ...
```

Serve static content - specify port
```
> mockaroon 8080
Serving HTTP on 0.0.0.0 port 8080 ...
```

Load configuration from a config file
```
> mockaroon --config config.json
Serving HTTPS on 0.0.0.0 port 8080 ...
Serving static files
adding route /path1
adding route /path2
adding route /path3
```

### How is this different from Python SimpleHTTPServer?

#### Config Files
Currently JSON and YAML config files are supported
```
> mockaroon --config=configfile.(json|yml) [optional port override]
```

#### Fields
**port** `default 8000`<br>
override default port
```
{
  "port": 7896
}
```

**staticFiles** `default true`<br>
override default of hosting static files
```
{
  "staticFiles": false
}
```

**ssl**<br>

Friendly reminder - to create SSL cert and private key
```
ssh-keygen -f private.key
openssl req -x509 -nodes -days 3065 -newkey rsa:2048 -keyout private.key -out cert.crt
```

If you dont want https omit it
```
{
  ssl: {
    "private": "private.key",
    "cert": "cert.cert"
  }
}
```

**endpoints**<br>
Latency allows you to randomly pick a delay between min and max milliseconds
```
{
  "endpoints": [
    {
      "paths": "",
      "methods": ["", ""],
      "headers": [{
        "Content-Type": "application/json"
      }],
      "body": "response"
      "latency": {
        "min": 100,
        "max": 1000
      }
    }
  ]
}
```

**proxies**<br>
```
{
  "proxies": [
    {
      "paths": ["/users/{.*}"],
      "destination": "https://somehost:1337"
    }
  ]
}
```

## TODO
- enable CORS
- better readme/docs
- add host address to bind too
- Unit Tests
- CI system (concourse)
- Add host routing feature
- Add more useful --help flag output
- verbose flag for debugging
- Compile Linux/Windows binary
- Scoop (for windows install)
- Terminal Output match python SimpleHTTPServer
- Use HTML templates to make file server look more like python SimpleHTTPServer
- randomize response
- Dependency management (IE godeps)
