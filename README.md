# Wrapping errors

Here are some ways of using the new Go errors wrappers.

## Sentinels

Sentinels are basic errors with nothing more than a type and a string value.

You can create a single `Sentinel` type and an `Is` method to handle all sentinels in your program:

```
type Sentinel string

func (e Sentinel) Error() string {
	return string(e)
}

func (e Sentinel) Is(err error) bool {
	sentinel, ok := err.(Sentinel)
	if !ok {
		return false
	}
	return sentinel == e
}
```

Create sentinels like this:

```
const ErrBroken = Sentinel("a specific thing broke")
```

So you can do this:

```
if errors.Is(err, ErrBroken) ...
```

If any error in the chain starting with `err` is an `ErrBroken`, then `Is` will be true.

You don't have to do anything special to use Sentinels in this way.

Your application or package can easily set up lots of sentinels, each of which can be distinguished with `errors.Is`:

```
const (
    ErrBroken = Sentinel("broke this way")
    ErrBusted = Sentinel("broke that way")
)
```

## Custom Error Types

The next level of complexity is wrapping an error in a custom type to "tag" an error with boolean attributes.
For example, you could have a custom error type that indicates that an error is temporary, or that it was the result of a timeout.

These kinds of error type require a constructor and a couple of methods.
But when these are written, the new error types fit right into the new techniques.

Here is an example of using a new type to indicate temporary errors.

First we create the new type, and all we have to do is embed `error`:

```
type TemporaryError struct {
    error
}
```

Now we build a constructor:

```
func NewTemporary(err error) error {
    return &TemporaryError{fmt.Errorf("temporary: %w", err)}
}
```
(You have to use `fmt.Errorf` to wrap the `err` because the returned value from `fmt.Errorf` has both `Error` and `Unwrap` methods at the right "level".)

Now we add `Is` and `Unwrap` methods on our new error type so that `errors.Is` and `errors.Unwrap` work right:

```
func (e *TemporaryError) Is(err error) bool {
	_, ok := err.(*TemporaryError)
	return ok
}

func (e *TemporaryError) Unwrap() error {
	return errors.Unwrap(e.error)
}
```

Now your new error type is a fully paid up member of the club.
You can do this to any error chain:

```
if errors.Is(err, &TemporaryError{}) ...
```

Or you can write a shortcut:

```
func IsTemporary(err error) bool {
	return errors.Is(err, &TemporaryError{})
}
```

and do this:

```
if IsTemporary(err) ...
```

## Complex error types

Sometimes you want to add more context to an error.
Suppose you want to set status codes that can be found at the top of your program.

First create a new error type by embedding `error` and adding any additional fields you need:

```
type ErrStatusCode struct {
	error
	statusCode int
}
```

Create a constructor:

```
func WithStatusCode(err error, code int) error {
	return &ErrStatusCode{
		error:      fmt.Errorf("code %d: %w", code, err),
		statusCode: code,
	}
}
```

And add the standard `Is` and `Unwrap` methods:

```
func (e *ErrStatusCode) Is(err error) bool {
	_, ok := err.(*ErrStatusCode)
	return ok
}

func (e *ErrStatusCode) Unwrap() error {
	return errors.Unwrap(e.error)
}
```

In addition, write any methods that make sense for your new type:

```
func (e *ErrStatusCode) StatusCode() int {
	return e.statusCode
}
```

Now at the top of your program you can grab an error's status code, if it has one, like this:

```
var errStatusCode *ErrStatusCode
var code int
if errors.As(err, &errStatusCode) {
    code = errStatusCode.StatusCode()
}
```

It doesn't matter how many times the `ErrStatusCode` has been wrapped; `As` will find it.
## First cause

With standardized `Unwrap` methods, you can easily find the original error in the chain:

```
func Cause(err error) error {
	type wrapper interface {
		Unwrap() error
	}
	for err != nil {
		cause, ok := err.(wrapper)
		if !ok {
			break
		}
		err = cause.Unwrap()
	}
	return err
}
```

## Visualizing the error chain

If you ever wonder what the error chain looks like, you can do this:

```
type Printer func(format string, args ...interface{})

func Chain(printf Printer, err error) {
	printf("error chain:\n")
	for err != nil {
		printf("\t%T %v\n", err, err)
		err = errors.Unwrap(err)
	}
}
```

Here is what the output can look like if you call it like this:

```
Chain(fmt.Printf, err)
```

```
error chain:
    *fmt.wrapError some annotation: code 400: temporary: original error
    *qerrors.ErrStatusCode code 400: temporary: original error
    *qerrors.TemporaryError temporary: original error
    qerrors.Sentinel original error
```

## Annotations

With the new `fmt.Errorf`, you can annotate errors without getting complicated, and without losing the benefits of `Is` and `As`, no matter how many annotations are added:

```
if err != nil {
    return fmt.Errorf("could get details for user %s: %w", user_id, err)
}
```

The original error is wrapped with an annotation, and no information is lost.
The entire chain can still be unwrapped at the top of the program.