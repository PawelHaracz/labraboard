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
- [ ] Add postgresql as backend
- [ ] Time Lease for the state
- [ ] Handle Destroy
- [ ] Access Token for Backend http
- [ ] Run plan using http backend
- [ ] Integrate plan with terraform plan association
- [ ] Handle other version than version 4.0 of tf
- [ ] Apply Mechanism to handle the state
- [ ] Integrate with the Git
- [ ] Handle multiple version of tf and tofu
- [ ] Policies and run on pre/post plan/apply
- [ ] Authenticate
- [ ] User configuration
- [ ] Add a web interface