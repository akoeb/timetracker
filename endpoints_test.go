package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	db = initDB("database.db")
)

/*func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}*/

//	apis.POST("/projects", createProject(db))
func TestCreateProject(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// test project
	tp := createTestProject(0)
	projectJSON, err := toJSON(*tp)
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("POST", "/projects", projectJSON, createProject(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusCreated, rec.Code)

		// get Project from JSON result
		rp, err := fromJSON(rec.Body.String(), Project{})
		assert.NoError(t, err)
		ret, ok := rp.(Project)
		assert.True(t, ok)

		// validate ID has been autogenerated:
		assert.Equal(t, ret.ID, 1, "ID autogenerated")

		// validate the rest of the object:
		tp.ID = ret.ID
		assert.Equal(t, *tp, ret, "Project not as expected: ", rp)
	}
}

//	apis.GET("/projects", showProjectList(db))
func TestShowProjectList(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp1, err := db.upsertProject(createTestProject(1))
	assert.NoError(t, err)

	// create a second test project
	tp2, err := db.upsertProject(createTestProject(2))
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("GET", "/projects", "", showProjectList(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Project from JSON result
		rp, err := fromJSON(rec.Body.String(), ProjectCollection{})
		assert.NoError(t, err)
		ret, ok := rp.(ProjectCollection)
		assert.True(t, ok)

		// validate number of projects:
		assert.Equal(t, len(ret.Projects), 2, "ProjectCollection length not es expected", rp)

		// validate the first project
		assert.Equal(t, tp1, &ret.Projects[0], "Project 1 not as expected: ", rp)

		// validate the second project
		assert.Equal(t, tp2, &ret.Projects[1], "Project 2 not as expected: ", rp)
	}
}

//	apis.GET("/projects/:projectid", showProjectDetails(db))
func TestShowSingleProject(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(0))
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("GET", fmt.Sprintf("/projects/%v", tp.ID), "", showProjectDetails(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Project from JSON result
		rp, err := fromJSON(rec.Body.String(), Project{})
		assert.NoError(t, err)
		ret, ok := rp.(Project)
		assert.True(t, ok)

		// validate ID has been autogenerated:
		assert.Equal(t, ret.ID, 1, "ID autogenerated")

		// validate the rest of the object:
		assert.Equal(t, *tp, ret, "Project not as expected: ", ret)
	}
}

//	apis.PUT("/projects/:projectid", updateProject(db))
func TestUpdateProject(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(0))
	assert.NoError(t, err)

	// change project
	tp.ClientName = "different client name"
	tp.Name = "different name"

	// serialize changed project
	projectJSON, err := toJSON(*tp)
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("PUT", fmt.Sprintf("/projects/%v", tp.ID), projectJSON, updateProject(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Project from JSON result
		rp, err := fromJSON(rec.Body.String(), Project{})
		assert.NoError(t, err)
		ret, ok := rp.(Project)
		assert.True(t, ok)

		// validate the rest of the object:
		assert.Equal(t, *tp, ret, "Project not as expected: ", ret)
	}

	// the update is ok, now verify that a subsequent Get also returns expected result:
	// execute
	rec, err = executeTest("GET", fmt.Sprintf("/projects/%v", tp.ID), "", showProjectDetails(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Project from JSON result
		rp, err := fromJSON(rec.Body.String(), Project{})
		assert.NoError(t, err)
		ret, ok := rp.(Project)
		assert.True(t, ok)

		// validate the rest of the object:
		assert.Equal(t, *tp, ret, "Project not as expected: ", ret)
	}

}

// Delete only works on projects with  no associated events
//	apis.DELETE("/projects/:projectid", deleteProject(db))
func TestDeleteProject(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp1, err := db.upsertProject(createTestProject(1))
	assert.NoError(t, err)

	// create a second test project
	tp2, err := db.upsertProject(createTestProject(2))
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("DELETE", fmt.Sprintf("/projects/%v", tp1.ID), "", deleteProject(db))

	// validate result
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// now validate that a subsequent project list only returns the rest
	rec, err = executeTest("GET", "/projects", "", showProjectList(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Project from JSON result
		rp, err := fromJSON(rec.Body.String(), ProjectCollection{})
		assert.NoError(t, err)
		ret, ok := rp.(ProjectCollection)
		assert.True(t, ok)

		// validate number of projects:
		assert.Equal(t, len(ret.Projects), 1, "ProjectCollection length not es expected", rp)

		// validate the first project
		assert.Equal(t, tp2, &ret.Projects[0], "Project 1 not as expected: ", rp)
	}
}

//	apis.GET("/projects/:projectid/events", showProjectEventHistory(db))
func TestShowProjectEventHistory(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(1))
	assert.NoError(t, err)

	ev1, err := db.upsertEvent(createTestEvent(tp.ID, 1))
	assert.NoError(t, err)

	ev2, err := db.upsertEvent(createTestEvent(tp.ID, 2))
	assert.NoError(t, err)

	ev3, err := db.upsertEvent(createTestEvent(tp.ID, 3))
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("GET", fmt.Sprintf("/projects/%v/events", tp.ID), "", showProjectEventHistory(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		var rp EventCollection
		res, err := fromJSON(rec.Body.String(), EventCollection{})
		assert.NoError(t, err)
		re, ok := res.(EventCollection)
		assert.True(t, ok)

		// validate number of projects:
		assert.Equal(t, len(re.Events), 3, "EventCollection length not es expected", rp)

		// validate the events
		assert.Equal(t, ev1, &re.Events[0], "Event 1 not as expected: ", rp)
		assert.Equal(t, ev2, &re.Events[1], "Event 2 not as expected: ", rp)
		assert.Equal(t, ev3, &re.Events[2], "Event 3 not as expected: ", rp)

	}
}

//	apis.POST("/projects/:projectid/events", createProjectEvent(db))
func TestCreateProjectEvent(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(1))
	assert.NoError(t, err)

	te := createTestEvent(tp.ID, 1)
	eventJSON, err := toJSON(*te)
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("POST", fmt.Sprintf("/projects/%v/events", tp.ID), eventJSON, createProjectEvent(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusCreated, rec.Code)

		// get Event from JSON result
		res, err := fromJSON(rec.Body.String(), Event{})
		assert.NoError(t, err)
		re, ok := res.(Event)
		assert.True(t, ok)

		// validate ID has been autogenerated:
		assert.Equal(t, re.ID, 1, "ID autogenerated")

		// validate the rest of the object:
		te.ID = re.ID
		assert.Equal(t, *te, re, "Event not as expected: ", re)
	}
}

//	apis.GET("/projects/:projectid/events/:eventid/", showProjectEvent(db))
func TestShowSingleProjectEvent(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(0))
	assert.NoError(t, err)

	te, err := db.upsertEvent(createTestEvent(tp.ID, 1))
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("GET", fmt.Sprintf("/projects/%v/events/%v", tp.ID, te.ID), "", showProjectEvent(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Event from JSON result
		res, err := fromJSON(rec.Body.String(), Event{})
		assert.NoError(t, err)
		re, ok := res.(Event)
		assert.True(t, ok)

		assert.Equal(t, *te, re, "Event not as expected: ", re)
	}
}

//	apis.PUT("/projects/:projectid/events/:eventid/", updateProjectEvent(db))
func TestUpdateProjectEvent(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(0))
	assert.NoError(t, err)

	te, err := db.upsertEvent(createTestEvent(tp.ID, 1))
	assert.NoError(t, err)

	// change event
	te.Code = "CHECK_OUT"
	te.Timestamp, err = time.Parse(time.RFC3339, "2012-10-10T13:47:23.522Z")
	assert.NoError(t, err)
	v, errlist := te.IsValid()
	assert.Nil(t, errlist)
	assert.True(t, v)

	// serialize changed event
	eventJSON, err := toJSON(*te)
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("PUT", fmt.Sprintf("/projects/%v/events/%v", tp.ID, te.ID), eventJSON, updateProjectEvent(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Event from JSON result
		rp, err := fromJSON(rec.Body.String(), Event{})
		assert.NoError(t, err)
		ret, ok := rp.(Event)
		assert.True(t, ok)

		// validate the rest of the object:
		assert.Equal(t, *te, ret, "Event not as expected: ", ret)
	}

	// the update is ok, now verify that a subsequent Get also returns expected result:
	// execute
	rec, err = executeTest("GET", fmt.Sprintf("/projects/%v/events/%v", tp.ID, te.ID), "", showProjectEvent(db))

	// validate
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		// get Project from JSON result
		re, err := fromJSON(rec.Body.String(), Event{})
		assert.NoError(t, err)
		ret, ok := re.(Event)
		assert.True(t, ok)

		// validate the rest of the object:
		assert.Equal(t, *te, ret, "Event not as expected: ", ret)
	}

}

//	apis.DELETE("/projects/:projectid/events/:eventid/", deleteProjectEvent(db))
func TestDeleteEvent(t *testing.T) {

	// reset db
	assert.NoError(t, db.resetDB())

	// create a test project
	tp, err := db.upsertProject(createTestProject(0))
	assert.NoError(t, err)

	te1, err := db.upsertEvent(createTestEvent(tp.ID, 1))
	assert.NoError(t, err)

	te2, err := db.upsertEvent(createTestEvent(tp.ID, 1))
	assert.NoError(t, err)

	// execute
	rec, err := executeTest("DELETE", fmt.Sprintf("/projects/%v/events/%v", tp.ID, te1.ID), "", deleteProjectEvent(db))

	// validate result
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// now validate that a subsequent event list only returns the rest
	rec, err = executeTest("GET", fmt.Sprintf("/projects/%v/events", tp.ID), "", showProjectEventHistory(db))
	if assert.NoError(t, err) {
		// HTTP Status Code
		assert.Equal(t, http.StatusOK, rec.Code)

		var rp EventCollection
		res, err := fromJSON(rec.Body.String(), EventCollection{})
		assert.NoError(t, err)
		re, ok := res.(EventCollection)
		assert.True(t, ok)

		// validate number of projects:
		assert.Equal(t, len(re.Events), 1, "EventCollection length not es expected", rp)

		// validate the events
		assert.Equal(t, te2, &re.Events[0], "Event 1 not as expected: ", rp)
	}
}

//	apis.GET("/reports/projects/:projectid", reportOnProject(db))
//	apis.GET("/reports/date/:datestr", reportOnTime(db))
