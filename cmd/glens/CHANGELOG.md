# Changelog

## [0.0.2](https://github.com/aydabd/glens/compare/cmd/glens/v0.0.1...cmd/glens/v0.0.2) (2026-02-28)


### Features

* extract pkg/logging as standalone workspace module with own go.mod and CI ([b032cfa](https://github.com/aydabd/glens/commit/b032cfa16c4179140286f7928c2f7ddc0f06c192))
* local open-source LLM support via Ollama (Mistral, Llama, Phi, Gemma) ([#9](https://github.com/aydabd/glens/issues/9)) ([0151fa6](https://github.com/aydabd/glens/commit/0151fa637502bddb52e4ec2262b6e75d86135166))
* move glens CLI into cmd/glens as isolated workspace module ([a146542](https://github.com/aydabd/glens/commit/a146542a1d519ce46fa29d684324ff51a1aab5d1))


### Bug Fixes

* errcheck/revive lint errors + per-module Makefiles + CI uses make targets ([10e499e](https://github.com/aydabd/glens/commit/10e499ee709571eaeece5fc6cb1ce6a7f15bc084))
* release workflow never triggered — bare semver tags not matched ([2cbfe7d](https://github.com/aydabd/glens/commit/2cbfe7da837776258771153e3f55eee33aa53b22))
* replace hardcoded 0.75 SuccessRate with real passed/total calculation ([45788d4](https://github.com/aydabd/glens/commit/45788d489daffb5202b2a873b126ad54f70e67a7))
* replace pre-commit in CI lint jobs with direct go commands, fix reporter_test.go formatting ([338edfc](https://github.com/aydabd/glens/commit/338edfcbeadffc628f87c57fd6e6efd7800bc994))
