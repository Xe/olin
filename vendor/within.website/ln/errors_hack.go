package ln

type wrapper interface {
	Unwrap() error
}

func unwrap(err error) error {
	u, ok := err.(wrapper)
	if !ok {
		return nil
	}

	return u.Unwrap()
}

func doGoError(err error, f F) {
	uw := unwrap(err)
	if uw != nil {
		f["wrapped_err"] = uw
	}
}
