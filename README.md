# GO-CMD

```go
import (
	cmd "github.com/yousysadmin/go-cmd"
	"log"
	"os"
)

func main() {
	command := "ping 127.0.0.1"
	logFile, _ := os.Create("./my.log")
	command := cmd.Cmd{
		Command: command,
		LogFile: logFile,
	}

	_, err := command.Run()
	if err != nil {
		log.Print(err)
	}
}
```
