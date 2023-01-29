# OpenFaaS Functions

<!-- toc -->

* [Open in devbox](#open-in-devbox)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

My collection of [OpenFaaS](https://openfaas.com) functions

## Open in devbox

* `curl -fsSL https://get.jetpack.io/devbox | bash`
* `devbox shell`
* `make serve`

After a few moments, you will have a running Kubernetes cluster with the
OpenFaaS gateway available on [localhost:8080](http://localhost:8080). To get the
username and password, run `make openfaas_login`.

**IMPORTANT** You will need to install [Skaffold](https://skaffold.dev/docs/install/)
manually. The version in nix is `2.0.2` which has a
[bug](https://github.com/GoogleContainerTools/skaffold/issues/8243). This is fixed
`2.0.4` and above.
