# MockServer

There's a lot going on here. We are attempting to mock GitHub.com, CircleCI.com,
and Git CLI remote responses. For CircleCI.com and GitHub.com APIs we can simply
capture responses and modify them suit our needs.

The tricky part is mocking the GIT CLI responses that come from GitHub.com. For
that we use a Docker container running a Go server Mux that serves as a proxy
to Git http-backend. This seems to work pretty well to fool Git that it has made
contact with github.com and recieved a response. Though I would guess they may
be a bit wonky without knowing detailed specs. But since there is no easy way
to capture those, this will do.

1. Spin up a container that has a network configured so that any Github.com or
CircleCI.com request will not traverse the internet, but rather serve those
responses.
2. Setup Git http-backend in the same container.
3. Run the avr application tests in this container.
4. Modify responses to suit test needs trying to make them as accurate as
possible.
