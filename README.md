# gitbm
A fast *(probably)*, lightweight *(I guess so)* Git branch bookmarking tool written in Go, powered by SQLite and FZF.

![image](https://github.com/user-attachments/assets/294dcc4b-b6bc-42a7-904f-a041a0c17d31)
<img width="1440" alt="Screenshot 2024-10-20 at 9 18 46 PM" src="https://github.com/user-attachments/assets/38dcbe2c-fb8b-43c6-8674-1be3a67495fe">

<img width="1440" alt="Screenshot 2024-10-20 at 9 15 05 PM" src="https://github.com/user-attachments/assets/ce0b8c7e-ba14-4378-90bd-8ada2c8f2ba6">

<img width="1440" alt="Screenshot 2024-10-20 at 9 18 57 PM" src="https://github.com/user-attachments/assets/a9fa14e5-be27-4fa8-aa2f-165d474efbbc">



## Why `gitbm`?
When you’re working on a big feature or in a team environment, managing multiple branches can become quite a challenge. Here’s why gitbm helps (or at least tries to help):
- You’re working on a huge feature split into multiple smaller tasks. Each task has its own branch, and you can’t remember all those branch names.
- If you have the memory of an elephant, this tool is not for you. But if you’re like me, who forgets branch names after a day or two, then **this tool is definitely for you**.
- Your team uses weird branch naming conventions like `JIRA-xyz123` or `CLICKUP-456`, and remembering these names is like remembering the 100th digit of Pi. 
- You just want to bookmark branches with meaningful alias like `boss-needs-it-tomorrow` and not worry about remembering the cryptic original names.

## How It Works?
- Everything happens **locally** in your Git repository. Your git configs are safe and untouched.
- A **SQLite database** (fancy but necessary!) is created in the `.git` directory of your repo to store branch names and aliases.
- Every time you run a `gitbm` command, the tool reads from this magical database, does its thing, and updates the database as needed.

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
- [x] Shell completion (because typing is hard).
- [x] Fuzzy search (FZF) for `remove` and `delete` commands.
- [ ] Add `recent` command with automatic branch tracking with git hooks.
- [ ] Better CLI output.
- [ ] Better error messages.
- [ ] Add some tests (maybe?)

## License
This project is licensed under the GPL-3.0. See [LICENSE](LICENSE.md) for the details (because legal stuff is important, apparently).
