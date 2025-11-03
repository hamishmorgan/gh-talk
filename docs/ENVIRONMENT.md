# Environment Variables

**Reference guide for all environment variables used by gh-talk**

## GitHub CLI Variables

These are automatically used by go-gh and require no additional code:

### `GH_TOKEN`

**Purpose:** GitHub authentication token  
**Used By:** go-gh API clients  
**Default:** Token from `gh auth login`

**Example:**

```bash
export GH_TOKEN=ghp_abc123xyz456
gh talk list threads
```

### `GH_HOST`

**Purpose:** GitHub host for API requests  
**Used By:** go-gh API clients  
**Default:** `github.com`

**Example (GitHub Enterprise):**

```bash
export GH_HOST=github.example.com
gh talk list threads
```

### `GH_REPO`

**Purpose:** Override current repository context  
**Used By:** go-gh repository detection  
**Default:** Detected from git remotes

**Example:**

```bash
export GH_REPO=owner/repo
gh talk list threads  # Uses owner/repo
```

### `GH_FORCE_TTY`

**Purpose:** Force terminal mode even when not in TTY  
**Used By:** go-gh term package  
**Default:** Auto-detected

**Example:**

```bash
export GH_FORCE_TTY=1
gh talk list threads | less  # Still shows table format
```

**Value Options:**

- `1` or `true` - Force TTY mode
- `<number>` - Force specific width (e.g., `80`)
- `<percentage>%` - Force percentage of actual width (e.g., `50%`)

### `GH_DEBUG`

**Purpose:** Enable debug logging for API requests  
**Used By:** go-gh API clients  
**Default:** Disabled

**Example:**

```bash
export GH_DEBUG=1
gh talk list threads
# Logs all GraphQL queries and responses
```

## Terminal Variables

These control color and terminal behavior:

### `NO_COLOR`

**Purpose:** Disable all color output  
**Standard:** <https://no-color.org>  
**Default:** Not set (colors enabled in TTY)

**Example:**

```bash
NO_COLOR=1 gh talk list threads
# Output without ANSI color codes
```

### `CLICOLOR`

**Purpose:** Control color support  
**Values:** `0` (disable) or `1` (enable)  
**Default:** Auto-detected from terminal

**Example:**

```bash
export CLICOLOR=0  # Disable colors
```

### `CLICOLOR_FORCE`

**Purpose:** Force color output even in non-TTY  
**Default:** Not set

**Example:**

```bash
export CLICOLOR_FORCE=1
gh talk list threads | cat  # Colors even though piped
```

### `TERM`

**Purpose:** Terminal type identifier  
**Used By:** Color capability detection  
**Default:** Set by terminal

**Common Values:**

- `xterm-256color` - 256 color support
- `xterm` - Basic color support
- `dumb` - No special features

### `COLORTERM`

**Purpose:** True color support indicator  
**Values:** `truecolor` or `24bit`  
**Default:** Set by modern terminals

## gh-talk Specific Variables

### `GH_TALK_CONFIG`

**Purpose:** Configuration file location  
**Default:** `~/.config/gh-talk/config.yml`

**Example:**

```bash
export GH_TALK_CONFIG=/path/to/custom-config.yml
gh talk list threads
```

### `GH_TALK_CACHE_DIR`

**Purpose:** Cache directory for API responses  
**Default:** `~/.cache/gh-talk`

**Example:**

```bash
export GH_TALK_CACHE_DIR=/tmp/gh-talk-cache
gh talk list threads
```

**Cache Contents:**

- Thread data (5 minute TTL)
- Repository metadata
- User information

### `GH_TALK_CACHE_TTL`

**Purpose:** Cache time-to-live in minutes  
**Default:** `5`

**Example:**

```bash
export GH_TALK_CACHE_TTL=10  # Cache for 10 minutes
gh talk list threads
```

**Values:**

- `0` - Disable caching
- `<number>` - Minutes to cache

### `GH_TALK_FORMAT`

**Purpose:** Default output format  
**Default:** `table`  
**Values:** `table`, `json`, `tsv`

**Example:**

```bash
export GH_TALK_FORMAT=json
gh talk list threads  # Always outputs JSON
```

**Override:**

```bash
# Environment says json, but flag overrides
export GH_TALK_FORMAT=json
gh talk list threads --format table  # Shows table
```

### `GH_TALK_EDITOR`

**Purpose:** Editor for message composition  
**Default:** Value of `$EDITOR`, falls back to `vim`

**Example:**

```bash
export GH_TALK_EDITOR=nano
gh talk reply --editor  # Opens nano
```

**Fallback Chain:**

1. `GH_TALK_EDITOR`
2. `EDITOR`
3. `vim` (hardcoded fallback)

## Standard Environment Variables

These are standard shell variables gh-talk respects:

### `EDITOR`

**Purpose:** Default text editor  
**Used By:** Editor integration  
**Default:** `vim`

**Example:**

```bash
export EDITOR=code
gh talk reply --editor  # Opens VS Code
```

### `HOME`

**Purpose:** User home directory  
**Used By:** Config and cache path resolution  
**Default:** Set by system

**Used For:**

- `~/.config/gh-talk/` → `$HOME/.config/gh-talk/`
- `~/.cache/gh-talk/` → `$HOME/.cache/gh-talk/`

### `XDG_CONFIG_HOME`

**Purpose:** XDG config directory (Linux standard)  
**Default:** `~/.config`

**Used For:**

```bash
export XDG_CONFIG_HOME=/custom/config
# Config: /custom/config/gh-talk/config.yml
```

### `XDG_CACHE_HOME`

**Purpose:** XDG cache directory (Linux standard)  
**Default:** `~/.cache`

**Used For:**

```bash
export XDG_CACHE_HOME=/custom/cache
# Cache: /custom/cache/gh-talk/
```

## Usage Examples

### Disable Caching

```bash
GH_TALK_CACHE_TTL=0 gh talk list threads
```

### Debug Mode

```bash
GH_DEBUG=1 gh talk list threads 2> debug.log
```

### Force JSON Output

```bash
GH_TALK_FORMAT=json gh talk list threads | jq '.[] | select(.isResolved == false)'
```

### Use Different GitHub Instance

```bash
export GH_HOST=github.example.com
export GH_TOKEN=ghp_enterpriseToken123
gh talk list threads
```

### Custom Config and Cache

```bash
export GH_TALK_CONFIG=/tmp/test-config.yml
export GH_TALK_CACHE_DIR=/tmp/test-cache
gh talk list threads
```

### No Colors in CI

```bash
# CI environments should set
export NO_COLOR=1
export GH_FORCE_TTY=0
```

## Environment Variable Priority

**For all settings, priority is:**

1. **Command-line flags** (highest priority)

   ```bash
   gh talk list threads --format json
   ```

2. **gh-talk environment variables**

   ```bash
   export GH_TALK_FORMAT=json
   ```

3. **Configuration file**

   ```yaml
   # ~/.config/gh-talk/config.yml
   defaults:
     format: json
   ```

4. **GitHub CLI environment variables**

   ```bash
   export GH_REPO=owner/repo
   ```

5. **Defaults** (lowest priority)
   - Built-in sensible defaults

**Example:**

```bash
# Config file says: format=table
# Environment says: GH_TALK_FORMAT=json
# Flag says: --format tsv
# Result: TSV (flag wins)
```

## Testing with Environment Variables

### In Tests

```go
func TestWithEnv(t *testing.T) {
    // Save original
    original := os.Getenv("GH_TALK_FORMAT")
    defer os.Setenv("GH_TALK_FORMAT", original)
    
    // Set test value
    os.Setenv("GH_TALK_FORMAT", "json")
    
    // Test
    format := getDefaultFormat()
    if format != "json" {
        t.Errorf("expected json, got %s", format)
    }
}
```

### In CI

```yaml
- name: Test with environment
  env:
    GH_TALK_FORMAT: json
    GH_TALK_CACHE_TTL: 0
  run: go test ./...
```

## Reference

### All Variables Quick Reference

| Variable | Purpose | Default | Example |
|----------|---------|---------|---------|
| `GH_TOKEN` | Auth token | From gh auth | `ghp_abc123` |
| `GH_HOST` | GitHub host | `github.com` | `github.example.com` |
| `GH_REPO` | Repository | From git | `owner/repo` |
| `GH_FORCE_TTY` | Force TTY | Auto-detect | `1` or `80` |
| `GH_DEBUG` | Debug mode | Off | `1` |
| `NO_COLOR` | Disable colors | Not set | `1` |
| `CLICOLOR` | Color support | Auto | `0` or `1` |
| `CLICOLOR_FORCE` | Force colors | Not set | `1` |
| `GH_TALK_CONFIG` | Config file | `~/.config/gh-talk/config.yml` | `/path/to/config.yml` |
| `GH_TALK_CACHE_DIR` | Cache dir | `~/.cache/gh-talk` | `/tmp/cache` |
| `GH_TALK_CACHE_TTL` | Cache TTL (min) | `5` | `10` or `0` |
| `GH_TALK_FORMAT` | Output format | `table` | `json` or `tsv` |
| `GH_TALK_EDITOR` | Text editor | `$EDITOR` or `vim` | `nano` or `code` |
| `EDITOR` | Default editor | `vim` | `nano` |

### Related Documentation

- [GO-GH.md](GO-GH.md) - How go-gh uses GH_* variables
- [DESIGN.md](DESIGN.md) - Output format decisions
- [ENGINEERING.md](ENGINEERING.md) - Testing with environment variables

---

**Last Updated**: 2025-11-02  
**Compatibility**: gh CLI v2.0+  
**Standards**: Follows XDG Base Directory Specification
