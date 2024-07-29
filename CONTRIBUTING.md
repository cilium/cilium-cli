# Contributing to Cilium CLI

## IMPORTANT Cilium CLI is moving 🚚📦📦📦

We are planning to merge Cilium CLI code into cilium/cilium repository to
simplify the overall development process after Cilium v1.16.0 gets released.
See [CFP-25694](https://github.com/cilium/design-cfps/pull/9) for details.
If you have questions, post a message in
[#development Cilium Slack channel](https://cilium.slack.com/archives/C2B917YHE).

## Contribution workflow

Cilium CLI uses GitHub for collaborative development. Please use GitHub issues
to discuss proposals and use pull requests to suggest changes. For more
information see the [Cilium Development
Guide](https://docs.cilium.io/en/latest/contributing/development/).

## An important note about mutating cluster resources

Cilium CLI must follow these general rules regarding mutating cluster resources:

- Cilium CLI must not mutate resources that are not managed by the Cilium Helm
  chart. Cilium CLI must use Cilium Helm values to configure Cilium installations.
- The only exception to the rule above is the connectivity test command, which
  may mutate resources within Kubernetes namespaces with the name prefix specified
  by --test-namespace flag. The connectivity test must not mutate resources that
  affect network connectivity outside the test namespaces.
- TODO: I'd like to deprecate the use of --include-unsafe-tests flag for adding
  tests that mutate resources outside the test namespaces:
  - These unsafe tests typically can't run in parallel.
  - It sets a bad precedent for new contributors. They might see test cases in
    the connectivity test command that mutate ressources outside the test
    namespaces, and assume that the connectivity test is the right command to
    add more unsafe tests. It only takes one careless review for these tests
    that are not guarded by --include-unsafe-tests to slip in.
  - My current thinking is to add a seprate command, maybe call it "cilium test"
    or "cilium integration test" or something that prints a big warnig at the
    beginning and prompt the user whether to proceed to make it clear it will
    mess up your cluster.

If you find use cases that you believe warrent making an exception to these
rules, start a discussion either by posting a message in [Cilium & eBPF Slack]
`#development` channel, or by opening a GitHub issue.

## Slack

Most developers are using the [Cilium & eBPF Slack] outside of GitHub.

## Code of conduct

All members of the Cilium community must abide by the [Cilium Community Code of
Conduct](https://github.com/cilium/cilium/blob/main/CODE_OF_CONDUCT.md). Only
by respecting each other can we develop a productive, collaborative community.
If you would like to report a violation of the code of contact, please contact
any of the maintainers or our mediator, Beatriz Martinez <beatriz@cilium.io>.

[Cilium & eBPF Slack]: https://docs.cilium.io/en/latest/community/community/#slack
