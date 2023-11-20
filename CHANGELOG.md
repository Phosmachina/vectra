# Changelog

## VNext

## 1.0.1

### Features

- Add a generator selector flag for CLI. This allows selecting multiple but not all 
  generators for reporting and generation.

### Style

- Minor changes in pug mixin.
- Improve logging and CLI outputs in general.

### Fixes

- Inform user when the configuration file has an invalid format.
- In cli output, print the full path instead of flag value when the path exists.
- Fix panic when the extraction of a body function fails.

## 1.0.0

### Features

- Make a cli, vectra, to perform generation and read reports.
- Make a generic generator system.
- Make some generator:
    - Base
    - Controller
    - Service
    - Types (storage, view, exchange)
    - I18n (types completion)
- ~~Add the default configuration for file's watcher on Jetbrains based IDE.~~
- Add a file watch capability to avoid dependency with IDE.

### Refactor

- Use ini files in place of csv to store i18n data.
- Migrate all template files in subfolder and embed it in executable.

### Docs

- (README): Rewrite it to be more simple and focus on essential.
