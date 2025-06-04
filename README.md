# Labraboard: The Intelligent IaC Management Platform

![Labraboard](.img/logo.png)
## About Labraboard

Labraboard is a state-aware Infrastructure as Code (IaC) platform specifically designed to manage Terraform configurations with ease and efficiency. Inspired by the Labrador, a versatile and intelligent working dog known for its roles as a guide, rescuer, and retriever, Labraboard embodies these qualities by becoming an indispensable tool in your infrastructure management toolkit.

## Key Features

* State Management: Efficiently handle Terraform states with support for InMemory storage, HTTP backend, PostgreSQL backend, and Redis queue.
* Lock Handling: Robust mechanisms to handle state locks, ensuring smooth and safe infrastructure updates.
* Custom Environment and Variables: Tailor your Terraform and OpenTofu environments with custom configurations for seamless deployments.
* Dynamic Backend Configuration: Automatically override backend configurations during plan or apply operations, simplifying the setup process.
* Project Operations: Manage projects with create, read, update, and delete functionalities, including the ability to specify Git references, set default values, and configure other essential parameters.
* Time Lease for State: Implement time-based leases to manage state longevity.
* User Configuration: Customize user settings and preferences.

## Why Labraboard?
Terraform is a powerful tool for managing infrastructure, but handling governance and surrounding processes can be challenging. Labraboard addresses these challenges, making it easier to run plans and apply infrastructure changes efficiently. With Labraboard, you can focus on building and deploying your infrastructure without worrying about the complexities of governance.

Labraboard - inspired by the intelligence and versatility of a Labrador, is here to be your guide in the world of infrastructure management.

## Demo
![Spped up demo](https://github.com/PawelHaracz/labraboard/assets/14162492/c2f6f8ab-9c3e-4419-8ccf-45582b602639)

if it is too fast there is a slower bit version: [link](https://github.com/PawelHaracz/labraboard/assets/14162492/74c482b9-05fd-4e53-9cd8-da1b46207836)

### .HTTP files

You can try by own this tool using [app.rest file](api.rest) please bear aware that you have to change values in file `http-client.env.json` and add a private value `ARM_CLIENT_SECRET` for azure. 
I enhance you to test in your terraform and your env variables :)

## Starting point

Swagger docs has to be updated to reflect the new endpoints. 
It can be done by using command line `swag init -g ./cmd/api/main.go -o ./docs`

## Swagger Page
[Swagger API link](https://api.labraboard.dev)

Upps here should be demo video - no worry, I will be soon! - Now check the API on [https://api.labraboard.dev](https://api.labraboard.dev)

### Checking changes 
`git log --pretty=format:"%h%x09%an%x09%ad%x09%s"`

## RoadMap list
- [X] Reading plan
- [X] InMemory storage
- [X] Trigger run plan
- [X] Override backend
- [X] Use custom Env and variables on to terraform
- [X] Http Backend (Get Put)
- [X] Handle Locks on the state 
- [X] Handle Destroy
- [X] Add PostgreSQL as backend
- [X] Redis queue
- [X] Project CRUD
- [x] Add reference to git sha in aggregate and relate with plans
- [x] Add mapping to decouple mapping between aggregates and repository
- [X] Add and remove env variables
- [X] Add and remove  variables
- [X] Fix handling returning changes from plan
- [X] Add unit testing of aggregates
- [X] ~~Bug fixing e2e testing~~ (manually)
- [X] Path for parameters in git folder
- [X] Create handlers cmd
- [X] Time Lease for the state
- [X] Refactor and move to use interfaces in handlers
- [X] Implement Logging
  - [X] Logger and move init code to init function
  - [X] Integrate logger with gin - use middleware for recordId, in future userId
  - [X] Propagate context values between loggers
  - [X] Replace all print to logger
  - [X] Add every method ctx to enrich logger
  - [X] Integrate handlers to use logger
- [X] Handle scheduled plan in TerraformPlanner
- [X] Run plan using http backend 
- [ ] Access Token for Backend http
- [ ] Clean solution to be more DDD
- [ ] Create Plan changes during run, what was happened during the time
- [X] Apply Mechanism to handle the state
  - [X] Apply based on the Plan
  - [X] Save outputs as deployment, handle errors
- [ ] Backup before apply using ApplyOptions.Backup
- [ ] Implement retries on apply
- [ ] Correlate Project, Deployment, Plans
- [ ] Integrate with the Git
- [ ] Handle other version than version 4.0 of tf 
- [ ] Handle multiple version of tf and tofu
- [ ] Policies and run on pre/post plan/apply
- [ ] Authenticate
- [ ] User configuration
- [ ] Add a web interface
- [ ] Encryption at rest
- [ ] fix end2end tests - `terraform_project_plan_test`
- [ ] fix bugs related passing refs instead of valu objects in array and fix mapping data from dao into aggregate

[//]: # (### Architecture )

[//]: # (#### Event Storming )

[//]: # ([Labraboard Event storming]&#40;https://miro.com/app/board/uXjVKHzpuQ4=/?share_link_id=741994614357&#41;)

[//]: # ()
[//]: # (![big-picture-1.png]&#40;.img/big-picture-1.png&#41;)

### Http Backend
Solution uses own delivered http backend where state is kept. During running plan or apply the backend configuration is 
added automatically by overriding the backend. 

#### Example of using own http Backend
please use your project id. Application use it identify terraform state.
```hcl
terraform {
  backend "http" {
    address = "http://localhost:8080/api/v1/state/terraform/bee3cf56-ecd1-4434-8e18-02b0ae2950cc"
    lock_address = "http://localhost:8080/api/v1/state/terraform/bee3cf56-ecd1-4434-8e18-02b0ae2950cc/lock"
    unlock_address = "http://localhost:8080/api/v1/state/terraform/bee3cf56-ecd1-4434-8e18-02b0ae2950cc/lock"
  }
}
```
## Build
### Prerequisites
- Go 1.x
- Node.js and Yarn
- Docker and Docker Compose (optional)
- PostgreSQL
- Redis

### Required Development Tools
The following tools will be automatically installed during development setup:
- `swag` - Swagger documentation generator
- `gosec` - Security scanner
- `trivy` - Container security scanner

### Environment Variables
The application can be configured using the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| CONNECTION_STRING | PostgreSQL connection string | - | Yes |
| HTTP_PORT | HTTP port to serve the application | 8080 | No |
| REDIS_HOST | Redis host | localhost | No |
| REDIS_PORT | Redis port | 6379 | No |
| REDIS_PASSWORD | Redis password | eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81 | No |
| REDIS_DB | Redis database number | 0 | No |
| LOG_LEVEL | Logging level | 1 | No |
| USE_PRETTY_LOGS | Use pretty logs instead of JSON | false | No |
| SERVICE_DISCOVERY | Service discovery URL | http://localhost | No |
| FRONTEND_PATH | Path to frontend files | /app/client | No |

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/PawelHaracz/labraboard.git
cd labraboard
```

2. Update Go modules:
```bash
make mod
```

3. Build the application:
```bash
make build
```

4. Run tests:
```bash
make test
```

5. Generate documentation:
```bash
make build-swagger
```

### Troubleshooting

If you encounter dependency issues:

1. Clean the Go module cache:
```bash
go clean -modcache
```

2. Update dependencies:
```bash
make update-dependencies
```

### Development Workflow

The project uses Make targets to streamline the development process. Here are the main commands:

#### Development
- `make dev-setup` - Setup development environment
- `make install` - Install application and dependencies
- `make fmt` - Format code
- `make lint` - Run linter
- `make vet` - Run go vet
- `make dependency-check` - Check for outdated dependencies
- `make update-dependencies` - Update dependencies

#### Testing
- `make test` - Run all tests
- `make test-unit` - Run unit tests
- `make test-cover` - Generate test coverage report

#### Building
- `make build` - Build all components
- `make build-api` - Build API server
- `make build-handlers` - Build handlers
- `make build-frontend` - Build frontend application

#### Docker Operations
- `make docker-build` - Build Docker image
- `make docker-push` - Push Docker image
- `make docker-compose-up` - Start services
- `make docker-compose-stop` - Stop services

#### Security
- `make security-scan` - Run security scans

#### Release Management
- `make release-prepare` - Prepare release artifacts
- `make release-publish` - Publish release

#### Maintenance
- `make clean` - Clean build artifacts
- `make clean-all` - Remove all generated artifacts

For a complete list of available commands, run:
```bash
make help
```

### Development Guidelines

1. **Code Style**
   - Follow Go standard formatting
   - Run `make fmt` before committing
   - Ensure code passes `make lint` and `make vet`

2. **Testing**
   - Write unit tests for new features
   - Maintain test coverage above 80%
   - Run `make test` before committing

3. **Documentation**
   - Update Swagger documentation for API changes
   - Keep README.md up to date
   - Document new features and changes

4. **Security**
   - Run security scans regularly
   - Keep dependencies updated
   - Follow security best practices

5. **Release Process**
   - Update version in Makefile
   - Create release notes
   - Run full test suite
   - Build and test Docker image
   - Publish release

## Disclaimer

Please note that this project is currently under active development and is not considered production-ready. We are continuously working to improve and stabilize its features, but it does not yet meet all the requirements for production use.