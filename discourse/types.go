package discourse

type GroupRequest struct {
  Group []*Group `json:"group,omitempty"`
}

type GroupCreateResponse struct {
  Group *Group `json:"basic_group,omitempty"`
}

type GroupReadResponse struct {
  Group *Group `json:"group,omitempty"`
}

type Group struct {
  Id int `json:"id,omitempty"`
  Name string `json:"name,omitempty"`
  MentionableLevel int `json:"mentionable_level,omitempty"`
  MessageableLevel int `json:"messageable_level,omitempty"`
  VisibilityLevel int `json:"visibility_level,omitempty"`
  PrimaryGroup bool `json:"primary_group,omitempty"`
  Title string `json:"title,omitempty"`
  GrantTrustLevel int `json:"grant_trust_level,omitempty"`
  FlairUrl string `json:"flair_url,omitempty"`
  FlairBackgroundColor string `json:"flair_bg_color,omitempty"`
  FlairColor string `json:"flair_color,omitempty"`
  PublicAdmission bool `json:"public_admission,omitempty"`
  PublicExit bool `json:"public_exit,omitempty"`
  AllowMembershipRequests bool `json:"allow_membership_requests,omitempty"`
  FullName string `json:"full_name,omitempty"`
  DefaultNotificationLevel int `json:"default_notification_level,omitempty"`
  MembersVisibilityLevel int `json:"members_visibility_level,omitempty"`
}

type SettingsResponse struct {
  Settings []*Setting `json:"site_settings,omitempty"`
}

type Setting struct {
  Name string `json:"setting,omitempty"`
  Description string `json:"description,omitempty"`
  Default string `json:"default,omitempty"`
  Value string `json:"value,omitempty"`
}

type Config struct {
  Host string
  Username string
  Token string
}

type ApiRequest struct {
  Method string
  Endpoint string
  JsonBody string
  Config Config
}

type DiscourseError struct {
  Message []string `json:"errors"`
  ErrorType string `json:"error_type"`
}
