eventsender
===========

A go library to stream events to an HTTP endpoint.

Example:

```golang
package main

import (
  "fmt"
  "time"
  "github.com/trusch/eventsender"
)

func main(){
  sender := eventsender.New("POST", "http://localhost:8080")  
  defer sender.Close()
  for i:=0; i<10; i++ {
    sender.SendEvent(i)
    time.Sleep(1 * time.Second)
    sender.SendEvent(fmt.Sprintf("%v", i))
    time.Sleep(1 * time.Second)
    sender.SendEvent(map[string]interface{}{"foo": 123})
  }
}
```
