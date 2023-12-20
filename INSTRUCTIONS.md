## Evaluation points in order of importance

- use of clean code, which is self documenting
- use of packages to achieve separation of concerns
- use of domain driven design
- use of golang idiomatic principles
- use of docker
- tests for business logic
- use of code quality checks such as linters and build tools
- use of git with appropriate commit messages
- documentation: README and inline code comments
- you MUST use go modules and a version of go >= 1.19

Time limitations: there are no hard time limits. Once again, do not spend more than 2 or 3 hours.

## Technical test

- Given a file with ports data (ports.json), write a port domain service that either creates a new record in a database, or updates the existing one (Hint: no need for delete or other methods).
- The file is of unknown size, it can contain several millions of records, you will not be able to read the entire file at once.
- The service has limited resources available (e.g. 200MB ram).
- The end result should be a storage containing the ports, representing the latest version found in the JSON. (Hint: use an in memory database to save time and avoid complexity).
- A Dockerfile should be used to contain and run the service (Hint: extra points for avoiding code building in docker).
- Provide at least one example per test type that you think are needed for your assignment. This will allow the reviewer to evaluate your critical thinking as well as your knowledge about testing.
- Your readme.md should explain how to run your program and test it.
- The service should handle certain signals correctly (e.g. a TERM or KILL signal should result in a graceful shutdown).

## Bonus points

- Address security concerns for Docker
- Code structure according to hexagonal architecture
