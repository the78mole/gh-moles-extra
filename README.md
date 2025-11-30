# gh-moles

A GitHub CLI extension that provides tools for managing GitHub repositories.

## Installation

Install this extension using the GitHub CLI:

```bash
gh extension install the78mole/gh-moles
```

## Usage

Once installed, use the extension with the `gh moles` command:

```bash
gh moles <command>
```

## Commands

### `gh moles run cleanup`

Delete old GitHub Actions workflow runs, keeping only the most recent ones OR delete all failed runs.

#### Flags

- `-y, --yes`: Skip confirmation prompt (auto-confirm deletion)
- `-f, --failed`: Delete all failed runs instead of keeping recent ones
- `KEEP_COUNT`: Number of runs to keep (default: 20, ignored with `--failed`)

#### Examples

```bash
# Keep 20 most recent runs (with confirmation)
gh moles run cleanup

# Keep 20 most recent runs (no confirmation)
gh moles run cleanup -y

# Keep 50 most recent runs (with confirmation)
gh moles run cleanup 50

# Keep 50 most recent runs (no confirmation)
gh moles run cleanup -y 50

# Delete all failed runs (with confirmation)
gh moles run cleanup --failed

# Delete all failed runs (no confirmation)
gh moles run cleanup -y -f
```

## Development

### Building from source

```bash
go build -o gh-moles .
```

### Testing locally

```bash
./gh-moles run cleanup --help
```

## License

MIT License - see [LICENSE](LICENSE) for details.
