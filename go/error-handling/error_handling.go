package erratum

import(
	"errors"
)

// open and use a resource
func Use(o ResourceOpener, input string) error{
	res, err := o()
	if nil != err {
		// presume this means the resource was not opened
		return errors.New("fixme")
	}else{
		defer res.Close()
		res.Frob(input)
		return nil
	}
}
