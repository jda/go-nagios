go-nagios
=========
go-nagios is a library for writing Nagios plugins in the Go programming
language.  go-nagios is free software, licensed under the liberal ISC license.
The API is currently liable to change, particularly with regards to performance
data.

Go has several properties that make it a particularly promising language for
tasks like writing Nagios plugins. Its niche seems to be somewhere between C
and Python, offering both efficient compilation as well as a quick and easy
development cycle.

Installation
------------
Installation is easy. After installing the Go distribution, simply run
`go get github.com/laziac/go-nagios`.

Documentation
-------------
You can view the documentation online [here][doc] or, once you've installed it,
locally with `godoc github.com/laziac/go-nagios`.

TODO
----
* figure out easier API for performance data
* unit tests
* example plugins

[doc]: http://go.pkgdoc.org/github.com/laziac/go-nagios
