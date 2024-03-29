# Tracified Blockchain

This project implements the tracified Blockchain-gateway.

# Git Workflow

Master branch is the main development branch. Do not commit directly to master branch.
Staging branch is for the staging environment which mimics production environment.

**TODO: figure out who, when and how to merge to staging branch**

Steps to contribute:

- Clone a local repo or if already cloned pull latest version with `git pull`
- Create a new feature branch for the work you're doing `git checkout -b feature-branch-name`
- Develop, Add, Commit as much as you like ðŸ˜Ž
- Once your feature is complete and all tests are passed, you can push your branch to bitbucket with `git push -u origin feature-branch-name`
- Goto bitbucket and create a new pull request to merge your feature branch
- Add a team member to review your pull request and ask to merge.

# Getting Started

- GO Installation
  https://docs.google.com/document/d/1yTEBXge5kwnq5u6AUcPC61OALerW1SFoAh6KqPme-MY/edit?usp=sharing

```

```

- Setting Environmental variables

```
    set GOPATH=E:\GO\
```

- Install packages

```
    go get -u github.com/golang/dep/cmd/dep
    dep init
    dep ensure


```

- Make env files

```
    Create a copy of default.env file.
    Rename it with your environment name and add your environment info
```

- Build

```
    go build
```

- Run

```
    ./tracified-gateway

    for windows run :- tracified-gateway.exe
```

- Run using Docker

```
    - with rabbitmq
    | ------------------------------------------------------------------- |
    steps
    create  docker network  tillitgatwaynetwork
    give hostname and network to rabbitmq container
    run rabbitmq
    build gateway
    run gateway

    docker network create tillitgatewaynetwork

    docker run -itd --rm --hostname rabbit-1 --name rabbit-1 --network tillitgatewaynetwork -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management

    docker build . -t [gateway tag name]

    docker run -itd --rm --net tillitgatewaynetwork  -e env=local -p 8000:8000 [gateway image id]
```

    | Note
    | ------------------------------------------------------------------- |
        In Local environment it should be
        RABBITUSER     = "guest"
        RABBITPASSWORD = "guest"
        RABBITHOSTNAME = "rabbit-1"
        RABBITPORT     = "5672"

        // Both rabbitmq container host nema and container name, and env file's RABBITHOSTNAME should be same

        env variable with docker
        RABBITHOSTNAME= "rabbit-1"

        without docker local environment
        RABBITHOSTNAME= "localhost"

```
    - old
    | ------------------------------------------------------------------- |
    docker run -e env=local -p 8000:8000 <image>
```

- Health Check

```
http://localhost:8000/health

```

# Folder Structure

| Name    | Description                                                |
| ------- | ---------------------------------------------------------- |
| **bin** | Contains the executable from the GO build.                 |
| **src** | Contains source code that will be compiled to the exe file |
| **pkg** | Contains packages that are/ has to be imported.            |

#Logging && Error

- Log errors exactly where they occur.
- DO NOT LOG class OR function names
- Pass errors to upper layers ONLY TO IDENTIFY WHAT WENT WRONG AND SEND CLIENT APPROPRIATE RESPONSE (not to log errors in upper layers)
- If we require to send message to upper layers regarding the type of errors, define them in the comments and pass that message In other instances if upper layers don't have to know exactly what types of erros can come in advance just pass them to upper layers.
