# labraboard

![Labraboard](.img/logo.png)
## About Labraboard

Labraboard is a state-aware Infrastructure as Code (IaC) platform that focuses on managing Terraform. It is designed to streamline the process of running plans and applying infrastructure changes. 

The platform provides a robust set of features including InMemory storage, HTTP backend, PostgreSQL backend, Redis queue, and more. It also supports handling locks on the state, destroying, and overriding the backend.

Labraboard is equipped with a custom environment and variables for Terraform, OpenTofu, making it a versatile tool for managing and deploying infrastructure. It also offers an HTTP backend where the state is kept, and during the running plan or apply, the backend configuration is added automatically by overriding the backend.

Future enhancements include project CRUD, time lease for the state, integration with Git, user configuration, web interface, and more.

With Labraboard, you can efficiently manage your infrastructure and streamline your IaC processes. 

## Starting point

Swagger docs has to be updated to reflect the new endpoints. 
It can be done by using command line `swag init -g ./cmd/main.go -o ./docs`

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
- [ ] ~~Add unit testing~~
- [ ] ~~Bug fixing e2e testing~~
- [ ] Path for parameters in git folder
- [ ] Time Lease for the state
- [ ] Handle scheduled plan in TerraformPlanner
- [ ] Access Token for Backend http
- [ ] Run plan using http backend
- [ ] Integrate plan with terraform plan association
- [ ] Clean solution to be more DDD
- [ ] Handle other version than version 4.0 of tf
- [ ] Apply Mechanism to handle the state
- [ ] Integrate with the Git
- [ ] Handle multiple version of tf and tofu
- [ ] Policies and run on pre/post plan/apply
- [ ] Authenticate
- [ ] User configuration
- [ ] Add a web interface
- [ ] Encryption at rest

### Http Backend
Solution uses own delivered http backend where state is kept. During running plan or apply the backend configuration is 
added automatically by overriding the backend. 

#### Example of using own http Backend
```hcl
terraform {
  backend "http" {
    address = "http://localhost:8080/api/v1/state/terraform/bee3cf56-ecd1-4434-8e18-02b0ae2950cc"
    lock_address = "http://localhost:8080/api/v1/state/terraform/bee3cf56-ecd1-4434-8e18-02b0ae2950cc/lock"
    unlock_address = "http://localhost:8080/api/v1/state/terraform/bee3cf56-ecd1-4434-8e18-02b0ae2950cc/lock"
  }
}
```

## Links

1. https://www.squash.io/optimizing-gin-in-golang-project-structuring-error-handling-and-testing/
2. https://github.com/swaggo/gin-swagger
3. https://github.com/eddycjy/go-gin-example
4. https://github.com/derekahn/ultimate-go/blob/master/language/interfaces/main.go
5. https://github.com/percybolmer/ddd-go
6. https://velocity.tech/blog/build-a-microservice-based-application-in-golang-with-gin-redis-and-mongodb-and-deploy-it-in-k8s
7. https://www.ompluscator.com/article/golang/practical-ddd-domain-repository/?source=post_page-----d308c9d79ba7--------------------------------

## Disclaimer

Please note that this project is currently under active development and is not considered production-ready. We are continuously working to improve and stabilize its features, but it does not yet meet all the requirements for production use.
