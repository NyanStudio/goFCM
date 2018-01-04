# goFCM : Firebase Cloud Messaging Library

## Usage

```
go get github.com/NyanStudio/goFCM
```

## Example

### Send notification to a registration token

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
	cm := new(fcm.Client)

	cm.SetServerKey(serverKey)

	cm.SetTo("REGISTRATION_TOKEN")

	cm.SetNotification(fcm.NotificationPayload{Title: "title", Body: "body"})

	rm, err := cm.SendMessage()
	if err != nil {
		fmt.Sprintln("ERR: %v", err)
	}

	fmt.Sprintln("RM: %v", rm)
}

```

### Send notification to registration tokens

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
	cm := new(fcm.Client)

	cm.SetServerKey(serverKey)

	cm.SetRegistrationIds([]string{"REGISTRATION_TOKEN1", "REGISTRATION_TOKEN2"})

	cm.SetNotification(fcm.NotificationPayload{Title: "title", Body: "body"})

	rm, err := cm.SendMessage()
	if err != nil {
		fmt.Sprintln("ERR: %v", err)
	}

	fmt.Sprintln("RM: %v", rm)
}

```

### Send notification to a topic

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
	cm := new(fcm.Client)

	cm.SetServerKey(serverKey)

	cm.SetTo("/topics/YOUR_TOPIC")

	cm.SetNotification(fcm.NotificationPayload{Title: "title", Body: "body"})

	rm, err := cm.SendMessage()
	if err != nil {
		fmt.Sprintln("ERR: %v", err)
	}

	fmt.Sprintln("RM: %v", rm)
}

```
