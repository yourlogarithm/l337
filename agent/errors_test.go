package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrBuilderParams(t *testing.T) {
	t.Run("error formatting", func(t *testing.T) {
		err := ErrBuilderParams{
			Param: "testParam",
			Msg:   "test error message",
		}
		
		expected := "error building agent: testParam: test error message"
		assert.Equal(t, expected, err.Error())
	})
	
	t.Run("error with empty param", func(t *testing.T) {
		err := ErrBuilderParams{
			Param: "",
			Msg:   "test error message",
		}
		
		expected := "error building agent: : test error message"
		assert.Equal(t, expected, err.Error())
	})
	
	t.Run("error with empty message", func(t *testing.T) {
		err := ErrBuilderParams{
			Param: "testParam",
			Msg:   "",
		}
		
		expected := "error building agent: testParam: "
		assert.Equal(t, expected, err.Error())
	})
	
	t.Run("error implements error interface", func(t *testing.T) {
		var err error = ErrBuilderParams{
			Param: "test",
			Msg:   "test message",
		}
		
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error building agent")
	})
}

func TestErrModelResponse(t *testing.T) {
	t.Run("error formatting", func(t *testing.T) {
		err := ErrModelResponse{
			Msg: "test model error",
		}
		
		expected := "error processing model response: test model error"
		assert.Equal(t, expected, err.Error())
	})
	
	t.Run("error with empty message", func(t *testing.T) {
		err := ErrModelResponse{
			Msg: "",
		}
		
		expected := "error processing model response: "
		assert.Equal(t, expected, err.Error())
	})
	
	t.Run("error implements error interface", func(t *testing.T) {
		var err error = ErrModelResponse{
			Msg: "test message",
		}
		
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error processing model response")
	})
}

func TestErrorsInRealScenarios(t *testing.T) {
	t.Run("ErrBuilderParams used in New function", func(t *testing.T) {
		_, err := New(nil)
		
		assert.Error(t, err)
		
		// Check that it's the correct error type
		var builderErr ErrBuilderParams
		assert.ErrorAs(t, err, &builderErr)
		assert.Equal(t, "model", builderErr.Param)
		assert.Equal(t, "model is required", builderErr.Msg)
	})
	
	t.Run("ErrBuilderParams used in BuildRunResponse", func(t *testing.T) {
		_, err := BuildRunResponse() // No messages provided
		
		assert.Error(t, err)
		
		// Check that it's the correct error type
		var builderErr ErrBuilderParams
		assert.ErrorAs(t, err, &builderErr)
		assert.Equal(t, "Messages", builderErr.Param)
		assert.Equal(t, "at least one message is required", builderErr.Msg)
	})
}