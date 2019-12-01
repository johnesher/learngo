package erratum


// open and use a resource
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
	defer func () {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case FrobError:
				res.Defrob(v.defrobTag)
				retval = v.inner
			default:
				retval = v.(error)
			}
		}
	}()
	res.Frob(input)
	return nil
}
