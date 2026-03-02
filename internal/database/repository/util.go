package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func toString(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func fromString(s pgtype.Text) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func toInt(i int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(i), Valid: i != 0}
}

func toBool(b *bool) pgtype.Bool {
	bb := false
	if b != nil {
		bb = *b
	}

	return pgtype.Bool{Bool: bb, Valid: b != nil}
}

// Uncomment when used
//
//	func toTime(t time.Time) pgtype.Timestamptz {
//		return pgtype.Timestamptz{Time: t, Valid: !t.IsZero()}
//	}
func fromTime(t pgtype.Timestamptz) time.Time {
	if t.Valid {
		return t.Time
	}

	return time.Time{}
}

// Uncomment when used
// func toDuration(d time.Duration) pgtype.Int8 {
// 	return pgtype.Int8{Int64: d.Nanoseconds(), Valid: d.Nanoseconds() > 0}
// }
