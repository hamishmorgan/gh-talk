# Cobra Framework Guide for gh-talk

**Practical guide to using Cobra for implementing gh-talk commands**

## Overview

Cobra is a powerful CLI framework for Go that provides command structure, flag parsing, and help generation. While `gh` CLI uses a custom framework, Cobra is the pragmatic choice for gh-talk due to its maturity, ecosystem, and development speed.

**Repository:** <https://github.com/spf13/cobra>  
**Documentation:** <https://cobra.dev>  
**Stars:** ~37k  
**Used By:** kubectl, Docker, Hugo, GitHub Actions CLI

## Why Cobra for gh-talk

**Benefits:**

- ✅ Automatic help generation
- ✅ Subcommand support built-in
- ✅ Persistent (global) flags
- ✅ Local (command-specific) flags
- ✅ Flag validation and requirements
- ✅ Shell completion (bash, zsh, fish, powershell)
- ✅ Pre/Post run hooks
- ✅ Args validation
- ✅ Huge ecosystem and examples
- ✅ Fast development (days vs weeks)

## Core Concepts

### Commands

**Command is the central structure:**

```go
type Command struct {
    Use   string   // Command usage string
    Short string   // Short description (one line)
    Long  string   // Long description (multi-line)
    RunE  func(cmd *Command, args []string) error  // Implementation
    Args  PositionalArgs  // Argument validation
}
```

**Example for gh-talk:**

```go
var replyCmd = &cobra.Command{
    Use:   "reply [thread-id] [message]",
    Short: "Reply to a review thread",
    Long: `Reply to a review thread with a message.

If thread-id is not provided, an interactive prompt will
ask you to select a thread. If message is not provided,
the --editor flag must be used.`,
    Args: cobra.RangeArgs(0, 2), // 0-2 arguments
    RunE: runReply,
}

func runReply(cmd *cobra.Command, args []string) error {
    // Implementation here
    return nil
}
```

### Flags

**Two Types:**

**Local Flags** (command-specific):

```go
replyCmd.Flags().BoolP("editor", "e", false, "Open editor for message")
replyCmd.Flags().Bool("resolve", false, "Resolve thread after replying")
```

**Persistent Flags** (inherited by subcommands):

```go
rootCmd.PersistentFlags().StringP("repo", "R", "", "Repository (OWNER/REPO)")
rootCmd.PersistentFlags().Int("pr", 0, "PR number")
```

**Getting Flag Values:**

```go
func runReply(cmd *cobra.Command, args []string) error {
    // Get flags
    useEditor, _ := cmd.Flags().GetBool("editor")
    shouldResolve, _ := cmd.Flags().GetBool("resolve")
    repo, _ := cmd.Flags().GetString("repo") // From persistent flags
    
    // Use values
    if useEditor {
        message = openEditor()
    }
    
    return nil
}
```

### Subcommands

**Parent-Child Relationship:**

```go
// Root command
var rootCmd = &cobra.Command{
    Use:   "gh-talk",
    Short: "Manage GitHub PR conversations",
}

// Parent command
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List conversations",
}

// Subcommands
var listThreadsCmd = &cobra.Command{
    Use:   "threads",
    Short: "List review threads",
    RunE:  runListThreads,
}

var listCommentsCmd = &cobra.Command{
    Use:   "comments",
    Short: "List comments",
    RunE:  runListComments,
}

// Build hierarchy
func init() {
    rootCmd.AddCommand(listCmd)
    listCmd.AddCommand(listThreadsCmd)
    listCmd.AddCommand(listCommentsCmd)
}

// Results in:
// gh-talk list threads
// gh-talk list comments
```

## gh-talk Command Structure with Cobra

### Root Command

**File: `internal/commands/root.go`**

```go
package commands

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "gh-talk",
    Short: "Manage GitHub PR and Issue conversations",
    Long: `gh-talk is a GitHub CLI extension for managing conversations
on Pull Requests and Issues from the terminal.

Features:
  • Reply to review threads
  • Add emoji reactions  
  • Resolve/unresolve threads
  • Hide comments
  • Filter conversations
  • Bulk operations`,
    SilenceUsage:  true, // Don't show usage on errors
    SilenceErrors: true, // We'll handle errors ourselves
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // Global persistent flags
    rootCmd.PersistentFlags().StringP("repo", "R", "", "Repository (OWNER/REPO)")
    rootCmd.PersistentFlags().Int("pr", 0, "PR number")
    rootCmd.PersistentFlags().Int("issue", 0, "Issue number")
    
    // Add all subcommands
    rootCmd.AddCommand(listCmd)
    rootCmd.AddCommand(replyCmd)
    rootCmd.AddCommand(resolveCmd)
    rootCmd.AddCommand(unresolveCmd)
    rootCmd.AddCommand(reactCmd)
    rootCmd.AddCommand(hideCmd)
    rootCmd.AddCommand(unhideCmd)
    rootCmd.AddCommand(showCmd)
}
```

### List Command

**File: `internal/commands/list.go`**

```go
package commands

import (
    "github.com/spf13/cobra"
    "github.com/hamishmorgan/gh-talk/internal/api"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List conversations",
}

var listThreadsCmd = &cobra.Command{
    Use:   "threads",
    Short: "List review threads",
    Long: `List review threads from a pull request.

By default, shows only unresolved threads. Use --all to see
both resolved and unresolved threads.

Examples:
  # List unresolved threads in current PR
  gh talk list threads
  
  # List all threads in PR #123
  gh talk list threads --pr 123 --all
  
  # List threads on specific file
  gh talk list threads --file src/api.go`,
    RunE: runListThreads,
}

func init() {
    listCmd.AddCommand(listThreadsCmd)
    
    // Filter flags
    listThreadsCmd.Flags().Bool("unresolved", true, "Show only unresolved threads")
    listThreadsCmd.Flags().Bool("resolved", false, "Show only resolved threads")
    listThreadsCmd.Flags().Bool("all", false, "Show all threads")
    listThreadsCmd.Flags().String("author", "", "Filter by author")
    listThreadsCmd.Flags().String("file", "", "Filter by file path")
    listThreadsCmd.Flags().String("since", "", "Show threads since date")
    
    // Output flags
    listThreadsCmd.Flags().String("format", "table", "Output format (table, json, tsv)")
    listThreadsCmd.Flags().StringSlice("json", nil, "Output JSON with specific fields")
    
    // Flag groups (mutually exclusive)
    listThreadsCmd.MarkFlagsMutuallyExclusive("unresolved", "resolved", "all")
}

func runListThreads(cmd *cobra.Command, args []string) error {
    // Get flags
    prNum, _ := cmd.Flags().GetInt("pr")
    repo, _ := cmd.Flags().GetString("repo")
    unresolved, _ := cmd.Flags().GetBool("unresolved")
    author, _ := cmd.Flags().GetString("author")
    format, _ := cmd.Flags().GetString("format")
    
    // Implementation
    client, err := api.NewClient()
    if err != nil {
        return fmt.Errorf("failed to create API client: %w", err)
    }
    
    // Get repository context
    owner, name, err := getRepository(repo)
    if err != nil {
        return err
    }
    
    // Get PR number
    if prNum == 0 {
        prNum, err = getCurrentPR()
        if err != nil {
            return fmt.Errorf("no PR specified and none found for current branch\nUse --pr <number> to specify a PR")
        }
    }
    
    // Fetch threads
    threads, err := client.ListThreads(owner, name, prNum)
    if err != nil {
        return fmt.Errorf("failed to fetch threads: %w", err)
    }
    
    // Filter
    if unresolved {
        threads = filterUnresolved(threads)
    }
    if author != "" {
        threads = filterByAuthor(threads, author)
    }
    
    // Output
    return renderThreads(threads, format)
}
```

### Reply Command

**File: `internal/commands/reply.go`**

```go
package commands

import (
    "fmt"
    
    "github.com/spf13/cobra"
    "github.com/hamishmorgan/gh-talk/internal/api"
)

var replyCmd = &cobra.Command{
    Use:   "reply [thread-id] [message]",
    Short: "Reply to a review thread",
    Long: `Reply to a review thread with a message.

Arguments:
  thread-id   Thread ID (PRRT_...), URL, or omit for interactive selection
  message     Reply message text, or use --editor

Examples:
  # Interactive mode (prompts for thread and message)
  gh talk reply
  
  # With full thread ID
  gh talk reply PRRT_kwDOQN97u85gQeTN "Fixed in latest commit"
  
  # With message and resolve
  gh talk reply PRRT_kwDOQN97u85gQeTN "Done!" --resolve
  
  # Using editor
  gh talk reply PRRT_kwDOQN97u85gQeTN --editor`,
    Args: cobra.RangeArgs(0, 2), // 0-2 arguments
    RunE: runReply,
}

func init() {
    replyCmd.Flags().BoolP("editor", "e", false, "Open editor for message composition")
    replyCmd.Flags().Bool("resolve", false, "Resolve thread after replying")
    replyCmd.Flags().StringP("message", "m", "", "Message text (alternative to argument)")
    
    // Flag conflicts
    replyCmd.MarkFlagsMutuallyExclusive("editor", "message")
}

func runReply(cmd *cobra.Command, args []string) error {
    // Parse arguments
    var threadID, message string
    
    switch len(args) {
    case 0:
        // Interactive mode
        id, err := selectThreadInteractively()
        if err != nil {
            return err
        }
        threadID = id
        
    case 1:
        // Thread ID provided, message via flag or editor
        threadID = args[0]
        
    case 2:
        // Both provided
        threadID = args[0]
        message = args[1]
    }
    
    // Get message if not provided
    useEditor, _ := cmd.Flags().GetBool("editor")
    msgFlag, _ := cmd.Flags().GetString("message")
    
    if message == "" {
        if msgFlag != "" {
            message = msgFlag
        } else if useEditor {
            msg, err := openEditor()
            if err != nil {
                return err
            }
            message = msg
        } else {
            return fmt.Errorf("no message provided\nUse --editor or --message flag, or provide message as argument")
        }
    }
    
    // Validate
    if message == "" {
        return fmt.Errorf("message cannot be empty")
    }
    
    // Execute
    client, err := api.NewClient()
    if err != nil {
        return err
    }
    
    err = client.ReplyToThread(threadID, message)
    if err != nil {
        return err
    }
    
    fmt.Printf("✓ Replied to thread %s\n", threadID)
    
    // Resolve if requested
    shouldResolve, _ := cmd.Flags().GetBool("resolve")
    if shouldResolve {
        err = client.ResolveThread(threadID)
        if err != nil {
            return fmt.Errorf("replied successfully but failed to resolve: %w", err)
        }
        fmt.Printf("✓ Resolved thread\n")
    }
    
    return nil
}
```

### Resolve Command

**File: `internal/commands/resolve.go`**

```go
package commands

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
    Use:   "resolve [thread-id...]",
    Short: "Resolve review threads",
    Long: `Mark one or more review threads as resolved.

Arguments:
  thread-id   One or more thread IDs, or omit for interactive selection

Examples:
  # Interactive selection
  gh talk resolve
  
  # Single thread
  gh talk resolve PRRT_kwDOQN97u85gQeTN
  
  # Multiple threads
  gh talk resolve PRRT_abc123 PRRT_def456 PRRT_ghi789
  
  # With message
  gh talk resolve PRRT_abc123 --message "Fixed in commit abc123"`,
    Args: cobra.MinimumNArgs(0), // 0 or more arguments
    RunE: runResolve,
}

func init() {
    resolveCmd.Flags().StringP("message", "m", "", "Message to post before resolving")
    resolveCmd.Flags().BoolP("yes", "y", false, "Skip confirmation for multiple threads")
}

func runResolve(cmd *cobra.Command, args []string) error {
    var threadIDs []string
    
    if len(args) == 0 {
        // Interactive selection
        ids, err := selectThreadsInteractively()
        if err != nil {
            return err
        }
        threadIDs = ids
    } else {
        threadIDs = args
    }
    
    // Confirm for multiple threads
    if len(threadIDs) > 1 {
        skipConfirm, _ := cmd.Flags().GetBool("yes")
        if !skipConfirm {
            confirmed, err := confirm(fmt.Sprintf("Resolve %d threads?", len(threadIDs)))
            if err != nil || !confirmed {
                return fmt.Errorf("cancelled")
            }
        }
    }
    
    // Post message first if provided
    message, _ := cmd.Flags().GetString("message")
    
    client, err := api.NewClient()
    if err != nil {
        return err
    }
    
    // Resolve each thread
    for _, id := range threadIDs {
        if message != "" {
            if err := client.ReplyToThread(id, message); err != nil {
                return fmt.Errorf("failed to add message to %s: %w", id, err)
            }
        }
        
        if err := client.ResolveThread(id); err != nil {
            return fmt.Errorf("failed to resolve %s: %w", id, err)
        }
        
        fmt.Printf("✓ Resolved %s\n", id)
    }
    
    return nil
}
```

## Advanced Features

### Args Validation

**Built-in Validators:**

```go
Args: cobra.NoArgs              // No arguments allowed
Args: cobra.ExactArgs(1)        // Exactly 1 argument
Args: cobra.MinimumNArgs(1)     // At least 1 argument
Args: cobra.MaximumNArgs(2)     // At most 2 arguments
Args: cobra.RangeArgs(0, 2)     // Between 0 and 2 arguments
```

**Custom Validation:**

```go
var replyCmd = &cobra.Command{
    Use:  "reply [thread-id] [message]",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) > 2 {
            return fmt.Errorf("too many arguments")
        }
        
        if len(args) >= 1 {
            // Validate thread ID format
            if !isValidThreadID(args[0]) {
                return fmt.Errorf("invalid thread ID format: %s", args[0])
            }
        }
        
        return nil
    },
    RunE: runReply,
}
```

### PreRun and PostRun Hooks

**Execution Order:**

```
PersistentPreRun  (parent)
↓
PreRun            (current command)
↓
Run/RunE          (current command)
↓
PostRun           (current command)
↓
PersistentPostRun (parent)
```

**Example:**

```go
var rootCmd = &cobra.Command{
    Use: "gh-talk",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Runs before EVERY command
        // Good for: setup, authentication checks, logging
        initializeClient()
    },
}

var replyCmd = &cobra.Command{
    Use: "reply",
    PreRun: func(cmd *cobra.Command, args []string) {
        // Runs before this command only
        // Good for: command-specific validation
        validateThreadID()
    },
    RunE: runReply,
    PostRun: func(cmd *cobra.Command, args []string) {
        // Runs after successful execution
        // Good for: cleanup, cache invalidation
        invalidateCache()
    },
}
```

**Use in gh-talk:**

```go
var rootCmd = &cobra.Command{
    Use: "gh-talk",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Validate we can create API client
        _, err := api.NewClient()
        if err != nil {
            return fmt.Errorf("authentication failed\nRun 'gh auth login' to authenticate")
        }
        return nil
    },
}
```

### Flag Binding and Validation

**Required Flags:**

```go
cmd.Flags().String("format", "", "Output format")
cmd.MarkFlagRequired("format")  // Error if not provided
```

**Mutually Exclusive:**

```go
cmd.Flags().Bool("resolved", false, "Show resolved")
cmd.Flags().Bool("unresolved", false, "Show unresolved")
cmd.Flags().Bool("all", false, "Show all")

cmd.MarkFlagsMutuallyExclusive("resolved", "unresolved", "all")
```

**Required Together:**

```go
cmd.Flags().String("hide", "", "Hide reason")
cmd.Flags().String("comment", "", "Comment ID")

cmd.MarkFlagsRequiredTogether("hide", "comment")
```

**One Required:**

```go
cmd.Flags().String("pr", "", "PR number")
cmd.Flags().String("issue", "", "Issue number")

cmd.MarkFlagsOneRequired("pr", "issue")
```

### Help and Usage

**Automatic Help:**

```bash
gh talk --help
gh talk list --help
gh talk list threads --help
```

**Custom Help:**

```go
cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
    // Custom help rendering
})

cmd.SetUsageFunc(func(cmd *cobra.Command) error {
    // Custom usage rendering
    return nil
})
```

**Example Template:**

```go
cmd.SetHelpTemplate(`{{.Long}}

Usage:
  {{.UseLine}}

{{if .HasAvailableSubCommands}}Available Commands:
{{range .Commands}}{{if .IsAvailableCommand}}  {{.Name | printf "%-15s"}} {{.Short}}
{{end}}{{end}}{{end}}
{{if .HasAvailableLocalFlags}}Flags:
{{.LocalFlags.FlagUsages}}{{end}}

Examples:
{{.Example}}
`)
```

## Common Patterns for gh-talk

### Pattern 1: Optional Interactive

**Allow argument OR interactive selection:**

```go
var reactCmd = &cobra.Command{
    Use:   "react <comment-id> <emoji>",
    Short: "Add emoji reaction to a comment",
    Args:  cobra.RangeArgs(0, 2), // Can be 0, 1, or 2 args
    RunE: func(cmd *cobra.Command, args []string) error {
        var commentID, emoji string
        
        switch len(args) {
        case 0:
            // Interactive: select comment and emoji
            id, err := selectComment()
            if err != nil {
                return err
            }
            commentID = id
            
            em, err := selectEmoji()
            if err != nil {
                return err
            }
            emoji = em
            
        case 1:
            // Comment ID provided, prompt for emoji
            commentID = args[0]
            em, err := selectEmoji()
            if err != nil {
                return err
            }
            emoji = em
            
        case 2:
            // Both provided
            commentID = args[0]
            emoji = args[1]
        }
        
        return addReaction(commentID, emoji)
    },
}
```

### Pattern 2: Context-Aware Defaults

**Use persistent flags with fallbacks:**

```go
func getRepository(cmd *cobra.Command) (owner, name string, err error) {
    // Try persistent flag
    repo, _ := cmd.Flags().GetString("repo")
    if repo != "" {
        r, err := repository.Parse(repo)
        if err != nil {
            return "", "", err
        }
        return r.Owner, r.Name, nil
    }
    
    // Fall back to git context
    r, err := repository.Current()
    if err != nil {
        return "", "", fmt.Errorf("cannot determine repository\nUse --repo OWNER/REPO flag")
    }
    
    return r.Owner, r.Name, nil
}
```

### Pattern 3: Flag Precedence

**Handle multiple ways to specify same thing:**

```go
func getMessage(cmd *cobra.Command, args []string) (string, error) {
    // Priority:
    // 1. Positional argument
    // 2. --message flag
    // 3. --editor flag
    // 4. Interactive prompt
    
    if len(args) >= 2 {
        return args[1], nil
    }
    
    msgFlag, _ := cmd.Flags().GetString("message")
    if msgFlag != "" {
        return msgFlag, nil
    }
    
    useEditor, _ := cmd.Flags().GetBool("editor")
    if useEditor {
        return openEditor()
    }
    
    return promptForMessage()
}
```

### Pattern 4: Bulk Operations

**Handle multiple arguments:**

```go
var resolveCmd = &cobra.Command{
    Use:   "resolve [thread-id...]",
    Short: "Resolve one or more threads",
    Args:  cobra.MinimumNArgs(0),
    RunE: func(cmd *cobra.Command, args []string) error {
        threadIDs := args
        
        if len(threadIDs) == 0 {
            // Interactive multi-select
            ids, err := selectThreadsMulti()
            if err != nil {
                return err
            }
            threadIDs = ids
        }
        
        // Confirm for bulk
        if len(threadIDs) > 1 {
            yes, _ := cmd.Flags().GetBool("yes")
            if !yes {
                if !confirmBulk(len(threadIDs)) {
                    return fmt.Errorf("cancelled")
                }
            }
        }
        
        // Process each
        for _, id := range threadIDs {
            if err := resolveThread(id); err != nil {
                return err
            }
        }
        
        return nil
    },
}
```

### Pattern 5: Type-Safe Flags

**Bind to variables:**

```go
var (
    flagRepo       string
    flagPR         int
    flagUnresolved bool
    flagFormat     string
)

func init() {
    listThreadsCmd.Flags().StringVarP(&flagRepo, "repo", "R", "", "Repository")
    listThreadsCmd.Flags().IntVar(&flagPR, "pr", 0, "PR number")
    listThreadsCmd.Flags().BoolVar(&flagUnresolved, "unresolved", true, "Unresolved only")
    listThreadsCmd.Flags().StringVar(&flagFormat, "format", "table", "Output format")
}

func runListThreads(cmd *cobra.Command, args []string) error {
    // Use variables directly (no Get* calls needed)
    if flagUnresolved {
        threads = filterUnresolved(threads)
    }
    
    return renderThreads(threads, flagFormat)
}
```

## Testing Commands

### Unit Testing Commands

```go
package commands_test

import (
    "bytes"
    "testing"
    
    "github.com/hamishmorgan/gh-talk/internal/commands"
)

func TestReplyCommand(t *testing.T) {
    // Capture output
    output := new(bytes.Buffer)
    
    // Create root command
    rootCmd := commands.NewRootCommand()
    rootCmd.SetOut(output)
    rootCmd.SetErr(output)
    
    // Set args
    rootCmd.SetArgs([]string{
        "reply",
        "PRRT_kwDOQN97u85gQeTN",
        "Test message",
    })
    
    // Execute
    err := rootCmd.Execute()
    
    // Assert
    if err != nil {
        t.Errorf("command failed: %v", err)
    }
    
    if !bytes.Contains(output.Bytes(), []byte("✓ Replied")) {
        t.Errorf("expected success message, got: %s", output.String())
    }
}
```

### Testing Flags

```go
func TestListThreadsFlags(t *testing.T) {
    cmd := commands.NewListThreadsCommand()
    
    // Test mutually exclusive flags
    cmd.SetArgs([]string{"--resolved", "--unresolved"})
    err := cmd.Execute()
    
    if err == nil {
        t.Error("expected error for mutually exclusive flags")
    }
}
```

## Shell Completion

### Generate Completion

**Cobra provides automatic completion:**

```go
// Add completion command to root
rootCmd.AddCommand(completionCmd)

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion script",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        switch args[0] {
        case "bash":
            return cmd.Root().GenBashCompletion(os.Stdout)
        case "zsh":
            return cmd.Root().GenZshCompletion(os.Stdout)
        case "fish":
            return cmd.Root().GenFishCompletion(os.Stdout, true)
        case "powershell":
            return cmd.Root().GenPowerShellCompletion(os.Stdout)
        }
        return fmt.Errorf("unsupported shell: %s", args[0])
    },
}
```

**Installation:**

```bash
# Bash
gh talk completion bash > /usr/local/etc/bash_completion.d/gh-talk

# Zsh  
gh talk completion zsh > "${fpath[1]}/_gh-talk"

# Fish
gh talk completion fish > ~/.config/fish/completions/gh-talk.fish
```

### Dynamic Completion

**Complete thread IDs:**

```go
var replyCmd = &cobra.Command{
    Use:   "reply [thread-id] [message]",
    ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        if len(args) == 0 {
            // Complete thread IDs
            threads, _ := fetchThreads()
            ids := make([]string, len(threads))
            for i, t := range threads {
                ids[i] = t.ID + "\t" + t.Preview // ID + description
            }
            return ids, cobra.ShellCompDirectiveNoFileComp
        }
        return nil, cobra.ShellCompDirectiveNoFileComp
    },
    RunE: runReply,
}
```

## Error Handling

### Cobra Error Patterns

**Return Errors from RunE:**

```go
RunE: func(cmd *cobra.Command, args []string) error {
    // Return errors, don't handle here
    if err := validateInput(); err != nil {
        return err  // Cobra handles printing
    }
    
    if err := doWork(); err != nil {
        return fmt.Errorf("failed to do work: %w", err)
    }
    
    return nil
}
```

**Custom Error Handling:**

```go
func Execute() error {
    err := rootCmd.Execute()
    if err != nil {
        // Custom error formatting
        fmt.Fprintf(os.Stderr, "\n✗ %v\n\n", err)
        os.Exit(1)
    }
    return nil
}
```

**Silence Default Error:**

```go
rootCmd.SilenceErrors = true  // Don't print errors twice
rootCmd.SilenceUsage = true   // Don't show usage on errors
```

## Integration with go-gh

### Combining Cobra with go-gh

**Perfect Together:**

```go
package commands

import (
    "github.com/spf13/cobra"
    "github.com/cli/go-gh/v2/pkg/api"
    "github.com/cli/go-gh/v2/pkg/repository"
    "github.com/cli/go-gh/v2/pkg/term"
    "github.com/cli/go-gh/v2/pkg/tableprinter"
)

var listThreadsCmd = &cobra.Command{
    Use:   "threads",
    Short: "List review threads",
    RunE: func(cmd *cobra.Command, args []string) error {
        // go-gh: Get repository
        repo, err := repository.Current()
        if err != nil {
            return err
        }
        
        // go-gh: Create API client
        client, err := api.DefaultGraphQLClient()
        if err != nil {
            return err
        }
        
        // Query threads
        threads, err := queryThreads(client, repo)
        if err != nil {
            return err
        }
        
        // go-gh: Terminal-adaptive output
        terminal := term.FromEnv()
        width, _, _ := terminal.Size()
        
        // go-gh: Table printer
        t := tableprinter.New(terminal.Out(), terminal.IsTerminalOutput(), width)
        
        for _, thread := range threads {
            t.AddField(thread.ID)
            t.AddField(thread.Path)
            t.EndRow()
        }
        
        return t.Render()
    },
}
```

**No Conflicts:**

- Cobra handles command structure
- go-gh handles GitHub integration
- Each does what it's best at

## Best Practices

### DO

✅ **Use RunE, not Run** - Return errors instead of handling  
✅ **Validate in Args** - Use built-in validators  
✅ **Short + Long descriptions** - Helpful help text  
✅ **Examples in Long** - Show real usage  
✅ **PersistentFlags for global** - --repo, --pr, --json  
✅ **Local flags for command-specific** - --editor, --resolve  
✅ **SilenceUsage and SilenceErrors** - Clean error output  
✅ **Use preRun for validation** - Check before executing  

### DON'T

❌ **Exit directly** - Return errors instead  
❌ **Print in commands** - Use cmd.Out() or return errors  
❌ **Skip validation** - Use Args and flag validators  
❌ **Ignore context** - Pass context to API calls  
❌ **Hardcode output** - Use term.FromEnv()  

## Quick Start for gh-talk

### 1. Add Dependency

```bash
go get github.com/spf13/cobra@latest
```

### 2. Create Root Command

```go
// internal/commands/root.go
package commands

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:           "gh-talk",
    Short:         "Manage GitHub PR conversations",
    SilenceUsage:  true,
    SilenceErrors: true,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "✗ %v\n", err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().StringP("repo", "R", "", "Repository (OWNER/REPO)")
    rootCmd.AddCommand(listCmd)
    // ... add other commands
}
```

### 3. Update main.go

```go
package main

import "github.com/hamishmorgan/gh-talk/internal/commands"

func main() {
    commands.Execute()
}
```

### 4. Add First Command

```go
// internal/commands/list.go
package commands

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List conversations",
}

var listThreadsCmd = &cobra.Command{
    Use:   "threads",
    Short: "List review threads",
    RunE:  runListThreads,
}

func init() {
    listCmd.AddCommand(listThreadsCmd)
    listThreadsCmd.Flags().Bool("unresolved", true, "Show unresolved only")
}

func runListThreads(cmd *cobra.Command, args []string) error {
    // Implementation
    fmt.Println("Listing threads...")
    return nil
}
```

### 5. Build and Test

```bash
go build
./gh-talk list threads --help
./gh-talk list threads --unresolved
```

## Resources

- **Cobra Repo:** <https://github.com/spf13/cobra>
- **Cobra Website:** <https://cobra.dev>
- **User Guide:** <https://github.com/spf13/cobra/blob/main/site/content/user_guide.md>
- **Examples:** <https://github.com/spf13/cobra/tree/main/site/content>

### Example Projects Using Cobra

- **kubectl:** <https://github.com/kubernetes/kubectl>
- **docker:** <https://github.com/docker/cli>
- **hugo:** <https://github.com/gohugoio/hugo>
- **gh actions:** GitHub's own CLI tool!

---

**Last Updated**: 2025-11-02  
**Cobra Version**: v1.8+  
**Context**: Implementation guide for gh-talk using Cobra  
**Decision**: Finalized - using Cobra for CLI framework
