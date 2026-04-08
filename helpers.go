package main

import (
	"database/sql"
	"time"
)

func convertToNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func parsePublishedAt(pubDate string) sql.NullTime {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, pubDate)
		if err == nil {
			return sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
	}

	return sql.NullTime{
		Valid: false,
	}
}
