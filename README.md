# gitbm
A fast (probably), lightweight (I guess so) Git branch bookmarking tool written in Go, powered by SQLite.

## Why `gitbm`?
- You’re working on a huge feature split into multiple smaller tasks. Each task has its own branch, and you can’t remember all those branch names (Chained branches are the worst, I know)
- If you have the memory of an elephant, this tool is not for you. But if you’re like me, who forgets branch names after a day or two, then **this tool is definitely for you**.
- Your team uses weird branch naming conventions like `JIRA-xyz123` or `CLICKUP-456`, and remembering these names is like remembering your long-lost Wi-Fi password.
- You just want to bookmark branches with a cool alias like `super-cool-feature` and not worry about remembering the cryptic original names.

## How It Works
- Everything happens **locally** in your Git repository. Your git configs are safe and untouched.
- A **SQLite database** (fancy, right?) is created in the `.git` directory of your repo to store branch names and aliases.
- Every time you run a `gitbm` command, the tool reads from this magical database, does its thing, and updates the database as needed.

## Installation
1. Make sure you’ve got Go installed (if not, [download it here](https://golang.org/dl/)).
2. Run the following command to install `gitbm`:
    ```bash
    go get -u github.com/devadathanmb/gitbm
    ```

## Usage
For full command details, just ask for help (don’t be shy):
```bash
gitbm --help
```

### Example Scenarios
- Bookmark your current branch with an alias:
    ```bash
    gitbm add "my-awesome-feature"
    ```

- Switch to a bookmarked branch in style with FZF:
    ```bash
    gitbm switch
    ```

- Forget what bookmarks you’ve made? List them:
    ```bash
    gitbm list bookmarks
    ```

## TODO
- [ ] Add `recent` command with automatic branch tracking (because let’s face it, we all need this).
- [ ] Fuzzy search (FZF) for `remove` and `delete` commands (because typing is hard).
- [ ] Add some tests, maybe... eventually.

## License
This project is licensed under the GPL-3.0. See [LICENSE](LICENSE) for the details (because legal stuff is important, apparently).