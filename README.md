<style>

html{
  text-align: justify;
}

.language{
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 1 0 300px;
}

</style>

[//]: # (<img src="docs/vectra_banner.svg" alt="Banner" style="width: 100%;padding: 3rem">)

<img src="docs/test.svg" width="100%" alt="test">

[//]: # (<img src="docs/overview.svg" width="100%" alt="Click to see the source">)


<div style="display: flex; justify-content: center">
<details>
<summary>
 Table of contents
</summary>

<!-- TOC -->

- [🎯 Overview](#-overview)
- [⚡️ Features](#️-features)
- [🚀 Getting started](#-getting-started)
- [🤝 Contributing](#-contributing)
- [🕘 What's next](#-whats-next)

<!-- TOC -->

</details>
</div>

## 🎯 Overview

The main goal of the project is to create a versatile multi-language template and toolkit
for website servers and backend systems. It strives to achieve this by integrating the
best
existing technologies, resulting in an efficient and fast server experience. The goal is
to minimize the complexity of development as much as possible. Ultimately, Vectra's goal
is to leverage the unique design of each technology to achieve specific goals without
investing excessive time and effort.

<div style="display: flex; gap: 20px; flex-wrap: wrap; 
padding-top: 20px" >

<div class="language">
<img src="docs/go_logo.svg" alt="Banner" style="height: 100px">
<p>
<strong>Go</strong> is a highly efficient and scalable programming language that enables 
rapid development of web applications. It provides a rich set of libraries and tools,
making it a popular choice for building server-side applications. With Go, you can create
robust and high-performance websites that handle large amounts of traffic without compromising speed or stability.
</p>
</div>

<div class="language">
<img src="docs/sass.svg" alt="Banner" style="height: 100px">
<p>
<strong>Sass</strong> (Syntactically Awesome Style Sheets) is a mature and widely adopted CSS 
extension language. It introduces powerful features like variables, mixins, and nested selectors, enabling developers to write clean, modular, and maintainable stylesheets. Sass seamlessly integrates with Vectra, allowing you to write reusable and easily customizable styles.
</p>

</div>

<div class="language">
<img src="docs/pug_logo.svg" alt="Banner" style="height: 100px">
<p>
<strong>Pug</strong> (formerly known as Jade) is a concise and expressive templating 
language that simplifies the creation of HTML markup. It provides a clean syntax with minimal clutter, reducing the amount of code you need to write. Pug supports reusable components, layout inheritance, and conditional rendering, allowing you to create dynamic and visually appealing web pages effortlessly.
</p>
</div>

<div class="language">
<img src="docs/badgerdb_logo.png" alt="Banner" style="height: 100px">
<p>
<strong>BadgerDB</strong> is a fast and efficient key-value store written in Go. It 
provides a simple and reliable database solution for storing and retrieving data 
within your backend. With Badger, you can easily manage your backend's data persistence, 
ensuring speedy access and efficient handling of user interactions.
</p>

</div>

</div>

By combining these technologies, Vectra offers a robust and streamlined development
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

## 🤝 Contributing

Your contributions are always valued and appreciated!

Thank you in advance for making this project even better. I'm excited to see your
contributions!

## 🕘 What's next

Improving and expanding this project is my perpetual goal.
Here's an insight into what I plan next:

- **HTML / Sass default components**: In the future, I plan to incorporate a set of
  default components into the project. This will help in establishing a
  consistent UI/UX throughout and will also save time and effort in design and
  development.
- **RBAC, ACL robust system**: Replace the current system with a robust and proven
  system like Casbin. This integration should help to deal with complex access management.

I value your ideas, contributions, and feedback. Stay tuned for the next steps on this
exciting journey!
