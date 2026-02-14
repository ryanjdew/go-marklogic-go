# UTCP Code-Mode Workspace Configuration

This folder contains manual definitions and examples for running a UTCP Code Mode MCP server with this workspace.

Overview:
- `.utcp_config.json` at the project root defines the code-mode configuration (what manuals to register and how to find MarkLogic endpoints using environment variables).
- `utcp/manuals/` contains simple example manual JSON files used by the code-mode bridge to register tools.
 - `utcp/manuals/` contains simple example manual JSON files used by the code-mode bridge to register tools, including a `go` manual to support Go developer workflows.
- You can add additional manuals, OpenAPI metadata, or change the call templates in `.utcp_config.json` to register more tools.

Examples:
- Run the MCP server:
  - Quick: `npx @utcp/code-mode-mcp`
  - Using helper script: `mcp/start-code-mode-mcp.sh` (Linux/macOS) or `mcp/start-code-mode-mcp.ps1` (Windows)

  CLI Tools Security Note:
  - The Go tools are registered as `cli` call templates in `.utcp_config.json` (go.build, go.test, go.mod.tidy, go.lint).
  - The code-mode MCP bridge requires an explicit opt-in to enable CLI tool execution (`registerCli()`), and CLI will run commands on your system; enable only on secure/trusted environments.

  Example: call Go workflows from TypeScript using Code Mode

  ```typescript
  const { result, logs } = await client.callToolChain(`
    // Build the project
    await go.build({ packages: "./...", output: "./bin/app" });
    // Run tests
    const testResults = await go.test({ packages: "./...", flags: "-v" });
    // Run a quick lint pass
    await go.lint({ packages: "./..." });
    return { tests: testResults };
  `);
  console.log('Execution logs:', logs);
  ```
