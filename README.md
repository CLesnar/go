# go

Go Library

Useful tools, middleware, types, and more, used for building web apps, commands, and various software.

## Dev-Container
This repo uses a dev container for development. Typical use case is in VS Code install the 'Dev Containers' extension reopen the project in the container specified in /.devcontainer/docker-compose.yml.

## Repo Structure
This repo is setup for multiple commands, packages, tools, assets, and more. Each command, package can import the various packages and make use of them.

### Install software
The dev container should be up to date, but this is how to recreate it with all tools.

#### Homebrew (brew)
sudo /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" 