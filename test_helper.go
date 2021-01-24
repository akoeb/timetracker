package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func createTestProject(cnt int) *Project {
	retVal := &Project{
		ID:         0,
		ClientName: fmt.Sprintf("CLIENT NAME %d", cnt),
		Name:       fmt.Sprintf("PROJECT NAME %d", cnt),
		Status:     fmt.Sprint("OPEN"),
	}
	return retVal
}
func createTestEvent(projectID int, cnt int) *Event {
	startTime, _ := time.Parse(time.RFC3339, "2010-10-10T13:47:23.522Z")
	eventTimeOffset, _ := time.ParseDuration(fmt.Sprintf("%vh", 13*cnt))
	retVal := &Event{
		ID:        0,
		ProjectID: projectID,
		Code:      fmt.Sprintf("CHECK_IN"),
		Timestamp: startTime.Add(eventTimeOffset),
		Note:      fmt.Sprintf("EVENT NOTE %d", cnt),
	}
	return retVal
}

func toJSON(obj interface{}) (string, error) {
	retVal := ""
	var err error
	var b []byte
	switch v := obj.(type) {
	case nil:
		err = fmt.Errorf("can not convert nil object to JSON")
	case Project, ProjectCollection, Event, ProjectEventsCollection:
		b, err = json.Marshal(v)
		retVal = string(b)
	default:
		err = fmt.Errorf("object has unknown type, cannot convert to JSON")
	}
	return retVal, err
}

// converts a json string to a model. use with type assertion
func fromJSON(jsonStr string, model interface{}) (interface{}, error) {
	switch model.(type) {
	case nil:
		return nil, fmt.Errorf("can not convert nil object to JSON")
	case Project:
		var v Project
		if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
			return nil, err
		}
		return v, nil
	case ProjectCollection:
		var v ProjectCollection
		if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
			return nil, err
		}
		return v, nil
	case Event:
		var v Event
		if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
			return nil, err
		}
		return v, nil
	case ProjectEventsCollection:
		var v ProjectEventsCollection
		if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
			return nil, err
		}
		return v, nil

	default:
		return nil, fmt.Errorf("object has unknown type, cannot convert to JSON")
	}
}

// func serializeTestProject(p *Project) (string, error) {
// 	b, err := json.Marshal(p)
// 	return string(b), err
// }
// func deserializeResultToTestProject(s string) (*Project, error) {
// 	var p Project
// 	b := []byte(s)
// 	err := json.Unmarshal(b, &p)
// 	return &p, err
// }
// func deserializeResultToTestProjectCollection(s string) (*ProjectCollection, error) {
// 	var p ProjectCollection
// 	b := []byte(s)
// 	err := json.Unmarshal(b, &p)
// 	return &p, err
// }

func executeTest(method string, path string, body string, hFunc echo.HandlerFunc) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	var req *http.Request
	req = httptest.NewRequest(method, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// add path to context
	c.SetPath(path)

	// now check path params
	rp, _ := regexp.Compile("/projects/(\\d+)")
	re, _ := regexp.Compile("/events/(\\d+)")
	if rp.MatchString(path) && re.MatchString(path) {
		c.SetParamNames("projectid", "eventid")
		c.SetParamValues(rp.FindStringSubmatch(path)[1], re.FindStringSubmatch(path)[1])
	} else if rp.MatchString(path) {
		c.SetParamNames("projectid")
		c.SetParamValues(rp.FindStringSubmatch(path)[1])
	}

	err := hFunc(c)
	return rec, err
}
