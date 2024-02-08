package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(request *racing.ListRacesRequest) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(request *racing.ListRacesRequest) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, request)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, request *racing.ListRacesRequest) (string, []interface{}) {
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
		tableFields := []string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}
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

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		/*	Add status to racing results
			All races that have an advertised_start_time
			in the past should reflect CLOSED.
		*/
		if advertisedStart.After(time.Now()) {
			race.Status = "OPEN"
		} else {
			race.Status = "CLOSED"
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
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
