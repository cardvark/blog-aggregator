package database

import (
	"database/sql"
	"time"
)

func GetNullTime(timeval time.Time) sql.NullTime {
	var nullTimeVal sql.NullTime
	nullTimeVal.Time = timeval
	nullTimeVal.Valid = true
	return nullTimeVal
}

func GetNullText(textval string) sql.NullString {
	var nullStringVal sql.NullString
	nullStringVal.String = textval
	nullStringVal.Valid = true
	return nullStringVal
}
