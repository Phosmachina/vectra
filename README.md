<h1 align="center">
    <img src="docs/vectra_banner.svg" alt="Banner" width="50%">
</h1>

<div align="center">

[![GoDoc](https://godoc.org/github.com/Phosmachina/vectra?status.svg)](https://pkg.go.dev/github.com/Phosmachina/vectra#section-documentation)
[![Go Report Card](https://goreportcard.com/badge/github.com/Phosmachina/vectra)](https://goreportcard.com/report/github.com/Phosmachina/vectra)

</div>

<details>
<summary>
 Table of contents
</summary>

<!-- TOC -->

* [🎯 Overview](#-overview)
* [⚡️ Features](#-features)
* [🚀 Getting started](#-getting-started)
    * [Prerequisite](#prerequisite)
    * [Install Vectra](#install-vectra)
    * [Deploy](#deploy)
    * [Run](#run)
* [🤝 Contributing](#-contributing)
* [🕘 What's next](#-whats-next)

<!-- TOC -->

</details>

## 🎯 Overview

The main goal of the project is to create a versatile multi-language template and toolkit
for website servers and backend systems. It strives to achieve this by integrating the
best existing technologies, resulting in an efficient and fast server experience. The
goal is to minimize the complexity of development as much as possible. Ultimately,
Vectra's goal is to leverage the unique design of each technology to achieve specific
goals without investing excessive time and effort.

By combining these technologies (Go, Pug, Sass, Badger, ...), Vectra offers a robust and 
streamlined development
environment. It reduces the need for complex setups and integrations, allowing you to
focus on building the core functionality and design of your website.

## ⚡️ Features

- **Code generation**
    - Controllers (view and service routes)
    - Types (storage, ajax, view)
    - Service (defines interface)
- **MVC architecture**
- **Pipeline for [Sass](https://sass-lang.com/) and [Pug](https://github.com/Joker/jade)**
    - All in one docker with needed tools
    - Jetbrains file watchers configuration

- **Web framework integrated: [Fiber](https://.gofiber.io)**
    - Separation for static and main app
    - Middlewares configured (log, compression, cache, csrf, ...)
- **Data validation with [Validator](https://github.com/go-playground/validator)
  and [Mold](https://github.com/go-playground/mold)**
- **KV helper, [FluentKV](https://github.com/phosmachina/FluentKV), for BadgerDB**

- **Connection system**
    - First connection mechanism
    - User and roles

- **JS helpers**
    - Ajax
    - Form data scrap
    - Svg sprite loader
    - Components

- **Integrated i18n system**

[//]: # (TODO add image magick command to compress image to AVIF)

[//]: # (TODO make i18n as an independant library?)

## 🚀 Getting started

### Prerequisite

- Docker
- Go SDK (or build and run in Docker)

### Install Vectra

If you are Go SDK, install with `go` command:

```shell
go install github.com/Phosmachina/vectra@latest
```

### Deploy

- This command permits writing a default config:
  ```shell
  vectra -p path/YourProject init
  ```

- Edit the configuration, `YourProject/.vectra/project.yml`, as your convenience.

- Run `vectra` for a full generation:
  ```shell
  vectra -p path/YourProject gen
  ```

- Launch watcher: the first time it might take some time because of container
  creation and image download:
  ```shell
  vectra -p path/YourProject watch
  ```

If you want to re-edit the configuration, maybe after that run a partial generation like 
this to avoid file overwriting:
```shell
vectra -p path/YourProject -s types,controlers,services gen
```

### Run

Now you can open the folder `path/YourProject`, which Vectra created as a project with
your IDE.

You need to make sure that the `*.pug` files are correctly transpiled to Go (there are
transpiled to `src/view/go/`).
Currently, with the file watchers, you need to make a change to the
files to trigger it.

After that, you can start your application.
This can be done manually by executing the following command (in the root directory of the
project):

```shell
go run app.go
```

## 🤝 Contributing

Your contributions are always valued and appreciated!

Thank you in advance for making this project even better. I'm excited to see your
contributions!

## 🕘 What's next

Improving and expanding this project is my perpetual goal.
Here's an insight into what I plan next:

- **Component architecture**: I want to provide a simple way to develop components by
  simplify the boilerplate between view and controller.
- **Default components**: In the future, I plan to incorporate a set of
  default components into the project. This will help in establishing a
  consistent UI/UX throughout and will also save time and effort in design and
  development.
- **RBAC, ACL robust system**: Replace the current system with a robust and proven
  system like Casbin. This integration should help to deal with complex access management.

I value your ideas, contributions, and feedback. Stay tuned for the next steps on this
exciting journey!
