# calendar-service [![Circle CI](https://circleci.com/gh/danesparza/calendar-service.svg?style=svg)](https://circleci.com/gh/danesparza/calendar-service)
A microservice to return Google calendar events in JSON format.  

[Download the latest release for your platform here](https://github.com/danesparza/calendar-service/releases/latest)

### Starting and testing the service
To start the service, just run `calendar-service`.  

If you need help, just run `calendar-service --help`.

There are a few command line parameters available:

Parameter       | Description
----------      | -----------
authEmail       | The OAuth 2.0 [service account email address](https://developers.google.com/console/help/new/#serviceaccounts), as listed in your project in the [Google Developer's console](https://console.developers.google.com)
authSubject     | The email address to impersonate.  The request will be made on behalf of this account, so you need ownership of this user account
keyFile         | The location of the PEM encoded keyfile to authenticate your service account
port            | The port the service listens on.  
allowedOrigins  | comma seperated list of [CORS](http://en.wikipedia.org/wiki/Cross-origin_resource_sharing) origins to allow.  In order to access the service directly from a javascript application, you'll need to specify the origin you'll be running the javascript site on.  For example: http://www.myjavascriptapplication.com

Once the service is up and running, you can connect to it using
`http://yourdomain.com:3000/calendar/calendarid` where `calendarid` is the Google calendar id you'd like information for.  

Example: `http://yourdomain:3000/calendar/mg2l41ag8ua062trmktgdq6v90@group.calendar.google.com`

To test your service quickly, you can use the [Postman Google Chrome Extension](https://chrome.google.com/webstore/detail/postman-rest-client/fdmmgilgnpjigdojojpjoooidkmcomcm?hl=en) to call the service and see the JSON return format.

Calendar information will be returned as a [JSON payload outlined on the Google Calendar developer website](https://developers.google.com/google-apps/calendar/v3/reference/events).

### Building
### Prerequisites
*To build, make sure you have the latest version of [Go](http://golang.org/) installed.  If you've never used Go before, it's a quick install and [there are installers for multiple platforms](http://golang.org/doc/install), including Windows, Linux and OSX.*

### Quick Start

Run the following commands to get latest and build.

```bash
go get github.com/danesparza/calendar-service
go build
```
