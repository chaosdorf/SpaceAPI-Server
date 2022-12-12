package main

import (
	"encoding/json"
)

func (s *SpaceApi) MarshalJSON() ([]byte, error) {
	s.Api = "0.13"

	type spaceApi *SpaceApi
	return json.Marshal(spaceApi(s))
}

func (s *SpaceApi) Reset() {
	s.State.Open = nil
	s.State.Message = ""
	s.State.LastChange = 0
}

type SpaceApi struct {
	Api              string   `json:"api"`
	ApiCompatibility []string `json:"api_compatibility"`
	Space            string   `json:"space"`
	Logo             string   `json:"logo"`
	Url              string   `json:"url"`
	Location         struct {
		Address string  `json:"address"`
		Lon     float64 `json:"lon"`
		Lat     float64 `json:"lat"`
	} `json:"location"`
	Contact struct {
		Email    string `json:"email"`
		Irc      string `json:"irc"`
		Twitter  string `json:"twitter"`
		Mastodon string `json:"mastodon"`
	} `json:"contact"`
	IssueReportChannels []string `json:"issue_report_channels"`
	Feeds               map[string]struct {
		Url  string `json:"url"`
		Type string `json:"type"`
	} `json:"feeds"`
	State struct {
		Open       *bool  `json:"open,omitempty"`
		Message    string `json:"message"`
		LastChange int64  `json:"lastchange,omitempty"`
	} `json:"state"`
	Projects []string `json:"projects"`
}
