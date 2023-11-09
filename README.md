# Technical Writer Application

This application automates the generation of documentation for code files using the OpenAI API.

## Overview

The Technical Writer application interfaces with the OpenAI API to generate documentation for each file within a given codebase. It reads the contents of the source files and utilizes a predefined instruction set to create relevant and helpful package-level documentation for external developers.

## Getting Started

### Prerequisites

- Docker must be installed on your system.
- You need an OpenAI API key with access to the GPT-4 API.

### Setup

1. Clone the repository to your local environment or pull the Docker image for this application.
2. Copy your OpenAI API key into your environment as `OPENAI_API_KEY`.

### Running the Application

To run the application:

1. Navigate to the directory containing the code base that you wish to document.
2. Execute the built Docker image and ensure the `OPENAI_API_KEY` environment variable is appropriately sourced.

    ```sh
    docker run -e OPENAI_API_KEY=${OPENAI_API_KEY} technical-writer
    ```

The application will then process all source files and apply generated documentation where applicable. 

### Development

If you wish to develop or modify the application, the source files include:

- `main.go`: Contains the core logic for file processing and API interaction.
- `action.yaml`: Describes the GitHub Action configuration for deployment.
- `Dockerfile`: Defines the Docker container configuration for building and running the application.

To build the application for local development:

```sh
go build -o technical-writer main.go
```

### Extending Functionality

To add more features or customize the way documentation is generated, you can modify `main.go` to alter the instruction set given to OpenAI or to change how files are processed.

## License

This application is open source and available under the MIT License.