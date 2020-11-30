package twitter

// Credentials stores the secrets/tokens/keys necessary for authentication with the Twitter API.
type Credentials struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}
