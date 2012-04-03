splint
======

`splint` is a little Go application to analyze Go source files.  It finds any functions that are
too long or have too many parameters or results.

These are typical signs that a function is doing too much.  We find `splint` to be a helpful tool
for detecting potential problem areas in our code, areas that should be refactored.  We tolerate long
functions and functions with long parameter/result lists when they are needed, but generally try to
keep them short.

Installation
------------

Use `go install`:

    go install github.com/stathat/splint

Usage
-----

Examples available at [www.stathat.com/src/splint](http://www.stathat.com/src/splint).

Contact us
----------

We'd love to hear from you if you are using `splint`!  We're on twitter: [@stat_hat](http://twitter.com/stat_hat) or [contact us here](http://www.stathat.com/docs/contact).

About
-----

Written by Patrick Crosby at [StatHat](http://www.stathat.com).  Twitter:  [@stat_hat](http://twitter.com/stat_hat)
