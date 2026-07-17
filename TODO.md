# TODO

- [x] execute `go mod tidy`
- [x] install a library to send emails
- [x] research if is needed an index for the system user
- [x] research how to create an index for nullable fields to faster null entries
- [ ] research how to assign a weight to roles
- [ ] research for good restrictions of the length of name, email, password fields
- [x] postgresql has default values when NULL in variable assignment? ([file](./db/migrations/000002_users_table.up.sql))
- [x] research about value objects
- [ ] add `application/json` `Content-Type` to the response
- [ ] add middleware to reject non-`application/json` `Content-Type` in requests
- [ ] implement graceful shutdown of the server
- [ ] implement **timeout** middleware to prevent the user to wait for long tasks
- [ ] implement **limit** middleware to prevent the system be called multiple times from the same IP
