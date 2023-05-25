# Mock Server

Mock a response to CircleCI, Git, and [GitHub API]. See [mock-server] code.

The GitHub and CircleCI mocking was not difficult as they are well documented.

Mocking Git response for remote repository fetch and pushes required at least
a weeks efforts to understand its Protocol and how to generate the correct
responses. The results is that all request are proxied from a Go server to
`git http-backend` via environment variables and command line arguments.

These links provided all that I needed to begin to understand how to mock Git
remote responses to pushes:

1. https://www.git-scm.com/docs/http-protocol
2. https://git-scm.com/docs/pack-protocol/en
3. https://git-scm.com/docs/git-upload-pack
4. https://git-scm.com/docs/gitprotocol-v2
5. https://stackoverflow.com/questions/38007361/is-there-anyway-to-create-null-terminated-string-in-go
6. https://git-scm.com/docs/git-receive-pack
7. https://git-scm.com/docs/git-http-backend
8. https://stackoverflow.com/questions/48472362/how-to-open-git-http-backend-as-a-http-server
9. https://github.com/git/git/blob/master/t/t5560-http-backend-noserver.sh

If more complex request need to be tested, it may be worth looking into using
git-daemon as a server and proxy request to that.

All this to mock responses and test that the Version Release Orb CLI works as
expected before performing real-world integration. These may seem like overkill
but for the amount of bugs that were caught before testing with partners and
having something to assert as code changes are happening is invaluable to
making a truly high-quality product. Not to mention just amount of confidence
you gain about the product you are delivering and how well you can speak to it
is just immeasurable. Don't every underestimate the power of automated testing.
The juice is worth the squeeze, unless the amount of bugs cause doe not make
sense for the amount of time spent getting the test to work; and they are very
brittle.

I can honestly say that the amount of time I spent on getting the Git mock
responses working was NOT worth it for this project alone. It just allowed me to 
go through all my test without having to write special code to get around the
git commands; so that I could test other parts that followed those command, thus
not reachable any other way during test. Because there is no way to mock a
package in Go, that I'm aware of. However, I will be able to use this of
knowledge of mocking git remote responses over many projects. Where it will
balance out over time.

---

[mock-server]: /vro/mock-server
[GitHub API]: https://docs.github.com/en/rest?apiVersion=2022-11-28
