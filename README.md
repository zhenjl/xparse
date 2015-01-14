timex
=====

Timex is a time parsing utility for Go. It exposes a set of time formats that it knows how to parse, and a single function `Parse()` to parse any time string.

```
import (
	"fmt"
	"time"
	"github.com/dataence/timex"
)

func main() {
	t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	t2, err := timex.Parse("2006-01-02T15:04:05+07:00")
	if err != nil {
		fmt.Println(err)
	} else if t1.UnixNano() != t2.UnixNano() {
		fmt.Println("%d != %d", t1.UnixNano(), t2.UnixNano())
	} else {
		fmt.Println(t2)
	}
}
```