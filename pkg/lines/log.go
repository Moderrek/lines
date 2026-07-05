package lines

import (
	"fmt"
	"os"
)

func (c *Counter) logVerbosef(format string, v ...any) {
	if c.Config.Verbose {
		fmt.Fprintf(os.Stderr, "info: "+format+"\n", v...)
	}
}
