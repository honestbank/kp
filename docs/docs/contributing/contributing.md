---
sidebar_position: 1
---

# Contributing
Keeping the core of KP small enables us to maintain this project a little easier.
We aim evolve KP through middleware without changing the core API.

But we do aim to add more middlewares.
If you think a middleware would help large group of people then feel free to request it or send a PR.

:::warning
If you have found a security problem, please reach out to engineering@honestbank.com
:::


### Contribution Guidelines {#guidelines}
Any contribution that adds or improves middlewares without changing the core API will be more likely to be accepted.

If you have a new feature idea that's not applicable to large community, you can simply keep it in your own repository.

:::tip
Searching through existing issues/PRs before attempting to send a PR is always the best idea.
Even if you're sure no one is working on the same idea as yours, making sure we're onboard with what you're planning will save you some time.
:::

## Code Contributions {#code-contributions}
First clone this repository, enter `v2` directory and run `go mod download && go mod vendor` this should download all dependencies to your vendor directory.

Make sure you install pre-commit hook, and you follow the same standard the project is already following.
If you're in doubt, you can always discuss in GitHub.

## Running tests {#running-tests}
Because of the nature of this project, we focused heavily on integration tests of the core API.

If you're working on a middleware though, feel free to run unit tests on just the middleware package.

Otherwise, feel free to use docker compose to start entire stack and run all tests including integration tests using `integration` build flag.
