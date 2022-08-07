package basal

func Sprintf(format string, a ...interface{}) string

func NewError(format string, a ...interface{}) error

func ToError(r interface{}) (err error)
