AppNexus REST API Go Client
===========================
*Connect to the AppNexus REST API through your Golang application*

[![Build Status](https://travis-ci.org/tnako/appnexus.svg?branch=master)](https://travis-ci.org/tnako/appnexus)

**Important** This is very much beta and subject to breaking changes while under development. Do not use in production just yet!

Of the published [AppNexus API Services](https://wiki.appnexus.com/display/adnexusdocumentation/API+Services) listed, this package covers the following:

* Auth API Service [Docs](https://wiki.appnexus.com/display/adnexusdocumentation/Auth+API+Service)
* Member API Service [Docs](https://wiki.appnexus.com/display/adnexusdocumentation/Member+Service)
* Segment Service [Docs](https://wiki.appnexus.com/display/adnexusdocumentation/Segment+Service)
* Publisher Service [Docs](https://wiki.appnexus.com/display/api/Publisher+Service)
* Site Service [Docs](https://wiki.appnexus.com/display/api/Site+Service)
* Placement Service [Docs](https://wiki.appnexus.com/display/api/Placement+Service)
* Deal Service [Docs](https://wiki.appnexus.com/display/api/Deal+Service)

Support for the remaining services should follow - pull requests welcome :)

Getting started
--------------
```Bash
cd $GOPATH
go get github.com/tnako/appnexus
go get -u all
```

Be sure to run the tests with `go test` and have a look at the [examples directory](./examples/) for a usage demonstration.
