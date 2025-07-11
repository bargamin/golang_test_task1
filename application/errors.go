package application

type HTTPError struct {
	Status      int
	Description string
}

func (e *HTTPError) Error() string { return e.Description }
