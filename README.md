# SpaceAPI-Server
SpaceAPI Server written in Go

This server can return a valid SpaceAPI-string in version 13 as specified 
[here](https://spacedirectory.org/pages/docs.html]) and provides further API endpoints for modification of the returned values.

## Features

### Implemented

*  Return of valid SpaceAPI strings
*  Modification of SpaceAPI
*  Persistence using a flatfile
*  Simple token authentication of modification-requests
*  Dockerfile

### Planned

*  Support for the whole SpaceAPI (with all specified fields) including modification

## Running

Start the application and send the payload via PUT 

payload.template.json
```
{
    "API":   "0.13",
    "Space": "vspace.one",
    "Logo":  "https://vspace.one/pic/logo_vspaceone.svg",
    "URL":   "https://vspace.one",
    "location": {
        "Address": "Wilhelm-Binder-Str. 19, 78048 VS-Villingen, Germany",
        "Lat":     48.065003,
        "Lon":     8.456495
    },
    "contact": {
        "Phone":   "+49 221 596196638",
        "Email":   "info@vspace.one",
        "Twitter": "@vspace.one"
    },
    "IssueReportChannels": [
        "email",
        "twitter"
    ]
}
```

* `docker run --name spaceapi-server -e TOKEN=yoursecrettoken -v /srv/spaceapi-server/:/go/src/github.com/chaosdorf/SpaceAPI-Server/data chaosdorf/spaceapi-server`

## API

### Getting SpaceAPI string

*GET on /spaceapi*  

Returns the whole SpaceAPI string

### Setting SpaceAPI values

*PUT on /spaceapi*

Makes it possible to send data similar to the SpaceAPI string to set 

**Note** that setting these values is only possible if the right token is specified in Header as `Authorization: Token $TOKEN`. When specifying a wrong token or none the server will respond with status 401.
