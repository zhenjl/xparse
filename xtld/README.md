xtld
====

xtld is a TLD parser that extracts the top-level-domain out from the given string. It uses the data set from
https://www.publicsuffix.org/list/effective_tld_names.dat. This package is inspired by https://github.com/jpillora/go-tld.

```
import "github.com/parse/xtld"

func main() {
	xtld.TLD("abc.com") // should return "com"
	xtld.TLD("中国")    // should return "中国"
	xtld.TLD("xxxxx")   // should return ""
}   
```