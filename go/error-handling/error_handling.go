package erratum

/*Use opens and uses a resource with error handling*/
func Use(o ResourceOpener, input string) (retval error) {
	var res Resource
	var err error
	for res, err = o(); nil != err; res, err = o() {
		if _, ok := err.(TransientError); ok {
			// sleep a bit then
			continue
		} else {
			return err
		}
	}
	// opened ok
	defer res.Close()
	defer func() {
		switch v := recover().(type) {
		case FrobError:
			res.Defrob(v.defrobTag)
			retval = v.inner
		case error:
			retval = v
		default:
			// ignore
		}
	}()
	// Finally call the target
	res.Frob(input)
	retval = nil
	return // retval
}
