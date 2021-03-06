package jsonpath

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/stretchr/testify/assert"
)

// Test that OptionFunction func conforms to Option interface
var _ Option = OptionFunction(func(selector *Selector) error { return nil })

func Test_OptionFunction(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		called := false
		option := OptionFunction(func(selector *Selector) error {
			called = true
			return nil
		})
		err := option.Apply(nil)
		assert.Nil(t, err)
		assert.True(t, called)
	})
	t.Run("error", func(t *testing.T) {
		called := false
		option := OptionFunction(func(selector *Selector) error {
			called = true
			return fmt.Errorf("fail")
		})
		err := option.Apply(nil)
		assert.EqualError(t, err, "fail")
		assert.True(t, called)
	})
}

func Test_ScriptEngine(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		engine := &testScriptEngine{value: 1}
		option := ScriptEngine(engine)
		selector := &Selector{}

		err := option.Apply(selector)
		assert.Nil(t, err)
		assert.Equal(t, engine, selector.engine)
	})
	t.Run("second", func(t *testing.T) {
		engine1 := &testScriptEngine{value: 1}
		engine2 := &testScriptEngine{value: 2}

		option1 := ScriptEngine(engine1)
		option2 := ScriptEngine(engine2)

		selector := &Selector{}

		err := option1.Apply(selector)
		assert.Nil(t, err)

		err = option2.Apply(selector)
		assert.Nil(t, err)

		assert.Equal(t, engine1, selector.engine)
		assert.NotEqual(t, engine2, selector.engine)
	})
}

func Test_QueryOptions(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		input := &option.QueryOptions{AllowMapReferenceByIndex: true, AllowStringReferenceByIndex: true}
		option := QueryOptions(input)
		selector := &Selector{}

		err := option.Apply(selector)
		assert.Nil(t, err)
		assert.Equal(t, input, selector.Options)
	})
	t.Run("second", func(t *testing.T) {
		input1 := &option.QueryOptions{AllowMapReferenceByIndex: true, AllowStringReferenceByIndex: true}
		input2 := &option.QueryOptions{AllowMapReferenceByIndex: false, AllowStringReferenceByIndex: false}

		option1 := QueryOptions(input1)
		option2 := QueryOptions(input2)

		selector := &Selector{}

		err := option1.Apply(selector)
		assert.Nil(t, err)

		err = option2.Apply(selector)
		assert.EqualError(t, err, "option already set")

		assert.Equal(t, input1, selector.Options)
		assert.NotEqual(t, input2, selector.Options)
	})
}
