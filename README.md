<a href="https://jgltechnologies.com/discord">
<img src="https://discord.com/api/guilds/844418702430175272/embed.png">
</a>

# SimpleFiles

SimpleFiles is a go library to make reading and writing to file easier. SimpleFiles uses a RWMutex to prevent errors when using a file in multiple gorotines.

<br>

```
go get github.com/JGLTechnologies/SimpleFiles
```

<br>

Reading JSON from a file

```go
package main

import (
	"fmt"
	"github.com/JGLTechnologies/SimpleFiles"
)

func main() {
	var j map[string]interface{}
	// If the file does not exist it will be created.
	f, err := SimpleFiles.New("test.json")
	if err != nil {
		panic(err)
	} else {
		err := f.ReadJSON(&j)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(j["name"])
		}
	}
}
```

<br>

Writing JSON to a file

```go
package main

import (
	"github.com/JGLTechnologies/SimpleFiles"
)

func main() {
	j := map[string]interface{}{"name": "Joe", "age": 47}
	// If the file does not exist it will be created.
	f, err := SimpleFiles.New("test.json")
	if err != nil {
		panic(err)
	} else {
		err := f.WriteJSON(j)
		if err != nil {
			panic(err)
		}
	}
}
```