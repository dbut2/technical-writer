# Technical Writer Bot

This repository contains the source code for a GitHub Action named *Technical Writer*. The action is designed to interact with the OpenAI API to generate technical documentation for files provided to it. External developers can integrate this action into their workflows to automatically generate package-level documentation for their projects.

## Action Specification: `action.yaml`

This is the action metadata file that configures the *Technical Writer* GitHub Action.

- `name`: The name of the GitHub Action, which is *Technical Writer*.
- `description`: The description of what the GitHub Action does, in this case using the OpenAI API to create technical documentation.
- `inputs`: The inputs section lists the expected input parameters for the action. Here, `openai_api_key` is a required input parameter which must be provided to use the OpenAI API.
- `runs`: This section defines the runtime environment for the action. It specifies that the action uses a Docker container and points to the Dockerfile used for building that container. Environmental variables, such as `OPENAI_API_KEY`, are also specified.

## Docker Configuration: `Dockerfile`

This is a multi-stage Dockerfile used to build the technical writer bot.

- It starts by creating a build environment using the `golang:alpine` base image.
- The application's Go module files (`go.mod`, `go.sum`) and the main application file (`main.go`) are copied into the image.
- The application is built, resulting in a binary named `/bin/technical-writer`.
- A second, lighter-weight `alpine:latest` image is used as a final stage.
- The built binary is copied over from the builder stage, and the entry point for the container is set to the path of this binary.

## Go Module Definition: `go.mod`

The `go.mod` file defines the application's module path and its dependencies. Dependencies are specified with versions:
- `github.com/sashabaranov/go-openai v1.17.4`: The OpenAI Golang SDK for communicating with OpenAI's API.
- `golang.org/x/sync v0.5.0`: Additional synchronization primitives for Go.

## Main Application: `main.go`

The `main.go` file is the entry point of the application, which:
- Sets up the OpenAI client using the API key retrieved from the environment variable.
- Lists all relevant files in the `/github/workspace` directory to be documented.
- Calls the `document` function to generate technical documentation for each listed file using OpenAI API.
- Implements file I/O operations such as reading from and writing to files, as needed.
- Provides two main functions, `listAllFiles` and `document`, to perform file listing and documentation generation, respectively.

## Usage

To use this GitHub Action in your workflow:

1. Set up a secret in your repository with your OpenAI API key.
2. Use the secret in your workflow file when configuring the action.
3. Add the action to a step in your GitHub workflow YAML file.

Sample workflow step:
```yaml
- name: Generate Documentation
  uses: your-repo/technical-writer@main
  with:
    openai_api_key: ${{ secrets.OPENAI_API_KEY }}
```

## Contributing

External developers are encouraged to contribute to this project by submitting pull requests and reporting issues or suggestions under the repository's [Issues](https://github.com/dbut2/technical-writer/issues) section.

For more significant changes, please raise an issue to discuss the changes before submitting a pull request. This will ensure that your valuable efforts align with the project direction and that merging can occur smoothly.

## License

Please include your project license here.

---
*Note: This README assumes the repository URL is `https://github.com/dbut2/technical-writer` and should be adjusted to match the actual repository URL.*