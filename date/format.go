package date

import "strings"

func Format(format string) string {
	format = strings.Replace(format, "MMMM", "January", -1)
	format = strings.Replace(format, "MMM", "Jan", -1)
	format = strings.Replace(format, "MM", "01", -1)
	format = strings.Replace(format, "M", "1", -1)

	format = strings.Replace(format, "ddd", "_2", -1)
	format = strings.Replace(format, "dd", "02", -1)
	format = strings.Replace(format, "d", "2", -1)

	format = strings.Replace(format, "HH", "15", -1)
	format = strings.Replace(format, "hh", "03", -1)
	format = strings.Replace(format, "hh", "3", -1)

	format = strings.Replace(format, "mm", "04", -1)
	format = strings.Replace(format, "m", "4", -1)

	format = strings.Replace(format, "ss", "05", -1)
	format = strings.Replace(format, "s", "5", -1)

	format = strings.Replace(format, "yyyy", "2006", -1)
	format = strings.Replace(format, "yy", "06", -1)
	format = strings.Replace(format, "y", "06", -1)

	format = strings.Replace(format, "SSS", "000", -1)

	format = strings.Replace(format, "a", "pm", -1)
	format = strings.Replace(format, "aa", "PM", -1)

	format = strings.Replace(format, "ZZ", "-0700", -1)
	format = strings.Replace(format, "Z", "-07", -1)

	format = strings.Replace(format, "zz:zz", "Z07:00", -1)
	format = strings.Replace(format, "zzzz", "Z0700", -1)
	format = strings.Replace(format, "z", "MST", -1)

	format = strings.Replace(format, "EEEE", "Monday", -1)
	format = strings.Replace(format, "E", "Mon", -1)

	return format
}

const (
	ISO_DATE = "2006-01-02" // Format("yyyy-MM-dd")
	ISO_TIME = "15:04:05" //Format("yyyy-MM-ddTHH:mm:ss.SSSZZ")
	ISO_TIMESTAMP = "2006-01-02T15:04:05.000-0700" //Format("yyyy-MM-ddTHH:mm:ss.SSSZZ")
	ISO_DATETIME = "2006-01-02 15:04:05" //Format("yyyy-MM-ddTHH:mm:ss")
	ISO_ZONEDDATETIME = "2006-01-02 15:04:05-07" //Format("yyyy-MM-ddTHH:mm:ssZ")
)
