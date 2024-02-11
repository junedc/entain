package db

const (
	eventsList = "list"
	event      = "get"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT *				
			FROM events
		`,
		event: `
			SELECT 
				*
			FROM events
			WHERE id = ?
		`,
	}
}
