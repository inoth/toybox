package httpserver

var (
	defaultOption = option{
		ReadTimeout:    10,
		WriteTimeout:   10,
		MaxHeaderBytes: 20,
		Port:           ":9060",
		handles:        make([]Handler, 0),
	}
)
