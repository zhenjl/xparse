parse
=====

A collection of parsing utilities:

* xip - IP (only v4 for now) parser that expands a given IP string to all the IP addresses it represents.
* xtime - time parsing utility in Go. It exposes a set of time formats that it knows how to parse, and a single function `Parse()` to parse any time string.
* xtld - TLD parser that extracts the top-level-domain out from the given string. It uses the data set from
https://www.publicsuffix.org/list/effective_tld_names.dat.