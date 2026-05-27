# Contributing

## Testing approach

1. Each new feature or bugfix must be covered with unit-tests.
1. Some features and bugs should require implementation of smoke test.

We try to keep smoke tests small and fast, thus in typical workflow most of your
work should be dedicated to unit tests and only few things should be done with
smoke tests. This concept is well known as [Test Pyramid][test-pyramid].

Typically it is difficult to decide when unit tests are enough or when smoke
test is required. In this project we try to solve the question in the following
way:

1. If feature is very core part of business logic (e.g. online recognition of
   audio), write a smoke test.
1. If feature/bugfix could not be reproduced without real interaction with other
   service (e.g. with `Triton server`), write a smoke test.
1. If feature/bugfix behavior is tricky, hard to mock or requires implementation
   of complex logic inside fake objects, write a smoke test.
1. In all other cases please prefer unit test.

[test-pyramid]: https://martinfowler.com/articles/practical-test-pyramid.html
