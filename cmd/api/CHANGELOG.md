# Changelog

## [0.0.2](https://github.com/aydabd/glens/compare/cmd/api/v0.0.1...cmd/api/v0.0.2) (2026-02-28)


### Features

* add Backend REST API module (cmd/api) — Phase 1 ([41b5e87](https://github.com/aydabd/glens/commit/41b5e87e790338dba6cc423eae4b6b53681cc9b8))
* add issue tracker provider abstraction (Phase 17) ([e09126a](https://github.com/aydabd/glens/commit/e09126a36ed601835db284689baba99c94999431))
* add Terraform fmt-check and per-module Go linting to root pre-commit ([e11fea2](https://github.com/aydabd/glens/commit/e11fea2e01603298bfee1304a35a05112aa3d09c))
* **api:** add endpoint safety categoriser (Phase 8) ([15f4a24](https://github.com/aydabd/glens/commit/15f4a246c4e37dd82e5a2c1ad83fc73ca4ffc4d1))
* **api:** add event schema definitions and publisher interface (Phase 11) ([da1a074](https://github.com/aydabd/glens/commit/da1a074928930a065d9d82b5b29aaf62c31309a1))
* RFC 9457 Problem Details error responses + Phase 18 auth plan ([222d49c](https://github.com/aydabd/glens/commit/222d49c0d3a8d10584e2591bec253157b8b37a54))


### Bug Fixes

* address all copilot reviewer feedback ([ff20aa8](https://github.com/aydabd/glens/commit/ff20aa83222de5afb79ae1320cd80bfc6f1c7374))
* gofmt struct alignment in analyze.go/mcp.go and use Go 1.25 in api.yml ([090cc8d](https://github.com/aydabd/glens/commit/090cc8ddd57e62f7080b6f4c961d711a07857684))
* keep MCP endpoint using JSON-RPC 2.0 error format (not RFC 9457) ([670f712](https://github.com/aydabd/glens/commit/670f712114c13d025e5d7f480a71ab5b09bfe522))
* resolve API lint (gosec G114 + revive) and Terraform validate failures ([fb61d16](https://github.com/aydabd/glens/commit/fb61d160af259c06565582904093131f1433aa32))
