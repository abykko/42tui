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

func fetchProfileCmd(m Model) ProfileData {
	log.Println("Iniciando fetchProfileCmd...")

	userLogged, err := conf.GetString("logged_with")
	if err != nil {
		log.Printf("Error obteniendo config 'logged_with': %v", err)
		return m.Profile
	}

	profile, err := getUserProfile(userLogged)
	if err != nil {
		log.Printf("Error obteniendo perfil: %v", err)
		return m.Profile
	}

	return profile
}

func getUserProfile(user string) (ProfileData, error) {
	endpoint := fmt.Sprintf("/users/profile?user=%s", user)
	log.Printf("Llamando a la API: %s", endpoint)

	resp, _, err := api.DoRequest(endpoint)
	if err != nil {
		return ProfileData{}, err
	}

	return parseProfile(resp, user)
}

func parseProfile(resp map[string]interface{}, userLogged string) (ProfileData, error) {
	key := fmt.Sprintf("https://intrapy.intra.42.fr/api/v1/users/%s", userLogged)

	raw, ok := resp[key]
	if !ok {
		return ProfileData{}, fmt.Errorf("clave no encontrada en la respuesta")
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

	return profile, nil
}

func logProfile(profile ProfileData) {
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		log.Printf("error serializando profile: %v", err)
		return
	}

	log.Println(string(data))
}

func ProfileService(m Model) (tea.Model, tea.Cmd) {

	lastUpdate := m.ProfileLastUpdate
	elapsedTime := time.Now().Unix() - lastUpdate

	if elapsedTime < 30 {
		log.Println("Tried to fetch profile. Is in cooldown!", elapsedTime, "seconds ago.")
		return m, nil
	}

	m.Profile = fetchProfileCmd(m)

	m.ProfileLastUpdate = time.Now().Unix()
		
	return m, nil
}