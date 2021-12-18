# proman

proman is a tool built to manage projects on a developers machine.

## Features

- [x] Configure project directory
- [x] List projects in project directory
- [x] Derive VCS repository URL

## Planned Features

- [ ] Welcome configuration screen
- [ ] Configure IDEs
- [ ] Open project directory with IDE
- [ ] Auto-detect IDEs on system
- [ ] Configure project name
- [ ] Configure project description
- [ ] View [git](https://git-scm.com/) configurations
- [ ] Edit [git](https://git-scm.com/) configurations

## Configuration

Configuration is stored in `~/.config/proman` as `config.toml`. Proman will create this file for you and only requires 
that `project_directory` be set.

## Development

### Prerequisites

- NodeJS
- Go >1.16
- [Wails](https://wails.app)

### Getting Started

Setup the frontend

```shell
cd frontend
npm i 
```

In another shell

```shell
go mod tidy
```

To run the backend

```shell
wails serve
```

To run the frontend

```shell
npm run serve
```

To build the app

```shell
wails build
```