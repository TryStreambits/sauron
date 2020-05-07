package sauron

// #region General

// TwitchGqlResponseUser is various user data from GQL
type TwitchGqlUser struct {
	ID                    string                     `json:"id"`
	BroadcastSettings     TwitchGqlBroadcastSettings `json:"broadcastSettings"`
	DisplayName           string                     `json:"displayName"`
	Login                 string                     `json:"login"`
	ProfileImageURL       string                     `json:"profileImageURL"`
	MediumProfileImageURL string                     `json:"medProfileImageUrl"`
	Roles                 TwitchGqlRoles             `json:"roles"`
}

// TwitchGqlBroadcastSettings is various broadcast settings
type TwitchGqlBroadcastSettings struct {
	ID       string        `json:"id"`
	Language string        `json:"language"`
	Game     TwitchGqlGame `json:"game,omitempty"`
	Title    string        `json:"title"`
}

// TwitchGqlGame is various game related settings
type TwitchGqlGame struct {
	ID          string `json:"id"`
	BoxArtURL   string `json:"boxArtURL"`
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
	TypeName    string `json:"__typename"`
}

// #endregion

// #region Channel

// TwitchGqlChannelResponse is some of the possible GQL response we get from the Twitch endpoint if it is a channel
type TwitchGqlChannelResponse struct {
	Data TwitchGqlData `json:"data"`
}

// TwitchGqlData is some of the possible GQL data
type TwitchGqlData struct {
	CurrentUser string          `json:"currentUser,omitempty"`
	Stream      TwitchGqlStream `json:"stream,omitempty"`
	User        TwitchGqlUser   `json:"user"`
}

type TwitchGqlRoles struct {
	IsAffiliate bool   `json:"isAffiliate"`
	IsPartner   bool   `json:"isPartner"`
	IsStaff     bool   `json:"isStaff,omitempty"`
	TypeName    string `json:"__typename"`
}

type TwitchGqlStream struct {
	Type string `json:"type,omitempty"`
}

// #endregion

// #region Clip

type TwitchGqlClipResponse struct {
	Data TwitchGqlClipRoot `json:"data"`
}

type TwitchGqlClipRoot struct {
	Clip TwitchGqlClipData `json:"clip"`
}

type TwitchGqlClipData struct {
	Broadcaster TwitchGqlUser `json:"broadcaster"`
	Game        TwitchGqlGame `json:"game,omitempty"`
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
}

// #endregion
