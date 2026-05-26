package levelfour

import "time"

// ParseDate parses a date string in "2006-01-02" format.
// Unlike MustParseDate, it returns an error instead of panicking.
func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

// ParseDateTime parses a datetime string in RFC3339 format.
// Unlike MustParseDateTime, it returns an error instead of panicking.
func ParseDateTime(datetime string) (time.Time, error) {
	return time.Parse(time.RFC3339, datetime)
}

// ParseDatePtr parses an optional date string. Returns nil if the
// input pointer is nil. Useful for API response fields that are
// nullable date values.
func ParseDatePtr(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	t, err := ParseDate(*s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ParseDateTimePtr parses an optional datetime string. Returns nil
// if the input pointer is nil. Useful for API response fields that
// are nullable datetime values.
func ParseDateTimePtr(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	t, err := ParseDateTime(*s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
