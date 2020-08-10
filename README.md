# auth0-exercise
`auth0-exercise` is a go service that allows users to provision a new service based on their specs

There are two parts to the service
1. The api that handles a bulk of the service provisioning logic
2. The cli that allows for interactions with the service

## Running

### GH Token
Since the service interacts with the Github API, you'll need a valid access token.

[creating-a-personal-access-token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)

### Create Directory

Make a new project directory where the service will be created for you
```
mkdir -p $GOPATH/src/github.com/{GH_USERNAME}/{PROJECT_NAME}
```

```
cd $GOPATH/src/github.com/{GH_USERNAME}/{PROJECT_NAME}
```

### Creater spec file (optional)

providing a spec file allows you to customize the service being created, otherwise, the default options will be used.

> Sample.yaml
```yaml
 name: "custom-yaml-repo"
 owner: "Juan"
 version: "1.0.0"
 hasCopyright: true
 hasLicense: true
 description: "A custom auth0 service from a yaml spec"
 entrypoint: "custom-yaml-service"
 hasGitignore: true
 isPrivate: true
 imports: '"fmt"'
 mainBranch: "main"
```

### Run CLI Command

For a default service
```
auth0-exercise generate --token={GH_TOKEN}
```

For a configured service
```
auth0-exercise generate --token={GH_TOKEN} --file={spec.yaml OR spec.json}
```

## Installation

## Architecture

## Dependencies

- `github.com/google/go-github`
- `golang.org/x/oauth2`
- `github.com/urfave/cli/v2`
- `gopkg.in/yaml.v2`
- `github.com/sirupsen/logrus`

## License
`auth0-exercise` is licensed under the MIT License. Please see the LICENSE file for details.

## Roadmap
- [ ] Create PR from the templated service
- [ ] Support permissions throughout the process
- [ ] Support GH Repo configuration
- [ ] Support editing/deleting in case of a mistake
- [ ] Service creation progress updates
- [ ] VSCode Extension