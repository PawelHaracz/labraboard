# labraboard

## Starting point

Swagger docs has to be updated to reflect the new endpoints. 
It can be done by using command line `swag init -g ./cmd/main.go -o ./docs`


### Checking changes 
`git log --pretty=format:"%h%x09%an%x09%ad%x09%s"`

## TODO list
- [X] Reading plan
- [X] InMemory storage
- [X] Trigger run plan
- [X] Override backend
- [X] Use custom Env and variables on to terraform
- [X] Http Backend (Get Put)
- [X] Handle Locks on the state 
- [X] Handle Destroy
- [ ] Add postgresql as backend
- [ ] Time Lease for the state
- [ ] Handle scheduled plan in TerraformPlanner
- [ ] Access Token for Backend http
- [ ] Run plan using http backend
- [ ] Integrate plan with terraform plan association
- [ ] Clean solution to more DDD
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