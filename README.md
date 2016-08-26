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
  "port": "7896"
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
If you dont want https omit it
```
{
  ssl: {
    "private": "test.key",
    "cert": "test.cert"
  }
}
```

**endpoints**<br>
override default port
```
{
  "endpoints": [
    "paths": "",
    "methods": ["", ""],
    "headers": [{
      "Content-Type": "application/json"
    }],
    "body": "response"
  ]
}
```

## TODO
- add gzip support https://gist.github.com/the42/1956518
- Add proxy feature to forward requests to another REST API
- Unit Tests
- CI system (concourse)
- add more useful --help flag output
- verbose flag for debugging
- Compile Linux/Windows binary
- Scoop (for windows install)
- Terminal Output match python SimpleHTTPServer
- Use HTML templates to make file server look more like python SimpleHTTPServer
