package fcm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	// Firebase Cloud Messaging HTTP endpoint
	httpEndpointURL = "https://fcm.googleapis.com/fcm/send"
)

// HTTPMessage -> Downstream HTTP messages
type HTTPMessage struct {
	To                    string              `json:"to,omitempty"`                      // This parameter specifies the recipient of a message.
	RegistrationIds       []string            `json:"registration_ids,omitempty"`        // This parameter specifies the recipient of a multicast message, a message sent to more than one registration token.
	Condition             string              `json:"condition,omitempty"`               // This parameter specifies a logical expression of conditions that determine the message target.
	NotificationKey       string              `json:"notification_key,omitempty"`        // This parameter is deprecated. Instead, use to to specify message recipients.
	CollapseKey           string              `json:"collapse_key,omitempty"`            // This parameter identifies a group of messages (e.g., with collapse_key: "Updates Available") that can be collapsed, so that only the last message gets sent when delivery can be resumed. This is intended to avoid sending too many of the same messages when the device comes back online or becomes active.
	Priority              string              `json:"priority,omitempty"`                // Sets the priority of the message. Valid values are "normal" and "high." On iOS, these correspond to APNs priorities 5 and 10.
	ContentAvailable      bool                `json:"content_available,omitempty"`       // On iOS, use this field to represent content-available in the APNs payload. When a notification or message is sent and this is set to true, an inactive client app is awoken, and the message is sent through APNs as a silent notification and not through the FCM connection server. Note that silent notifications in APNs are not guaranteed to be delivered, and can depend on factors such as the user turning on Low Power Mode, force quitting the app, etc. On Android, data messages wake the app by default. On Chrome, currently not supported.
	MutableContent        bool                `json:"mutable_content,omitempty"`         // Currently for iOS 10+ devices only. On iOS, use this field to represent mutable-content in the APNs payload. When a notification is sent and this is set to true, the content of the notification can be modified before it is displayed, using a Notification Service app extension.
	DelayWhileIdle        bool                `json:"delay_while_idle,omitempty"`        // This parameter is deprecated.
	TimeToLive            int                 `json:"time_to_live,omitempty"`            // This parameter specifies how long (in seconds) the message should be kept in FCM storage if the device is offline. The maximum time to live supported is 4 weeks, and the default value is 4 weeks.
	RestrictedPackageName string              `json:"restricted_package_name,omitempty"` // This parameter specifies the package name of the application where the registration tokens must match in order to receive the message.
	DryRun                bool                `json:"dry_run,omitempty"`                 // This parameter, when set to true, allows developers to test a request without actually sending a message.
	Data                  interface{}         `json:"data,omitempty"`                    // This parameter specifies the custom key-value pairs of the message's payload.
	Notification          NotificationPayload `json:"notification,omitempty"`            // This parameter specifies the predefined, user-visible key-value pairs of the notification payload.
}

// NotificationPayload -> Notification payload support
type NotificationPayload struct {
	Title            string `json:"title,omitempty"`              // The notification's title.(iOS,Android,Web)
	Body             string `json:"body,omitempty"`               // The notification's body text.(iOS,Android,Web)
	AndroidChannelID string `json:"android_channel_id,omitempty"` // The notification's channel id (new in Android O).(Android)
	Icon             string `json:"icon,omitempty"`               // The notification's icon.(Android,Web)
	Sound            string `json:"sound,omitempty"`              // The sound to play when the device receives the notification.(iOS,Android)
	Badge            string `json:"badge,omitempty"`              // The value of the badge on the home screen app icon.(iOS)
	Tag              string `json:"tag,omitempty"`                // Identifier used to replace existing notifications in the notification drawer.(Android)
	Color            string `json:"color,omitempty"`              // The notification's icon color, expressed in #rrggbb format.(Android)
	ClickAction      string `json:"click_action,omitempty"`       // The action associated with a user click on the notification.(iOS,Android,Web)
	Subtitle         string `json:"subtitle,omitempty"`           // The notification's subtitle.(iOS)
	BodyLocKey       string `json:"body_loc_key,omitempty"`       // The key to the body string in the app's string resources to use to localize the body text to the user's current localization.(iOS,Android)
	BodyLocArgs      string `json:"body_loc_args,omitempty"`      // Variable string values to be used in place of the format specifiers in body_loc_key to use to localize the body text to the user's current localization.(iOS,Android)
	TitleLocKey      string `json:"title_loc_key,omitempty"`      // The key to the title string in the app's string resources to use to localize the title text to the user's current localization.(iOS,Android)
	TitleLocArgs     string `json:"title_loc_args,omitempty"`     // Variable string values to be used in place of the format specifiers in title_loc_key to use to localize the title text to the user's current localization.(iOS,Android)
}

// HTTPResponse -> Downstream HTTP message response
type HTTPResponse struct {
	StatusCode   int
	MulticastID  int64                `json:"multicast_id"`
	Success      int                  `json:"success"`
	Failure      int                  `json:"failure"`
	CanonicalIds int                  `json:"canonical_ids"`
	Results      []HTTPResponseResult `json:"results,omitempty"`
	MessageID    int64                `json:"message_id,omitempty"`
	Error        string               `json:"error,omitempty"`
}

// HTTPResponseResult -> Downstream HTTP message response result
type HTTPResponseResult struct {
	MessageID      int64  `json:"message_id"`
	RegistrationID int    `json:"registration_id"`
	Error          string `json:"error,omitempty"`
}

// Client stores the key and the Message
type Client struct {
	ServerKey string
	Message   HTTPMessage
}

// SendMessage to firebase
func (c *Client) SendMessage() (rm *HTTPResponse, err error) {
	rm = new(HTTPResponse)

	var jsonByte []byte
	if jsonByte, err = json.Marshal(c.Message); err != nil {
		return rm, err
	}

	request, _ := http.NewRequest("POST", httpEndpointURL, bytes.NewBuffer(jsonByte))
	request.Header.Set("Authorization", "key="+c.ServerKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	var response *http.Response
	response, err = client.Do(request)

	rm.StatusCode = response.StatusCode

	if err != nil {
		return rm, err
	}
	defer response.Body.Close()

	var responseBody []byte
	if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
		return rm, err
	}

	if response.StatusCode != 200 {
		return rm, nil
	}

	if err = json.Unmarshal([]byte(responseBody), &rm); err != nil {
		return rm, err
	}

	return rm, nil
}

// SetCollapseKey identifies a group of messages
func (c *Client) SetCollapseKey(collapseKey string) {
	c.Message.CollapseKey = collapseKey
}

// SetCondition sets a logical expression of conditions that determine the message target
func (c *Client) SetCondition(condition string) {
	c.Message.Condition = condition
}

// SetContentAvailable use this field to represent content-available in the APNs payload
func (c *Client) SetContentAvailable(contentAvailable bool) {
	c.Message.ContentAvailable = contentAvailable
}

// SetData sets data payload
func (c *Client) SetData(body interface{}) {
	c.Message.Data = body
}

// SetDryRun allows developers to test a request without actually sending a message
func (c *Client) SetDryRun(dryRun bool) {
	c.Message.DryRun = dryRun
}

// SetMutableContent use this field to represent mutable-content in the APNs payload
func (c *Client) SetMutableContent(mutableContent bool) {
	c.Message.MutableContent = mutableContent
}

// SetNotification sets the notification payload
func (c *Client) SetNotification(title, body, androidChannelID, icon, sound, badge, tag, color, clickAction, subtitle, bodyLocKey, bodyLocArgs, titleLocKey, titleLocArgs string) {
	c.Message.Notification.Title = title
	c.Message.Notification.Body = body
	c.Message.Notification.AndroidChannelID = androidChannelID
	c.Message.Notification.Icon = icon
	c.Message.Notification.Sound = sound
	c.Message.Notification.Badge = badge
	c.Message.Notification.Tag = tag
	c.Message.Notification.Color = color
	c.Message.Notification.ClickAction = clickAction
	c.Message.Notification.Subtitle = subtitle
	c.Message.Notification.BodyLocKey = bodyLocKey
	c.Message.Notification.BodyLocArgs = bodyLocArgs
	c.Message.Notification.TitleLocKey = titleLocKey
	c.Message.Notification.TitleLocArgs = titleLocArgs
}

// SetPriority sets the priority of the message
func (c *Client) SetPriority(priority string) {
	if priority != "high" {
		priority = "normal"
	}
	c.Message.Priority = priority
}

// SetRegistrationIds sets the recipient of a multicast message
func (c *Client) SetRegistrationIds(registrationIds []string) {
	c.Message.RegistrationIds = append(c.Message.RegistrationIds, registrationIds...)
}

// SetRestrictedPackageName sets the package name of the application where the registration tokens must match in order to receive the message
func (c *Client) SetRestrictedPackageName(restrictedPackageName string) {
	c.Message.RestrictedPackageName = restrictedPackageName
}

// SetServerKey sets the server key
func (c *Client) SetServerKey(serverKey string) {
	c.ServerKey = serverKey
}

// SetTimeToLive sets how long (in seconds) the message should be kept in FCM storage if the device is offline.
func (c *Client) SetTimeToLive(timeToLive int) {
	if timeToLive > 2419200 {
		c.Message.TimeToLive = 2419200
	} else {
		c.Message.TimeToLive = 2419200
	}
}

// SetTo sets the recipient of a message
func (c *Client) SetTo(to string) {
	c.Message.To = to
}
