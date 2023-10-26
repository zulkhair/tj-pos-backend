package dateutil

import "time"

func TimeFormat() string {
	return "2006-01-02 15:04:05"
}

func DateFormat() string {
	return "2006-01-02"
}

func DateFormatResponse() string {
	return "02-01-2006"
}

func TimeFormatResponse() string {
	return "02-01-2006 15:04:05"
}

func MonthToString(month int) string {
	switch month {
	case 1:
		return "Januari"
	case 2:
		return "Febuari"
	case 3:
		return "Maret"
	case 4:
		return "April"
	case 5:
		return "Mei"
	case 6:
		return "Juni"
	case 7:
		return "Juli"
	case 8:
		return "Agustus"
	case 9:
		return "September"
	case 10:
		return "Oktober"
	case 11:
		return "November"
	case 12:
		return "Desember"
	default:
		return ""
	}
}

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
