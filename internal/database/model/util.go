package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func fromString(s pgtype.Text) string {
	if s.Valid {
		return s.String
	}

	return ""
}

// Uncomment when used
// func fromInt(i pgtype.Int4) int {
// 	if i.Valid {
// 		return int(i.Int32)
// 	}
//
// 	return 0
// }

// Uncomment when used
// func fromBool(b pgtype.Bool) *bool {
// 	if b.Valid {
// 		return &b.Bool
// 	}
//
// 	return nil
// }

// Uncomment when used
// func fromTime(t pgtype.Timestamptz) time.Time {
// 	if t.Valid {
// 		return t.Time
// 	}
//
// 	return time.Time{}
// }

// Uncomment when used
// func fromDuration(i pgtype.Int8) time.Duration {
// 	if i.Valid {
// 		return time.Duration(i.Int64)
// 	}
//
// 	return time.Duration(0)
// }
