package examples

import (
	"context"
	"fmt"
	"io"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/providers/openai"
	"github.com/yourlogarithm/l337/types"
)

func StreamingExample() {
	model := openai.NewModel("gpt-4o")
	agent, err := agent.New(
		model,
		agent.WithName("streamer_jedi"),
	)
	if err != nil {
		panic(err)
	}

	run, channel, err := agent.RunStreamingWithParams(
		context.Background(),
		0, // bufferSize as in go channels
		types.WithTextMessage(types.RoleUser, "What was your greatest failure, and what did you learn from it?"),
	)
	if err != nil {
		panic(err)
	}

	for {
		chunk, err := channel.Next()
		if err == io.EOF {
			// stream ended naturally
			break
		} else if err != nil {
			// something bad happened
			panic(err)
		}
		// do something with the chunk
		fmt.Print(chunk.Content)
	}

	// `run` is considered complete only after `io.EOF` is received, otherwise it is incomplete
	// If you don't need the chunks anymore - can call `channel.Drain()` before using `run`
	// `channel.Drain()` will wait for stream completion, discarding all chunks:
	// ```go
	// n_chunks := channel.Drain()
	// ```

	fmt.Println(run.Content())
}
