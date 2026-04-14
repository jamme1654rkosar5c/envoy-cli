# envoy-cli

> A CLI tool for managing and validating `.env` files across multiple project environments with secret diffing support.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/envoy-cli/releases) page.

---

## Usage

```bash
# Validate a .env file against a template
envoy validate --env .env --template .env.example

# Diff secrets between two environments
envoy diff --from .env.staging --to .env.production

# Check for missing or extra keys across all environments
envoy audit --dir ./environments
```

**Example output:**

```
✔ .env.staging   — 12/12 keys matched
✘ .env.production — missing: DATABASE_URL, REDIS_HOST
```

Run `envoy --help` for a full list of commands and flags.

---

## Commands

| Command    | Description                                      |
|------------|--------------------------------------------------|
| `validate` | Check a `.env` file against a template           |
| `diff`     | Show key/value differences between env files     |
| `audit`    | Scan a directory of env files for inconsistencies|

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername