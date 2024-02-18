package optional

import (
	"errors"
	"testing"

	"github.com/duke-git/lancet/v2/internal"
)

func TestEmpty(t *testing.T) {
	assert := internal.NewAssert(t, "TestEmpty")
	opt := Empty[int]()

	assert.ShouldBeTrue(opt.IsEmpty())
}

func TestOf(t *testing.T) {
	assert := internal.NewAssert(t, "TestOf")
	value := 42
	opt := Of(value)

	assert.ShouldBeTrue(opt.IsPresent())
	assert.Equal(opt.Get(), value)
}

func TestOfNullable(t *testing.T) {
	assert := internal.NewAssert(t, "TestOfNullable")
	var value *int = nil
	opt := OfNullable(value)

	assert.ShouldBeFalse(opt.IsPresent())

	value = new(int)
	*value = 42
	opt = OfNullable(value)

	assert.ShouldBeTrue(opt.IsPresent())
}

func TestOrElse(t *testing.T) {
	assert := internal.NewAssert(t, "TestOrElse")
	optEmpty := Empty[int]()
	defaultValue := 100

	val := optEmpty.OrElse(defaultValue)
	assert.Equal(val, defaultValue)

	optWithValue := Of(42)
	val = optWithValue.OrElse(defaultValue)
	assert.Equal(val, 42)
}

func TestOrElseGetHappyPath(t *testing.T) {
	assert := internal.NewAssert(t, "TestOrElseGetHappyPath")
	optWithValue := Of(42)
	supplier := func() int { return 100 }

	val := optWithValue.OrElseGet(supplier)
	assert.Equal(val, 42)
}

func TestOrElseGet(t *testing.T) {
	assert := internal.NewAssert(t, "TestOrElseGet")
	optEmpty := Empty[int]()
	supplier := func() int { return 100 }

	val := optEmpty.OrElseGet(supplier)
	assert.Equal(val, supplier())
}

func TestOrElseThrow(t *testing.T) {
	assert := internal.NewAssert(t, "TestOrElseThrow")
	optEmpty := Empty[int]()
	_, err := optEmpty.OrElseThrow(func() error { return errors.New("no value") })

	assert.Equal(err.Error(), "no value")

	optWithValue := Of(42)
	val, err := optWithValue.OrElseThrow(func() error { return errors.New("no value") })

	assert.IsNil(err)
	assert.Equal(val, 42)
}

func TestIfPresent(t *testing.T) {
	assert := internal.NewAssert(t, "TestIfPresent")
	called := false
	action := func(value int) { called = true }

	optEmpty := Empty[int]()
	optEmpty.IfPresent(action)

	assert.ShouldBeFalse(called)

	called = false // Reset for next test
	optWithValue := Of(42)
	optWithValue.IfPresent(action)

	assert.ShouldBeTrue(called)
}

func TestIfPresentOrElse(t *testing.T) {
	assert := internal.NewAssert(t, "TestIfPresentOrElse")

	// Test when value is present
	calledWithValue := false
	valueAction := func(value int) { calledWithValue = true }
	emptyAction := func() { t.Errorf("Empty action should not be called when value is present") }

	optWithValue := Of(42)
	optWithValue.IfPresentOrElse(valueAction, emptyAction)

	assert.ShouldBeTrue(calledWithValue)

	// Test when value is not present
	calledWithEmpty := false
	valueAction = func(value int) { t.Errorf("Value action should not be called when value is not present") }
	emptyAction = func() { calledWithEmpty = true }

	optEmpty := Empty[int]()
	optEmpty.IfPresentOrElse(valueAction, emptyAction)

	assert.ShouldBeTrue(calledWithEmpty)
}

func TestGetWithPanicStandard(t *testing.T) {
	assert := internal.NewAssert(t, "TestGetWithPanicStandard")

	// Test when value is present
	optWithValue := Of(42)
	func() {
		defer func() {
			r := recover()
			assert.IsNil(r)
		}()
		val := optWithValue.Get()
		if val != 42 {
			t.Errorf("Expected Get to return 42, got %v", val)
		}
	}()

	// Test when value is not present
	optEmpty := Empty[int]()
	func() {
		defer func() {
			r := recover()
			assert.IsNotNil(r)
		}()
		_ = optEmpty.Get()
	}()
}
