package ire

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/adayoung/ada-bot/settings"
)

var APIURL string // This is read and set from config.yaml by main.init()

// Represents a single event per IRE's gamefeed
type Event struct {
	ID          int    `json:"id"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

// Gamefeed is a collection of Events, ahead of LastID
type Gamefeed struct {
	LastID int
	Events *[]Event
}

type eventsByDate []Event            // Implements sort.Interface
func (d eventsByDate) Len() int      { return len(d) }
func (d eventsByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d eventsByDate) Less(i, j int) bool {
	return d[i].Date < d[j].Date
}

// Get the latest events from API endpoint, returning them
func (g *Gamefeed) Sync() ([]Event, error) {
	url := fmt.Sprintf("%s/gamefeed.json", APIURL)
	g.LastID = settings.Settings.IRE.LastID
	if g.LastID > 0 {
		url = fmt.Sprintf("%s?id=%d", url, g.LastID)
	}

	var deathsights []Event

	if !settings.Settings.IRE.DeathsightEnabled { // Oops, we're disabled, bail out
		return deathsights, nil
	}

	if err := getJSON(url, &g.Events); err == nil {
		for _, event := range *g.Events {
			go logEvent(event)
			if event.ID > g.LastID {
				g.LastID = event.ID
			}

			if event.Type == "DEA" {
				deathsights = append(deathsights, event)
			}

                        if event.Type == "DUE" {
                                deathsights = append(deathsights, event)
                        }
		}
	} else {
		return nil, err // Error at http.Get() call
	}

	settings.Settings.IRE.LastID = g.LastID
	sort.Sort(eventsByDate(deathsights))
	return deathsights, nil
}

// Represents a player per IRE's API
type Player struct {
	Name         string `json:"name"`
	Fullname     string `json:"fullname"`
	City         string `json:"city"`
	Guild        string `json:"guild"`
	Level        string `json:"level"`
	Race         string `json:"race"`
	Kills        string `json:"kills"`
	Deaths       string `json:"deaths"`
	ExplorerRank string `json:"explorerrank"`
}


//var explorer_rank [21]string

//	explorer_rank[0] = "a Vagrant"
//      explorer_rank[1] = "a Pedestrian"
//        explorer_rank[2] = "a Landloper"
//        explorer_rank[3] = "an Itinerant Traveler"
//        explorer_rank[4] = "a Pilgrim"
//        explorer_rank[5] = "a Valleyrunner"
//        explorer_rank[6] = "a Riverwalker"
//        explorer_rank[7] = "a Roamer of the Basin"
//        explorer_rank[8] = "a Visionary of Distance"
//        explorer_rank[9] = "a Walker of Lost Lands"
//        explorer_rank[10] = "a Seeker of Avechna"
//        explorer_rank[11] = "an Ethereal Wanderer"
//        explorer_rank[12] = "a Pioneer of the Unknown"
//        explorer_rank[13] = "an Elemental Tracker"
//        explorer_rank[14] = "a Worldwalker"
//        explorer_rank[15] = "a Planar Drifter"
//        explorer_rank[16] = "a Cosmic Wayfarer"
//        explorer_rank[17] = "a Voyager to the Beyond"
//        explorer_rank[18] = "an Astral Traveller"
//        explorer_rank[19] = "a Messenger of the Fates"
//        explorer_rank[20] = "a Planeswalker"

func (s *Player) String() string {
	player := fmt.Sprintf(`
           Name: %s
       Fullname: %s
           Race: %s
            Org: %s
	      Guild: %s
          Kills: %s
         Deaths: %s
   ExplorerRank: %s`,
		s.Name, s.Fullname, s.Race, strings.Title(s.City), 
		strings.Title(s.Guild), s.Kills, s.Deaths, s.ExplorerRank)
	return player
}

// Lookup and retrieve a player from IRE's API
func GetPlayer(player string) (*Player, error) {
	if !(len(player) > 0) {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid player name supplied: %s", player))
	}
	if match, err := regexp.MatchString("(?i)[^a-z]+", player); err == nil {
		if !match {
			url := fmt.Sprintf("%s/characters/%s.json", APIURL, player)
			_player := &Player{}
			if err := getJSON(url, &_player); err == nil {
				return _player, nil
			} else {
				return nil, err // Error at getJson() call
			}
		} else {
			return nil, fmt.Errorf(fmt.Sprintf("Invalid player name supplied: %s", player))
		}
	} else {
		return nil, err // Error at regexp.MatchString() call
	}
}

