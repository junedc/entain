package db

import (
	"math/rand"
	"time"

	"syreclabs.com/go/faker"
)

func (r *eventsRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, category TEXT, status TEXT,  advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	statuses := []string{"ABANDON", "CLOSED", "SETUP", "READY", "PAUSED", "PREMATCH", "INPLAY", "MATCHOVER"}
	games := []string{"AFL", "BASKETBALL", "BASEBALL", "GOLF", "RUGBY LEAGUE", "SOCCER", "TENNIS"}
	for i := 1; i <= 100; i++ {

		rand.Seed(time.Now().UnixNano()) // seed or it will be set to 1

		status := statuses[rand.Intn(len(statuses))]
		game := games[rand.Intn(len(games))]
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO events(id, name, category, status, advertised_start_time) VALUES (?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Name()+" @ "+faker.Team().Name(),
				game,
				status,
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 10)).Format(time.RFC3339),
			)
		}
	}

	return err
}
