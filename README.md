# goFCM : Firebase Cloud Messaging Library

## Usage

```
go get github.com/NyanStudio/goFCM
```

## Example

### Send Notification

```go

package main

import (
	"fmt"
	"github.com/NyanStudio/goFCM"
)

const (
	serverKey = "YOUR_SERVER_KEY"
)

func main() {
	cm := new(cloudMessaging.Client)

	cm.SetServerKey(serverKey)

	cm.SetTo("REGISTRATION_TOKEN")

	cm.SetNotification("title", "body", "", "", "", "", "", "", "", "", "", "", "", "")

	rm, err := cm.SendMessage()
	if err != nil {
		fmt.Sprintln("ERR: %v", err)
	}

	fmt.Sprintln("RM: %v", rm)
}

```
