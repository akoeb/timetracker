package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// apis.GET("/projects", showProjectList(db))
// TODO: status filters in querystring
func showProjectList(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		projects, err := db.getProjectList()
		if err != nil {
			ctx.Logger().Infof("showAllProjects: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read projects")
		}
		if len(projects.Projects) < 1 {
			return echo.NewHTTPError(http.StatusNotFound, "Project not Found")
		}
		return ctx.JSON(http.StatusOK, projects)
	}
}

// apis.POST("/projects", createProject(db))
func createProject(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		project := &Project{}

		// bind body into struct
		err := ctx.Bind(project)
		if err != nil {
			ctx.Logger().Infof("createProject: bind error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// field validation
		if ok, errors := project.IsValid(); !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Wrong Input: %v", errors))
		}

		// creates can not have an id
		if project.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create project that has an id")
		}

		// write to db
		err = db.upsertProject(project)
		if err != nil {
			ctx.Logger().Infof("createProject: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write project")
		}

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, project)
	}
}

// apis.GET("/projects/:projectid", showProjectDetails(db))
func showProjectDetails(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		project, err := db.getProjectByID(id)
		if err != nil {
			ctx.Logger().Infof("showOneProject: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read project")
		}
		if project.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Project not Found")
		}
		return ctx.JSON(http.StatusOK, project)
	}
}

// apis.PUT("/projects/:projectid", updateProject(db))
func updateProject(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// read parameter
		id, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// bind body into struct
		project := &Project{}
		err = ctx.Bind(&project)
		if err != nil {
			ctx.Logger().Infof("updateProject: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		// some validation
		if project.ID != id {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("project with id %d can not be updated in path with id %d", project.ID, id))
		}

		// do database operation
		err = db.upsertProject(project)
		if err != nil {
			ctx.Logger().Infof("updateProject: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change project")

		}

		// inform the client
		return ctx.JSON(http.StatusOK, project)
	}
}

// // Delete only works on projects with  no associated events
// apis.DELETE("/projects/:projectid", deleteProject(db))
func deleteProject(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		projectid, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// do the database magic
		_, err = db.deleteProjectByID(projectid)
		if err != nil {
			ctx.Logger().Infof("deleteProject: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete project")
		}

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}

// apis.GET("/projects/:projectid/events",showProjectEventHistory(db))
func showProjectEventHistory(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		projectid, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		events, err := db.getProjectByID(projectid)
		if err != nil {
			ctx.Logger().Infof("showProjectEventHistory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get event list")
		}
		return ctx.JSON(http.StatusOK, events)
	}
}

// apis.POST("/projects/:projectid/events",createProjectEvent(db))
func createProjectEvent(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		event := &Event{}

		// get parameter
		projectid, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// bind body into struct
		err = ctx.Bind(event)
		if err != nil {
			ctx.Logger().Infof("createEvent: bind error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// field validation
		if ok, errors := event.IsValid(); !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Wrong Input: %v", errors))
		}

		// creates can not have an id
		if event.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create event that has an id")
		}

		// overwrite projectID from JSON object, the correct one is in the path
		event.ProjectID = projectid

		// write to db
		err = db.upsertEvent(event)
		if err != nil {
			ctx.Logger().Infof("createEvent: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write event")
		}

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, event)
	}
}

// apis.GET("/projects/:projectid/events/:eventid/",showProjectEvent(db))
func showProjectEvent(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		projectid, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter %v", projectid)
		}
		eventid, err := strconv.Atoi(ctx.Param("eventid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter %v", eventid)
		}
		event, err := db.getProjectEventByID(projectid, eventid)
		if event.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Event not Found")
		}
		return ctx.JSON(http.StatusOK, event)
	}
}

// apis.PUT("/projects/:projectid/events/:eventid/",updateProjectEvent(db))
func updateProjectEvent(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		projectid, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter %v", projectid)
		}
		eventid, err := strconv.Atoi(ctx.Param("eventid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter %v", eventid)
		}
		// bind body into struct
		event := &Event{}
		err = ctx.Bind(&event)
		if err != nil {
			ctx.Logger().Infof("updateProjectEvent: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// make sure, malevolent users did not overwrite IDs in JSON:
		if event.ID != eventid {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("event with id %d can not be updated in path with id %d", event.ID, eventid))
		}
		if event.ProjectID != projectid {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("event with project_id %d can not be updated in path with id %d", event.ProjectID, projectid))
		}
		// make sure the event exists in the project
		dbevent, err := db.getProjectEventByID(projectid, eventid)
		if dbevent.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Event not Found")
		}

		err = db.upsertEvent(event)
		if err != nil {
			ctx.Logger().Infof("updateEvent: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change event")
		}

		// inform the client
		return ctx.JSON(http.StatusOK, event)
	}
}

// apis.DELETE("/projects/:projectid/events/:eventid/", deleteProjectEvent(db))
func deleteProjectEvent(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		projectid, err := strconv.Atoi(ctx.Param("projectid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		eventid, err := strconv.Atoi(ctx.Param("eventid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter %v", eventid)
		}

		// do the database magic
		deleted, err := db.deleteEventByProjectAndID(projectid, eventid)
		if err != nil {
			ctx.Logger().Infof("deleteProject: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete project")
		}
		if deleted < 1 {
			return echo.NewHTTPError(http.StatusNotFound, "Event not Found")
		}
		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}
