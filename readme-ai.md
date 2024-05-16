<p align="center">
  <img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" width="100" alt="project-logo">
</p>
<p align="center">
    <h1 align="center">SWIFT-BANK</h1>
</p>
<p align="center">
    <em><code>► INSERT-TEXT-HERE</code></em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/zde37/Swift-Bank?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/zde37/Swift-Bank?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/zde37/Swift-Bank?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/zde37/Swift-Bank?style=default&color=0080ff" alt="repo-language-count">
<p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>

<br><!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary><br>

- [ Overview](#-overview)
- [ Features](#-features)
- [ Repository Structure](#-repository-structure)
- [ Modules](#-modules)
- [ Getting Started](#-getting-started)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Tests](#-tests)
- [ Project Roadmap](#-project-roadmap)
- [ Contributing](#-contributing)
- [ License](#-license)
- [ Acknowledgments](#-acknowledgments)
</details>
<hr>

##  Overview

<code>► INSERT-TEXT-HERE</code>

---

##  Features

<code>► INSERT-TEXT-HERE</code>

---

##  Repository Structure

```sh
└── Swift-Bank/
    ├── .github
    │   └── workflows
    ├── Makefile
    ├── README.md
    ├── config
    │   ├── application.go
    │   └── config.go
    ├── controller
    │   ├── handler
    │   └── middlewares
    ├── database
    │   ├── database.go
    │   └── migrations
    ├── go.mod
    ├── go.sum
    ├── helpers
    │   ├── helpers.go
    │   └── random.go
    ├── main.go
    ├── mock
    │   ├── repository_provider.go
    │   └── service_provider.go
    ├── models
    │   ├── entities.go
    │   └── requests.go
    ├── repository
    │   ├── repo.go
    │   ├── repo_impl.go
    │   ├── repo_main_test.go
    │   └── repo_test.go
    └── service
        ├── service.go
        ├── service_impl.go
        └── service_test.go
```

---

##  Modules

<details closed><summary>.</summary>

| File                                                                    | Summary                                                                                                                                                                                                                                                                                                                                                                                                     |
| ---                                                                     | ---                                                                                                                                                                                                                                                                                                                                                                                                         |
| [go.sum](https://github.com/zde37/Swift-Bank/blob/master/go.sum)     | This code file, part of the Swift-Bank repository, plays a critical role in configuring the application settings. By managing key configurations for the project, it enables smooth operation and customization of the banking application. This aspect of the codebase is pivotal for tailoring the applications behavior and ensuring optimal performance within the overarching repository architecture. |
| [main.go](https://github.com/zde37/Swift-Bank/blob/master/main.go)   | Initiates server setup by loading config, connecting to Postgres, instantiating repository, service, and handler for HTTP endpoints, then starting the server. This orchestrates the core functionality of Swift-Banks backend system.                                                                                                                                                                      |
| [go.mod](https://github.com/zde37/Swift-Bank/blob/master/go.mod)     | Defines external dependencies and versions for the Swift Bank project. Validates the required packages for the application to function correctly.                                                                                                                                                                                                                                                           |
| [Makefile](https://github.com/zde37/Swift-Bank/blob/master/Makefile) | <code>► INSERT-TEXT-HERE</code>                                                                                                                                                                                                                                                                                                                                                                             |

</details>

<details closed><summary>models</summary>

| File                                                                                 | Summary                         |
| ---                                                                                  | ---                             |
| [entities.go](https://github.com/zde37/Swift-Bank/blob/master/models/entities.go) | <code>► INSERT-TEXT-HERE</code> |
| [requests.go](https://github.com/zde37/Swift-Bank/blob/master/models/requests.go) | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>database</summary>

| File                                                                                   | Summary                         |
| ---                                                                                    | ---                             |
| [database.go](https://github.com/zde37/Swift-Bank/blob/master/database/database.go) | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>database.migrations</summary>

| File                                                                                                                              | Summary                         |
| ---                                                                                                                               | ---                             |
| [000001_init_schema.down.sql](https://github.com/zde37/Swift-Bank/blob/master/database/migrations/000001_init_schema.down.sql) | <code>► INSERT-TEXT-HERE</code> |
| [000001_init_schema.up.sql](https://github.com/zde37/Swift-Bank/blob/master/database/migrations/000001_init_schema.up.sql)     | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>repository</summary>

| File                                                                                                 | Summary                         |
| ---                                                                                                  | ---                             |
| [repo_main_test.go](https://github.com/zde37/Swift-Bank/blob/master/repository/repo_main_test.go) | <code>► INSERT-TEXT-HERE</code> |
| [repo_test.go](https://github.com/zde37/Swift-Bank/blob/master/repository/repo_test.go)           | <code>► INSERT-TEXT-HERE</code> |
| [repo_impl.go](https://github.com/zde37/Swift-Bank/blob/master/repository/repo_impl.go)           | <code>► INSERT-TEXT-HERE</code> |
| [repo.go](https://github.com/zde37/Swift-Bank/blob/master/repository/repo.go)                     | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>.github.workflows</summary>

| File                                                                                  | Summary                         |
| ---                                                                                   | ---                             |
| [ci.yml](https://github.com/zde37/Swift-Bank/blob/master/.github/workflows/ci.yml) | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>config</summary>

| File                                                                                       | Summary                         |
| ---                                                                                        | ---                             |
| [application.go](https://github.com/zde37/Swift-Bank/blob/master/config/application.go) | <code>► INSERT-TEXT-HERE</code> |
| [config.go](https://github.com/zde37/Swift-Bank/blob/master/config/config.go)           | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>mock</summary>

| File                                                                                                     | Summary                         |
| ---                                                                                                      | ---                             |
| [service_provider.go](https://github.com/zde37/Swift-Bank/blob/master/mock/service_provider.go)       | <code>► INSERT-TEXT-HERE</code> |
| [repository_provider.go](https://github.com/zde37/Swift-Bank/blob/master/mock/repository_provider.go) | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>controller.handler</summary>

| File                                                                                                               | Summary                         |
| ---                                                                                                                | ---                             |
| [handler.go](https://github.com/zde37/Swift-Bank/blob/master/controller/handler/handler.go)                     | <code>► INSERT-TEXT-HERE</code> |
| [handler_test.go](https://github.com/zde37/Swift-Bank/blob/master/controller/handler/handler_test.go)           | <code>► INSERT-TEXT-HERE</code> |
| [handler_main_test.go](https://github.com/zde37/Swift-Bank/blob/master/controller/handler/handler_main_test.go) | <code>► INSERT-TEXT-HERE</code> |
| [handler_impl.go](https://github.com/zde37/Swift-Bank/blob/master/controller/handler/handler_impl.go)           | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>controller.middlewares</summary>

| File                                                                                       | Summary                         |
| ---                                                                                        | ---                             |
| [jwt.go](https://github.com/zde37/Swift-Bank/blob/master/controller/middlewares/jwt.go) | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>service</summary>

| File                                                                                          | Summary                         |
| ---                                                                                           | ---                             |
| [service_impl.go](https://github.com/zde37/Swift-Bank/blob/master/service/service_impl.go) | <code>► INSERT-TEXT-HERE</code> |
| [service.go](https://github.com/zde37/Swift-Bank/blob/master/service/service.go)           | <code>► INSERT-TEXT-HERE</code> |
| [service_test.go](https://github.com/zde37/Swift-Bank/blob/master/service/service_test.go) | <code>► INSERT-TEXT-HERE</code> |

</details>

<details closed><summary>helpers</summary>

| File                                                                                | Summary                         |
| ---                                                                                 | ---                             |
| [helpers.go](https://github.com/zde37/Swift-Bank/blob/master/helpers/helpers.go) | <code>► INSERT-TEXT-HERE</code> |
| [random.go](https://github.com/zde37/Swift-Bank/blob/master/helpers/random.go)   | <code>► INSERT-TEXT-HERE</code> |

</details>

---

##  Getting Started

**System Requirements:**

* **Go**: `version x.y.z`

###  Installation

<h4>From <code>source</code></h4>

> 1. Clone the Swift-Bank repository:
>
> ```console
> $ git clone https://github.com/zde37/Swift-Bank
> ```
>
> 2. Change to the project directory:
> ```console
> $ cd Swift-Bank
> ```
>
> 3. Install the dependencies:
> ```console
> $ go build -o myapp
> ```

###  Usage

<h4>From <code>source</code></h4>

> Run Swift-Bank using the command below:
> ```console
> $ ./myapp
> ```

###  Tests

> Run the test suite using the command below:
> ```console
> $ go test
> ```

---

##  Project Roadmap

- [X] `► INSERT-TASK-1`
- [ ] `► INSERT-TASK-2`
- [ ] `► ...`

---

##  Contributing

Contributions are welcome! Here are several ways you can contribute:

- **[Report Issues](https://github.com/zde37/Swift-Bank/issues)**: Submit bugs found or log feature requests for the `Swift-Bank` project.
- **[Submit Pull Requests](https://github.com/zde37/Swift-Bank/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.
- **[Join the Discussions](https://github.com/zde37/Swift-Bank/discussions)**: Share your insights, provide feedback, or ask questions.

<details closed>
<summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your github account.
2. **Clone Locally**: Clone the forked repository to your local machine using a git client.
   ```sh
   git clone https://github.com/zde37/Swift-Bank
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to github**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.
8. **Review**: Once your PR is reviewed and approved, it will be merged into the main branch. Congratulations on your contribution!
</details>

<details closed>
<summary>Contributor Graph</summary>
<br>
<p align="center">
   <a href="https://github.com{/zde37/Swift-Bank/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=zde37/Swift-Bank">
   </a>
</p>
</details>

---

##  License

This project is protected under the [SELECT-A-LICENSE](https://choosealicense.com/licenses) License. For more details, refer to the [LICENSE](https://choosealicense.com/licenses/) file.

---

##  Acknowledgments

- List any resources, contributors, inspiration, etc. here.

[**Return**](#-overview)

---
