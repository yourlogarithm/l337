# l337 Architecture

## Project Structure

The l337 project follows Go best practices for package organization. Here's an overview of the structure:

```
l337/
├── agent/              # Core agent implementation
│   ├── remote/         # Remote agent functionality
│   │   ├── client/     # HTTP client for remote agents
│   │   └── server.go   # HTTP server for agents
│   ├── agent.go        # Main agent struct and logic
│   ├── builder.go      # Builder pattern for agent creation
│   ├── interface.go    # AgentImpl interface
│   ├── run.go          # Agent execution logic
│   └── ...             # Other agent-related files
├── chat/               # Chat-related types and utilities
│   ├── message.go      # Message types
│   ├── content.go      # Content types (text, image, audio)
│   ├── options.go      # Chat options and parameters
│   └── ...             # Other chat-related files
├── docs/               # Documentation assets
│   └── logo.svg        # Project logo
├── examples/           # Example usage (runnable examples)
│   ├── agent.go        # Simple agent example
│   ├── team.go         # Multi-agent team example
│   └── tools.go        # Tool usage example
├── internal/           # Internal packages (not for external use)
│   └── logging/        # Logging utilities
├── metrics/            # Metrics collection and reporting
│   └── metrics.go      # Metrics types and utilities
├── provider/           # LLM provider interfaces and implementations
│   ├── base.go         # Provider interfaces
│   ├── types.go        # Common provider types
│   ├── ollama/         # Ollama provider implementation
│   └── openai/         # OpenAI provider implementation
└── tools/              # Tool definitions and integrations
    ├── tool.go         # Tool interface and creation
    ├── toolkit.go      # Tool management utilities
    └── mcp.go          # Model Context Protocol integration
```

## Package Overview

### Core Packages

#### `agent`
The agent package provides the main functionality for creating and managing AI agents. It includes:
- Agent builder with options pattern
- Agent execution (run and streaming)
- System message generation
- Subordinate agent management (for multi-agent systems)
- Remote agent capabilities (HTTP-based)

#### `chat`
Contains all chat-related types and utilities:
- Message structures (user, assistant, system, tool)
- Content types (text, image, audio)
- Chat parameters and options
- Run response handling

#### `provider`
Defines interfaces for LLM providers and contains implementations:
- Base provider interface (`ModelImpl`)
- Provider-agnostic request/response types
- Streaming support
- Provider-specific implementations (OpenAI, Ollama)

#### `tools`
Tool system for extending agent capabilities:
- Type-safe tool definitions
- JSON schema generation
- Model Context Protocol (MCP) integration
- Tool execution and management

#### `metrics`
Metrics collection for monitoring agent and model performance:
- Token usage tracking
- Timing information
- Provider-specific metrics

### Supporting Packages

#### `internal/logging`
Internal logging utilities using Go's `log/slog` package. Not intended for external use.

#### `examples`
Runnable examples demonstrating library usage:
- Simple agent creation
- Multi-agent teams
- Tool integration

## Design Principles

1. **Separation of Concerns**: Each package has a clear, focused responsibility
2. **Provider Agnostic**: Core abstractions work with any LLM provider
3. **Type Safety**: Strong typing with Go's type system
4. **Builder Pattern**: Flexible agent configuration using functional options
5. **Streaming Support**: First-class support for streaming responses
6. **Extensibility**: Easy to add new providers, tools, and capabilities

## Package Dependencies

```
examples → agent → provider → chat
                 ↓           ↓
               tools      metrics
                 ↑
              internal/logging
```

- `examples` depends on `agent` and demonstrates usage
- `agent` coordinates between `provider`, `chat`, and `tools`
- `provider` implementations use `chat` types
- All packages can use `internal/logging`
- `metrics` is used by providers to report performance data

## Extension Points

### Adding a New Provider

1. Implement the `provider.ModelImpl` interface
2. Create a new package under `provider/`
3. Add conversion functions between provider's types and `chat` types

### Adding a New Tool

1. Use `tools.New()` for simple tools
2. Use `tools.NewWithArgs[T]()` for tools with typed parameters
3. Register with agent using `agent.WithTool()`

### Creating a Remote Agent

1. Wrap an agent with `remote.AgentServer`
2. Serve over HTTP
3. Connect using `remote/client.RemoteAgent`
