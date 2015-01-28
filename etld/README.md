etld [DEPRECATED]
====

etld was an attempt for a TLD parser that extracts the top-level-domain out from the given string. It uses the data set from
https://www.publicsuffix.org/list/effective_tld_names.dat. Unfortunately it's still not 100% and there's already a working
implementation at golang.org/x/net/publicsuffix. So I recommend using that instead.
