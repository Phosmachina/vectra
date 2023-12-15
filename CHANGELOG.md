# Changelog

## VNext

## 1.1.0

### Features

- Add templating for configuration and separate fields used for production, and those 
  only used for dev.
- Add svg sprite generator: 
  -️ Recursively takes svg files,
  -️ Auto name symbol from a path,
  -️ Make mixin for simple usage,
  -️ Minify the final result.
- Add watcher configuration permitting to disable watcher individually.
- Improve Pug Watcher: now you can define in config your layout files and all normal 
  pug files are transpiled when layout is edited.

### Refactor

- Move reports of generators in a subdirectory `.vectra/report/`.
  This permits a better separation between configuration and reports.
- Inline some yaml tag for composed types.
  Edit your config file and inline all `base` tags.  

### Fixes

- On template view types, add a parameters, `IsPageCtx`, to add the possibility to define
  multiple `GlobalCtx` (It's currently impossible to define the return type of
  constructor: it's inferred by name).
- For exchange type generation: now the json tag is build with a Camel to Snake
  transformation, not just `ToLower`.

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
