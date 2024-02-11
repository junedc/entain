package db

import (
	"cruzes.co/junedc/entain/sport/proto/sport"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
)

// eventsRepo provides repository access to races.
type EventsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of races.
	List(request *sport.ListEventsRequest) ([]*sport.Event, error)

	Get(request *sport.GetEventRequest) (*sport.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NeweventsRepo creates a new races repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *eventsRepo) List(request *sport.ListEventsRequest) ([]*sport.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = r.applyFilter(query, request)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

// Get - get race by id
func (r *eventsRepo) Get(request *sport.GetEventRequest) (*sport.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[event]

	args = append(args, request.Id)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	//reuse scanEvents
	races, err := r.scanEvents(rows)

	if len(races) > 0 {
		//return the first record
		return races[0], nil
	} else {
		//if race is not found return an empty race
		//just like declaring a string which has a default value of ''
		//TODO:: might change to nil or return an error depends on detailed specs
		return new(sport.Event), nil
	}

}

func (r *eventsRepo) applyFilter(query string, request *sport.ListEventsRequest) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if request == nil {
		return query, args
	}

	if len(request.Filter.MeetingIds) > 0 {
		//count the number of meeting ids and create a SQL query with parameter denoted by ?
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(request.Filter.MeetingIds)-1)+"?)")

		//pass the actual value of the meeting ids as argument
		for _, meetingID := range request.Filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	//call ListRaces asking for races that are visible only
	if request.Filter.Visible != nil && request.Filter.GetVisible() {
		clauses = append(clauses, "visible = true")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	//check if sort_by parameter is supplied
	if request.SortBy != nil && request.GetSortBy() != "" {
		sortByFields := strings.Split(request.GetSortBy(), ",")
		tableFields := []string{"id", "name", "category", "status", "advertised_start_time"}
		validSortByFields := checkSortByFields(sortByFields, tableFields)
		if len(validSortByFields) == 0 {
			query += " ORDER BY advertised_start_time ASC"
		} else {
			//generate SQL command from valid sort_by fields
			query += " ORDER BY " + strings.Join(validSortByFields, ",")
		}
		return query, args
	}

	//add default ordering
	query += " ORDER BY advertised_start_time ASC"

	return query, args
}

func (m *eventsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sport.Event, error) {
	var races []*sport.Event

	for rows.Next() {
		var sportEvent sport.Event
		var advertisedStart time.Time

		if err := rows.Scan(&sportEvent.Id, &sportEvent.Name, &sportEvent.Category, &sportEvent.Status, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		sportEvent.AdvertisedStartTime = ts

		races = append(races, &sportEvent)
	}

	return races, nil
}

// return valid sort_by fields - fields that exist in table
// copied from https://go.dev/play/p/eGGcyIlZD6y
func checkSortByFields(tableFields []string, orderFields []string) []string {
	out := []string{}
	bucket := map[string]bool{}

	for _, i := range tableFields {
		for _, j := range orderFields {
			if i == j && !bucket[i] {
				out = append(out, i)
				bucket[i] = true
			}
		}
	}

	return out
}
