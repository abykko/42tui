package service

import (
	"encoding/json"
	"fmt"
	"log"

	"42cli/api"
	"42cli/conf"
)

type Milestone struct {
	Deadline    string  `json:"deadline"`
	Level       int     `json:"level"`
	MilestoneID int     `json:"milestone_id"`
	UserID      int     `json:"user_id"`
	ValidatedAt string `json:"validated_at"`
}

type Pace struct {
	ActivatedAt            string      `json:"activated_at"`
	CursusBeginDate        string      `json:"cursus_begin_date"`
	Deadline               string      `json:"deadline"`
	EndOfCursusMaxDeadline string      `json:"end_of_cursus_max_deadline"`
	EtaEndOfCursus         string      `json:"eta_end_of_cursus"`
	IsActivated            bool        `json:"is_activated"`
	IsCampusActivated      bool        `json:"is_campus_activated"`
	Milestone              int         `json:"milestone"`
	Milestones             []Milestone `json:"milestones"`
	Pace                   int         `json:"pace"`
	PaceSpeedUp            *int        `json:"pace_speed_up"`
	ProbationaryPeriod     bool        `json:"probationary_period"`
}

type Team struct {
	FinalMark     int    `json:"final_mark"`
	IsValidated   bool   `json:"is_validated"`
	LastEventDate string `json:"last_event_date"`
	Occurrence    int    `json:"occurrence"`
}

type Project struct {
	Children       []any  `json:"children"`
	FinalMark      int    `json:"final_mark"`
	IsValidated    bool   `json:"is_validated"`
	LastEventDate  string `json:"last_event_date"`
	Occurrence     int    `json:"occurrence"`
	ProjectName    string `json:"project_name"`
	ProjectSlug    string `json:"project_slug"`
	ProjectsUserID int    `json:"projects_user_id"`
	Teams          []Team `json:"teams"`
}

type ProfileData struct {
	ID               int     `json:"id"`
	UserDataID       int     `json:"user_data_id"`
	Login            string  `json:"login"`
	DisplayedLogin   string  `json:"displayed_login"`
	Email            string  `json:"email"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Phone            string  `json:"phone"`
	Image            string  `json:"image"`
	ProfilePicture   string  `json:"profile_picture"`
	IsActive         bool    `json:"is_active"`
	Wallet           int     `json:"wallet"`
	EvaluationPoints int     `json:"evaluation_points"`
	AlumnizedAt      *string `json:"alumnized_at"`
	Close            *string `json:"close"`
	DataErasureDate  string  `json:"data_erasure_date"`
	Location         string `json:"location"`
	Groups           []any   `json:"groups"`
	Titles           []any   `json:"titles"`
	Projects         []Project `json:"projects"`
	Pace             Pace    `json:"pace"`
	Language         string  `json:"language"`
}

func parseProfile(resp map[string]interface{}, userLogged string) (ProfileData, error) {

	// Most of the profile page data
	key := "profile"

	raw, ok := resp[key]
	if !ok {
		log.Printf("[ProfileService] Error: User key '%s' not found in response map", key)
		return ProfileData{}, fmt.Errorf("key not found")
	}

	bytes, err := json.Marshal(raw)
	if err != nil {
		log.Printf("[ProfileService] Error marshaling raw user data: %v", err)
		return ProfileData{}, err
	}

	var profile ProfileData
	if err := json.Unmarshal(bytes, &profile); err != nil {
		log.Printf("[ProfileService] Error unmarshaling profile structural data: %v", err)
		return ProfileData{}, err
	}
	log.Println("[ProfileService] Base profile data successfully parsed")

	rawPace, ok := resp["pace"]
	if ok {
		log.Println("[ProfileService] Pace data found, attempting to parse...")
		paceBytes, err := json.Marshal(rawPace)
		if err == nil {
			var pace Pace
			if json.Unmarshal(paceBytes, &pace) == nil {
				profile.Pace = pace
				log.Println("[ProfileService] Successfully parsed pace data")
			} else {
				log.Println("[ProfileService] Warning: Failed to unmarshal pace JSON")
			}
		} else {
			log.Printf("[ProfileService] Warning: Failed to marshal raw pace data: %v", err)
		}
	} else {
		log.Println("[ProfileService] No pace field found in response")
	}

	rawProjects, ok := resp["projects"]
	if ok {
		log.Println("[ProfileService] Project data found, attempting to parse...")
		projBytes, err := json.Marshal(rawProjects)
		if err == nil {
			var projects []Project
			if json.Unmarshal(projBytes, &projects) == nil {
				profile.Projects = projects
				log.Printf("[ProfileService] Successfully parsed %d projects", len(projects))
			} else {
				log.Println("[ProfileService] Warning: Failed to unmarshal projects JSON")
			}
		} else {
			log.Printf("[ProfileService] Warning: Failed to marshal raw projects data: %v", err)
		}
	} else {
		log.Println("[ProfileService] No projects field found in response")
	}

	summary, ok := resp["summary"]
	if ok {
		log.Println("[ProfileService] Summary data found, attempting to parse language...")
		if s, ok := summary.(map[string]interface{}); ok {
			if lang, ok := s["language"].(string); ok {
				profile.Language = lang
				log.Printf("[ProfileService] Language identified: %s", lang)
			} else {
				log.Println("[ProfileService] Warning: 'language' field missing or not a string in summary")
			}
		} else {
			log.Println("[ProfileService] Warning: 'summary' field is not a valid map structure")
		}
	}

	log.Printf("[ProfileService] Parsing completed for user: %s", userLogged)
	return profile, nil
}

func fetchProfile(user string) (ProfileData, error) {
	log.Printf("[ProfileService] Request received to fetch profile for user: '%s'", user)

	if user == "" {
		log.Println("[ProfileService] User not provided, retrieving default fallback logged_with user...")
		loggedUser, err := conf.GetString("logged_with")
		if err != nil {
			log.Printf("[ProfileService] Error retrieving logged_with from configuration: %v", err)
			return ProfileData{}, err
		}
		user = loggedUser
		log.Printf("[ProfileService] Fallback user resolved to: '%s'", user)
	}

	endpoint := fmt.Sprintf("/users/profile?user=%s", user)
	log.Printf("[ProfileService] Making API request to endpoint: %s", endpoint)

	resp, _, err := api.DoRequest(endpoint)
	if err != nil {
		log.Printf("[ProfileService] API request failed for endpoint %s: %v", endpoint, err)
		return ProfileData{}, err
	}
	log.Println("[ProfileService] API request succeeded, proceeding to parse payload")

	return parseProfile(resp, user)
}

func isFreeze() (bool, error) {
	endpoint := "/freeze"
	log.Printf("[ProfileService] Making API request to endpoint: %s", endpoint)

	resp, _, err := api.DoRequest(endpoint)
	if err != nil {
		log.Printf("[ProfileService] API request failed for endpoint %s: %v", endpoint, err)
		return false, err
	}

	freeze, ok := resp["freeze"].(bool)
	if !ok {
		return false, fmt.Errorf("invalid type for freeze")
	}

	return freeze, nil
}

func Profile(user string) (ProfileData, error) {
	log.Printf("[ProfileService] External call to Profile() invoked for user: '%s'", user)

	freeze, err := isFreeze()
	if err != nil {
		return ProfileData{}, err
	}

	if freeze {
		return ProfileData{}, fmt.Errorf("freeze")
	}

	return fetchProfile(user)
}