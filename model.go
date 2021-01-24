package main

import (
	"fmt"
	"time"
)

// Project is the container that we count the times for
type Project struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ClientName string `json:"client_name"`
	Status     string `json:"status"`
}

// ProjectCollection is a collection of projects (duh)
type ProjectCollection struct {
	Projects []Project `json:"projects"`
}

// Event is either start or stop of the timer
type Event struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Code      string    `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note"`
}

// EventCollection collection if events
//type EventCollection struct {
//Events []Event `json:"events"`
//}

// ProjectEventsCollection combines a project and the collection of its events
type ProjectEventsCollection struct {
	Project
	Events []Event `json:"events"`
}

// IsValid on project validates the project
func (p Project) IsValid() (bool, []error) {
	var errList []error
	switch p.Status {
	case "OPEN", "CLOSED", "ACTIVE":
		return true, nil
	default:
		errList = append(errList, fmt.Errorf("Invalid project status code %v", p.Status))
	}
	// TODO: check other errors here
	if len(errList) > 0 {
		return false, errList
	}
	return true, nil
}

// IsValid on event validates the event
func (e Event) IsValid() (bool, []error) {
	var errList []error
	switch e.Code {
	case "CHECK_IN", "CHECK_OUT":
		return true, nil
	default:
		errList = append(errList, fmt.Errorf("Invalid event action code %v", e.Code))
	}
	// TODO: check other errors here
	if len(errList) > 0 {
		return false, errList
	}
	return true, nil
}
