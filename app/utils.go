package app

type ActionError struct {
	message  string
	notFound bool
}

func (err *ActionError) Error() string {
	return err.message
}

func (err *ActionError) IsNotFound() bool {
	return err.notFound
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
