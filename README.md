# gitbm
A fast *(probably)*, lightweight *(I guess so)* Git branch bookmarking tool written in Go, powered by SQLite, and FZF with the magic of Git hooks.

![image](https://github.com/user-attachments/assets/294dcc4b-b6bc-42a7-904f-a041a0c17d31)

## Why `gitbm`?
When you’re working on a big feature or in a team environment, managing multiple branches can become quite a challenge. Here’s why gitbm helps (or at least tries to help):
- You’re working on a huge feature split into multiple smaller tasks. Each task has its own branch, and you can’t remember all those branch names.
- If you have the memory of an elephant, this tool is not for you. But if you’re like me, who forgets branch names after a day or two, then **this tool is definitely for you**.
- Your team uses weird branch naming conventions like `JIRA-xyz123` or `CLICKUP-456`, and remembering these names is like remembering the 100th digit of Pi. 
- You just want to bookmark branches with meaningful alias like `boss-needs-it-tomorrow` and not worry about remembering the cryptic original names.
- It knows your pain branches and the ones you’ve cried over, so you can easily checkout your favorites and weep with joy!


## How It Works?
- Everything happens **locally** in your Git repository. Your git configs are safe and untouched.
- A **SQLite database** (fancy but necessary!) is created in the `.git` directory of your repo to store branch names and aliases.
- Every time you run a `gitbm` command, the tool reads from this magical database, does its thing, and updates the database as needed.
- With those Git hooks doing their thing, gitbm remembers your recent and frequent branches so you don’t have to!

## Installation
1. Make sure you’ve got Go installed (if not, [download it here](https://golang.org/dl/)).
2. Run the following command to install `gitbm`:

    ```bash
    go install github.com/devadathanmb/gitbm@latest
    ```
3. For basic shell completions, run the following command:

    ```bash
    gitbm completion zsh >> ~/.zshrc # Example for Zsh, you can put completion files in more appropriate locations.
    ```

## Usage

For full command details, just ask for help:
```bash
gitbm --help
```

### Some Cool Stuff You Can Do:
- Fuzzy checkout to one of your latest top 10 branches:
    ```bash
    gitbm recent
    ```

- Fuzzy checkout to one of your top 10 most frequently checked out branches:
    ```bash
    gitbm frequent
    ```

- Create bookmark groups and add branches to them:
    ```bash
    gitbm create "group-name"
    gitbm add
    ```

- Fuzzy checkout to a bookmarked branch:
    ```bash
    gitbm checkout
    ```

- Fuzzy remove a bookmarked branch:
    ```bash
    gitbm remove
    ```

And many more! Check out the help command for more details.

## TODO
- [x] Shell completion (because typing is hard).
- [x] Fuzzy search (FZF) for `remove` and `delete` commands.
- [x] Add `recent` command with automatic branch tracking with git hooks.
- [ ] Track new branches automatically.
- [ ] Add reset command to `recent` and `frequent` commands. 
- [ ] Better CLI output.
- [ ] Better error messages.
- [ ] Add some tests (maybe?)

## License
This project is licensed under the GPL-3.0. See [LICENSE](LICENSE.md) for the details (because legal stuff is important, apparently).
