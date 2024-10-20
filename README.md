# gitbm
A fast *(probably)*, lightweight *(I guess so)* Git branch bookmarking tool written in Go, powered by SQLite.

![image](https://github.com/user-attachments/assets/294dcc4b-b6bc-42a7-904f-a041a0c17d31)

## Why `gitbm`?
When you’re working on a big feature or in a team environment, managing multiple branches can become quite a challenge. Here’s why gitbm helps (or at least tries to help):
- You’re working on a huge feature split into multiple smaller tasks. Each task has its own branch, and you can’t remember all those branch names.
- If you have the memory of an elephant, this tool is not for you. But if you’re like me, who forgets branch names after a day or two, then **this tool is definitely for you**.
- Your team uses weird branch naming conventions like `JIRA-xyz123` or `CLICKUP-456`, and remembering these names is like remembering the 100th digit of Pi. 
- You just want to bookmark branches with meaningful alias like `boss-needs-it-tomorrow` and not worry about remembering the cryptic original names.

## How It Works
- Everything happens **locally** in your Git repository. Your git configs are safe and untouched.
- A **SQLite database** (fancy but necessary!) is created in the `.git` directory of your repo to store branch names and aliases.
- Every time you run a `gitbm` command, the tool reads from this magical database, does its thing, and updates the database as needed.

## Installation
1. Make sure you’ve got Go installed (if not, [download it here](https://golang.org/dl/)).
2. Run the following command to install `gitbm`:
    ```bash
    go get -u github.com/devadathanmb/gitbm
    ```

## Usage
https://github.com/user-attachments/assets/78a32e1d-2cc3-48c3-abf9-2cfa65474054

For full command details, just ask for help:
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
- [ ] Add `recent` command with automatic branch tracking.
- [ ] Shell completion (because typing is hard).
- [ ] Fuzzy search (FZF) for `remove` and `delete` commands.
- [ ] Better CLI output.
- [ ] Better error messages.
- [ ] Add some tests (maybe?)

## License
This project is licensed under the GPL-3.0. See [LICENSE](LICENSE.md) for the details (because legal stuff is important, apparently).
