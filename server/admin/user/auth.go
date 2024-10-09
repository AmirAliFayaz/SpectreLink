package user

import "SpectreLink/log"

func (s *TelnetSession) DoAuthenticate() bool {
	if err := s.Titlef("SpectreLink Login"); err != nil {
		return false
	}
	
	prompt, err := s.Promptf("Username » ")
	if err != nil {
		log.Exception(err, "Failed to read username")
		return false
	}
	
	password, err := s.Password("Password » ")
	if err != nil {
		log.Exception(err, "Failed to read password")
		return false
	}
	
	log.Infof("username %s, password %s", prompt, password)
	return true
}
