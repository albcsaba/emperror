package emperror

import (
	"fmt"
	"testing"

	"emperror.dev/errors"
)

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected to recover from a panic, but nothing panicked")
		}

		err, ok := r.(error)
		if !ok {
			t.Fatal("expected to recover an error from a panic")
		}

		if err == nil {
			t.Fatal("expected to the recovered error to be an error, received nil")
		}

		var stt stackTracer
		if !errors.As(err, &stt) {
			t.Fatal("error is expected to carry a stack trace")
		}

		st := stt.StackTrace()

		if got, want := fmt.Sprintf("%n", st[0]), "TestPanic"; got != want {
			t.Errorf("function name does not match the expected one\nactual:   %s\nexpected: %s", got, want)
		}

		if got, want := fmt.Sprintf("%s", st[0]), "panic_test.go"; got != want {
			t.Errorf("file name does not match the expected one\nactual:   %s\nexpected: %s", got, want)
		}

		if got, want := fmt.Sprintf("%d", st[0]), "46"; got != want {
			t.Errorf("line number does not match the expected one\nactual:   %s\nexpected: %s", got, want)
		}
	}()

	Panic(errors.New("error"))
}

func TestPanic_NoError(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Fatalf("unexpected panic, received: %v", r)
		}
	}()

	Panic(nil)
}

func createRecoverFunc(p interface{}) func() error {
	return func() (err error) {
		defer func() {
			err = Recover(recover())
		}()

		panic(p)
	}
}

func assertRecoveredError(t *testing.T, err error, msg string) {
	t.Helper()

	var stt stackTracer
	if !errors.As(err, &stt) {
		t.Fatal("error is expected to carry a stack trace")
	}

	st := stt.StackTrace()

	if got, want := fmt.Sprintf("%n", st[0]), "createRecoverFunc.func1"; got != want {
		t.Errorf("function name does not match the expected one\nactual:   %s\nexpected: %s", got, want)
	}

	if got, want := fmt.Sprintf("%s", st[0]), "panic_test.go"; got != want {
		t.Errorf("file name does not match the expected one\nactual:   %s\nexpected: %s", got, want)
	}

	if got, want := fmt.Sprintf("%d", st[0]), "66"; got != want {
		t.Errorf("line number does not match the expected one\nactual:   %s\nexpected: %s", got, want)
	}

	if got, want := err.Error(), msg; got != want {
		t.Errorf("error message does not match the expected one\nactual:   %s\nexpected: %s", got, want)
	}
}

func TestRecover_ErrorPanic(t *testing.T) {
	err := fmt.Errorf("internal error")

	f := createRecoverFunc(err)

	v := f()

	assertRecoveredError(t, v, "internal error")
}

func TestRecover_StringPanic(t *testing.T) {
	f := createRecoverFunc("internal error")

	v := f()

	assertRecoveredError(t, v, "internal error")
}

func TestRecover_AnyPanic(t *testing.T) {
	f := createRecoverFunc(123)

	v := f()

	assertRecoveredError(t, v, "unknown panic, received: 123")
}

func TestRecover_Nil(t *testing.T) {
	f := createRecoverFunc(nil)

	v := f()

	if got, want := v, error(nil); got != want {
		t.Errorf("the recovered value is expected to be nil\nactual: %v", got)
	}
}

func TestHandleRecover(t *testing.T) {
	err := errors.New("error")

	defer HandleRecover(ErrorHandlerFunc(func(err error) {
		if got, want := err.Error(), "error"; got != want {
			t.Errorf("error does not match the expected one\nactual:   %s\nexpected: %s", got, want)
		}
	}))

	panic(err)
}
