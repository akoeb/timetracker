package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Database connection and statement methods
type Database struct {
	conn     *sql.DB
	filepath string
	models   map[string]DbModel
}

// DbModel holds the most commonly used queries
type DbModel struct {
	fieldList   []string
	table       string
	selectQuery string
	insertQuery string
	updateQuery string
	deleteQuery string
}

func (m *DbModel) fields() string {
	return strings.Join(m.fieldList, ", ")
}
func (m *DbModel) questionMarks() string {
	var arr2 []string
	for range m.fieldList {
		arr2 = append(arr2, "?")
	}
	return strings.Join(arr2, ", ")
}
func (m *DbModel) fieldsAndQuestionMarks() string {
	return strings.Join(m.fieldList, " = ?, ") + " = ?"
}

// NewDbModel is the factory method to create a DbModel
func NewDbModel(table string, fieldList []string) DbModel {
	m := &DbModel{
		table:     table,
		fieldList: fieldList,
	}
	m.selectQuery = fmt.Sprintf("SELECT id, %s FROM %s", m.fields(), table)
	m.insertQuery = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, m.fields(), m.questionMarks())
	m.updateQuery = fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", table, m.fieldsAndQuestionMarks())
	m.deleteQuery = fmt.Sprintf("DELETE FROM %s", table)
	return *m
}

func (db *Database) connect() {
	conn, err := sql.Open("sqlite3", "file:"+db.filepath+"?foreign_keys=on")
	if err != nil {
		panic(err)
	}
	db.conn = conn
	// Here we check for any db errors then exit
}

func initDB(filepath string) *Database {
	db := &Database{
		filepath: filepath,
		models:   make(map[string]DbModel),
	}
	db.connect()

	// initialize the DB queries for the models Event and Project:
	db.models["project"] = NewDbModel("projects", []string{"name", "client_name", "status"})
	db.models["event"] = NewDbModel("events", []string{"project_id", "code", "timestamp", "note"})

	return db
}

// MIGRATION is now done outside app with migrate

// now the db functions:
func (db *Database) getProjectList() (*ProjectCollection, error) {
	result := ProjectCollection{}
	rows, err := db.conn.Query(db.models["project"].selectQuery)

	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return &result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		item := Project{}
		err = rows.Scan(&item.ID, &item.Name, &item.ClientName, &item.Status)
		// Exit if we get an error
		if err != nil {
			return &result, err
		}
		result.Projects = append(result.Projects, item)
	}
	return &result, nil
}
func (db *Database) upsertProject(project *Project) (*Project, error) {

	doInsert := true
	if project.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = db.models["project"].insertQuery
	} else {
		query = db.models["project"].updateQuery
	}

	// Create a prepared SQL statement
	stmt, err := db.conn.Prepare(query)
	// Exit if we get an error
	if err != nil {
		return nil, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	var result sql.Result
	if doInsert {
		result, err = stmt.Exec(project.Name, project.ClientName, project.Status)
	} else {
		result, err = stmt.Exec(project.Name, project.ClientName, project.Status, project.ID)
	}
	// Exit if we get an error
	if err != nil {
		return nil, err
	}

	// in insert, read the autoincremented id back into struct
	if doInsert {
		id64, err := result.LastInsertId()
		project.ID = int(id64)
		if err != nil {
			return nil, err
		}
	}
	return project, nil
}

func (db *Database) getProjectByID(id int) (Project, error) {
	result := Project{}
	sql := db.models["project"].selectQuery + " WHERE id = ?"
	rows, err := db.conn.Query(sql, id)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result.ID, &result.Name, &result.ClientName, &result.Status)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// TODO: error if project not empty
func (db *Database) deleteProjectByID(id int) (int, error) {

	// first all items in the category:
	sql := db.models["project"].deleteQuery + " WHERE id = ?"

	// Create a prepared SQL statement
	stmt, err := db.conn.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	result, err := stmt.Exec(id)

	// Exit if we get an error
	if err != nil {
		return 0, err
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(numDeleted), nil
}

func (db *Database) getProjectEventList(projectid int) ([]Event, error) {
	result := []Event{}
	sql := db.models["event"].selectQuery + " WHERE project_id = ?"

	rows, err := db.conn.Query(sql, projectid)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		item := Event{}
		err = rows.Scan(&item.ID, &item.ProjectID, &item.Code, &item.Timestamp, &item.Note)

		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}
func (db *Database) getProjectEventByID(projectID int, eventID int) (*Event, error) {
	result := Event{}
	sql := db.models["event"].selectQuery + " WHERE id = ? AND project_id = ?"
	rows, err := db.conn.Query(sql, eventID, projectID)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return nil, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result.ID, &result.ProjectID, &result.Code, &result.Timestamp, &result.Note)
		// Exit if we get an error
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}
func (db *Database) upsertEvent(event *Event) (*Event, error) {

	doInsert := true
	if event.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = db.models["event"].insertQuery
	} else {
		query = db.models["event"].updateQuery
	}

	// Create a prepared SQL statement
	stmt, err := db.conn.Prepare(query)
	// Exit if we get an error
	if err != nil {
		return nil, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	var result sql.Result
	if doInsert {
		result, err = stmt.Exec(event.ProjectID, event.Code, event.Timestamp, event.Note)
	} else {
		result, err = stmt.Exec(event.ProjectID, event.Code, event.Timestamp, event.Note, event.ID)
	}
	// Exit if we get an error
	if err != nil {
		return nil, err
	}

	// in insert, read the autoincremented id back into struct
	if doInsert {
		id64, err := result.LastInsertId()
		event.ID = int(id64)
		if err != nil {
			return nil, err
		}
	}
	return event, nil
}
func (db *Database) deleteEventByProjectAndID(projectID int, eventID int) (int, error) {

	// first all items in the category:
	sql := db.models["event"].deleteQuery + " WHERE project_id = ? AND id = ?"

	// Create a prepared SQL statement
	stmt, err := db.conn.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	result, err := stmt.Exec(projectID, eventID)

	// Exit if we get an error
	if err != nil {
		return 0, err
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(numDeleted), nil
}

// get last event in project
func (db *Database) getLastEventByProject(projectID int) (*Event, error) {
	result := &Event{}
	sql := db.models["event"].selectQuery + " WHERE project_id = ? ORDER BY timestamp desc"

	rows, err := db.conn.Query(sql, projectID)

	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	// only read the last item (first when sorting desc)
	if rows.Next() {
		err = rows.Scan(&result.ID, &result.ProjectID, &result.Code, &result.Timestamp, &result.Note)

		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// get next possible event type
func (db *Database) nextEventType(projectID int) (string, error) {
	event, err := db.getLastEventByProject(projectID)
	if err != nil {
		return "", err
	}

	//no events yet:
	if event.ID == 0 {
		return "CHECK_IN", nil
	}

	if event.Code == "CHECK_IN" {
		return "CHECK_OUT", nil
	} else if event.Code == "CHECK_OUT" {
		return "CHECK_IN", nil
	}

	// should not happen:
	return "", fmt.Errorf("unexpected Event Code found in Database: %v (id: %v)", event.Code, event.ID)
}

// reset the database completely:
func (db *Database) resetDB() error {

	sqlList := [4]string{
		db.models["event"].deleteQuery,
		db.models["project"].deleteQuery,
		fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s'", db.models["event"].table),
		fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s'", db.models["project"].table),
	}

	for _, sql := range sqlList {

		// Create a prepared SQL statement
		stmt, err := db.conn.Prepare(sql)
		// Exit if we get an error
		if err != nil {
			return err
		}
		// Make sure to cleanup after the program exits
		defer stmt.Close()

		// Execute
		_, err = stmt.Exec()

		// Exit if we get an error
		if err != nil {
			return err
		}

	}

	return nil
}
