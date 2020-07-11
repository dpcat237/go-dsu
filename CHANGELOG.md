# v0.9.0 (2020.07.11)

#### Features
- Change to display all modules in one column
- Check vulnerabilities require auth data to prevent server error
- Check vulnerabilities in analyze and preview commands change to optional

# v0.8.0 (2020.07.03)

#### Features
- Implement optional prompt for update confirmation for updates with changes

# v0.7.0 (2020.06.22)

#### Features
- Add progress bar for preview command
- Add critical severity
- Implement vulnerabilities check from OSS Index for preview command
- Implement modules download with Git
- Implement analyze command

# v0.6.0 (2020.06.15)

#### Features
- Allow passing project path for preview command

#### Fixes
- Fix empty directory path for new modules

# v0.5.0 (2020.06.14)

#### Features
- Implement in preview command check of changes in license of direct and indirect dependencies

# v0.4.0 (2020.05.31)

#### Features
- Allow optionally run local tests before and after an update of each module with rollback if tests fail

# v0.3.0 (2020.05.27)

#### Features
- Implement update of only direct modules
- Allow select interactively direct modules to update
- Add updated modules to vendor folder if it exists

# v0.2.0 (2020.05.27)

#### Features
- Check internet connection before starting the update process
- Implement clean and preview as separate commands

# v0.1.0 (2020.05.26)

#### Features
- Implement command to add missing, remove unused and update modules