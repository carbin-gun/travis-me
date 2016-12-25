package try_test

import (
	"testing"
	"log"
	"fmt"

	"github.com/carbin-gun/travis-me/toolbox/try"
	"errors"
)

func TestRetryExample(t *testing.T) {
	try.MaxTries = 20
	SomeFunction := func() (string, error) {
		return "", nil
	}
	var value string
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		value, err = SomeFunction()
		return attempt < 5, err
	})
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func TestRetryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("somehing badly happened")
	}
	var value string
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("received panic:%v", r))
			}
		}()
		value, err = SomeFunction()
		return
	})
	//every time ,it retried.but every time until the last one,there is always an error received from panic.so there is an error
	//from panic.
	if err != nil {
		//	log.Fatalln("error:", err)
	}
}

func TestTrySuccess(t *testing.T) {
	var callCount = 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})
	if err != nil {
		log.Fatalln("error:", err)
	}
	if callCount != 1 {
		log.Fatalln("error callCount!=1")
	}
}

func TestTryFail(t *testing.T) {
	var callCount = 0
	theErr := errors.New("something went wrong!")
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})
	if err != theErr {
		log.Fatalln("error:", err)
	}
	if callCount != 5 {
		log.Fatalln("error callCount!=5")
	}
}

func TestTryLimit(t *testing.T) {
	err := try.Do(func(attempt int) (bool, error) {
		return true, errors.New("Error")
	})
	if !try.IsMaxTry(err) {
		log.Fatalln("error not equal try.MaxRetries")
	}
}

func TestTryPanics(t *testing.T) {
	theErr := errors.New("something went wrong")
	var callCount = 0
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("received panic:%v", r))
			}
		}()
		callCount ++
		if attempt > 2 {
			panic("I don't like three")
		}
		err = theErr //会使用panic而不会使用到theErr
		return

	})
	if err.Error() != "received panic:I don't like three" {
		log.Fatalln("err.Error()!='received panic:I don't like three'")
	}
	if callCount != 5 {
		log.Fatalln("callCount!=5")
	}
}
