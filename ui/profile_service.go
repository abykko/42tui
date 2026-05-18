package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"

	"42cli/api"
	"42cli/conf"
)

func parseProfile(resp map[string]interface{}, userLogged string) (ProfileData, error) {

	// Most of the profile page data
	key := fmt.Sprintf("https://intrapy.intra.42.fr/api/v1/users/%s", userLogged)

	raw, ok := resp[key]
	if !ok {
		return ProfileData{}, fmt.Errorf("key not found")
	}

	bytes, err := json.Marshal(raw)
	if err != nil {
		return ProfileData{}, err
	}

	var profile ProfileData
	if err := json.Unmarshal(bytes, &profile); err != nil {
		return ProfileData{}, err
	}

	rawProjects, ok := resp["projects"]
	if ok {
		projBytes, err := json.Marshal(rawProjects)
		if err == nil {
			var projects []Project
			if json.Unmarshal(projBytes, &projects) == nil {
				profile.Projects = projects
			}
		}
	}

	summary, ok := resp["summary"]
	if ok {
		if s, ok := summary.(map[string]interface{}); ok {
			if lang, ok := s["language"].(string); ok {
				profile.Language = lang
			}
		}
	}

	return profile, nil
}

func fetchProfile(m Model, user string) ProfileData {

	if user == "" {
		// Use logged user as parrameter
		loggedUser, err := conf.GetString("logged_with")
		if err != nil {
			log.Printf("Error obteniendo config 'logged_with': %v", err)
			return m.Profile
		}
		user = loggedUser
	}

	// Fetch data
	endpoint := fmt.Sprintf("/users/profile?user=%s", user)
	log.Printf("[ProfileService] llamando a la API: %s", endpoint)

	resp, _, err := api.DoRequest(endpoint)
	if err != nil {
		return ProfileData{}
	}

	parsedProfile, err := parseProfile(resp, user)
	if err != nil {
		return ProfileData{}
	}
	return parsedProfile
}

func ProfileService(m Model) (tea.Model, tea.Cmd) {

	// Fetch data and update m.Profile
	m.Profile = fetchProfile(m, "")
	m.ProfileLastUpdate = time.Now().Unix()

	// Set project list content
	m.ProjectsViewport.SetContent(projects(m))
	
	return m, nil
}