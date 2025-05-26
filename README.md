# strikerGo

**StrikerGo** is a high-performance Blackjack simulation engine written in Go. It evaluates various card counting strategies over large numbers of hands and deck configurations.

---

## 📦 Build

```bash
make                # Builds the binary in ./bin/striker-go
```

---

## ▶️ Run

```bash
make run STRATEGY=basic DECKS=single-deck HANDS=10000000 THREADS=8
```

This runs a simulation with the specified strategy and deck configuration. Logs are saved to:

```
~/Striker/Simulations/YYYY/MM/DD/strikerGo-HHMMSS.log
```

---

## 🔁 Run All Strategy/Deck Combinations

```bash
make run-all
```

---

## 🧹 Clean

```bash
make clean
```

Removes build artifacts.

---

## 🧼 Lint

```bash
make lint
```

Runs `golangci-lint`. Requires installation:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## 🛠️ Install

```bash
make install
```

Installs the binary to:

```
~/Striker/bin
```

---

## 🎮 Strategies

Pass these as `--<strategy>` arguments:

- `--mimic`
- `--linear`
- `--polynomial`
- `--neural`
- `--basic`
- `--high-low`
- `--wong`

---

## 🃏 Deck Configurations

Pass these as `--<deck>` arguments:

- `--single-deck`
- `--double-deck`
- `--six-shoe`

---

## ⚙️ Example Aliases

Define helpful shortcuts in your shell:

```bash
alias sg='./bin/striker-go'

# Run common configurations
alias sg1='sg --mimic --single-deck --number-of-hands 5000000 --number-of-threads 8'
alias sg2='sg --linear --double-deck --number-of-hands 10000000 --number-of-threads 16'
alias sg6='sg --neural --six-shoe --number-of-hands 20000000 --number-of-threads 24'

# Strategy shorthand
alias sgb='sg --basic'
alias sgh='sg --high-low'
alias sgw='sg --wong'
```

---

## 📁 Project Layout

```
.
├── Makefile
├── Makefile.run
├── main.go
├── internal/
│   ├── arguments/
│   ├── cards/
│   ├── constants/
│   ├── table/
│   └── simulator/
├── bin/
└── README.md
```

---

## ✅ Requirements

- Go 1.20+
- Optional: `golangci-lint` for linting
- Unix-like system for `make` support

---

## ✍️ Contributing

Contributions are welcome! Feel free to fork, patch, and open pull requests.

---

## 📜 License

MIT License – see `LICENSE` file for details.

