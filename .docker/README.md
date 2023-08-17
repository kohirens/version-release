# Container Environment

Container environments are provided for releases but also for developing.

The `vr` container image is used for releasing and running test against.

There is also a more complex mock-server. It allows testing features that
call GitHub or CircleCI, mocking the responses expected, thus allowing
observation of application behavior.

## How To Run Mock Container

The mock container works in tandem with the vr (Version Release) container. You
run them so that the container network points to the mock-server for request
going to public GitHub and CircleCI. This allows those responses to be mocked
and the Version Release application tested for how it handles those responses.

```shell
docker run -it `
    --rm `
    --add-host "api.circleci.com:127.0.0.1" `
    --add-host "app.circleci.com:127.0.0.1" `
    --add-host "github.com:127.0.0.1" `
    --add-host "api.github.com:127.0.0.1" `
    --env-file .docker/mock-server/integration-test.env `
    mock-server
```
