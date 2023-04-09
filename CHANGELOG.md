# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.4.0](https://github.com/padok-team/YATAS/compare/v1.3.3...v1.4.0) (2023-04-09)


### Features

* **cli:** added documentation ([5ef60b3](https://github.com/padok-team/YATAS/commit/5ef60b3ff9d90bc0816f68e9ac09d7395eafee4d))
* **config:** added validation ([c170337](https://github.com/padok-team/YATAS/commit/c170337a0af0b14f833c64fedc61d08aa4ead98c))
* **global:** remove all AWS references ([#102](https://github.com/padok-team/YATAS/issues/102)) ([d559814](https://github.com/padok-team/YATAS/commit/d559814ccf80f77a194c322d156531ba2f58f03e))
* **golang:** updated to 1.20 ([f16bf3d](https://github.com/padok-team/YATAS/commit/f16bf3d657249978b0d5295dc0dd8c90e6bea786))
* **internal:** improved logging ([71850f2](https://github.com/padok-team/YATAS/commit/71850f25e9919fdc1bc70b517cae94556881a6f9))
* **labeler:** updated labeler ([e9e503b](https://github.com/padok-team/YATAS/commit/e9e503ba282f28de1c316c04f1655ea5b2b8656a))
* **logger:** improved coloring and default to error ([3ccfcac](https://github.com/padok-team/YATAS/commit/3ccfcac75bcb663c8e0c71580277d9afbe907162))
* **logger:** moved to own package ([75cdb23](https://github.com/padok-team/YATAS/commit/75cdb23f5bc8a182add0b26737d1ae086b3d0b91))
* **manager:** improved logging and readability of code ([2920936](https://github.com/padok-team/YATAS/commit/2920936bf203a0d2ca525cf7aab56aece4f2df3a))
* **panic:** removed the last ones :) ([a675028](https://github.com/padok-team/YATAS/commit/a675028afe493b2393cee42e271eef9f9552f246))
* **releaser:** updated to latest ([399da1a](https://github.com/padok-team/YATAS/commit/399da1a7d470f2a070189fe0357876a99419b62d))
* **tests:** added test ([19f46bb](https://github.com/padok-team/YATAS/commit/19f46bb14f081841a598e99ad32736790de3b635))
* **tests:** added test for plugin installation ([ef17889](https://github.com/padok-team/YATAS/commit/ef17889fbc64925ca4dab17982dd717f8fbca257))
* **tests:** added tests on types shared so that it can be changed ([0a35d93](https://github.com/padok-team/YATAS/commit/0a35d9348c0133481848887748f541add7e9ff79))
* **types:** added type checking check ([8546be1](https://github.com/padok-team/YATAS/commit/8546be189cac38f6ea8eaf3cf38b4e6d25b1a1af))


### Bug Fixes

* **release:** updated to 1.20 ([7c241d2](https://github.com/padok-team/YATAS/commit/7c241d21f51bfff2423a0f305554c24138dcce3d))

### [1.3.3](https://github.com/padok-team/YATAS/compare/v1.3.2...v1.3.3) (2022-12-23)

### [1.3.2](https://github.com/padok-team/YATAS/compare/v1.3.1...v1.3.2) (2022-12-15)

### [1.3.1](https://github.com/padok-team/YATAS/compare/v1.3.0...v1.3.1) (2022-10-25)


### Bug Fixes

* fix yatas-aws plugin source ([#84](https://github.com/padok-team/YATAS/issues/84)) ([a7cc241](https://github.com/padok-team/YATAS/commit/a7cc2413fa9ed8cedaa84e8ce3615b34eb795432))
* **name:** changed name ([f822e58](https://github.com/padok-team/YATAS/commit/f822e580155ebd14f0ea218173e84864d1396d65))

## [1.3.0](https://github.com/padok-team/YATAS/compare/v1.2.0...v1.3.0) (2022-10-12)


### Features

* **yatas:** added debug log ([eb32a4c](https://github.com/padok-team/YATAS/commit/eb32a4cffa8f869611a59976ec76cce45a07329a))


### Bug Fixes

* **readme-gen:** added cognito because readme generator was wrong ([83a2c45](https://github.com/padok-team/YATAS/commit/83a2c453d896e3df00c8722d288c722369fca9de))

## [1.2.0](https://github.com/padok-team/YATAS/compare/v1.1.0...v1.2.0) (2022-09-27)


### Features

* **categories:** added categories in init check ([69c7afb](https://github.com/padok-team/YATAS/commit/69c7afba8ba7a5e46d738e207c455aff10c3cb5c))

## [1.1.0](https://github.com/padok-team/YATAS/compare/v1.0.0...v1.1.0) (2022-09-27)


### Features

* **checks:** update init check method to add new category initialisation ([87c0230](https://github.com/padok-team/YATAS/commit/87c0230b38c898e081737b30e04a5ccecb3f9223))

## [1.0.0](https://github.com/padok-team/YATAS/compare/v0.11.10...v1.0.0) (2022-09-26)


### Features

* **plugins:** added better error message ([2f8d566](https://github.com/padok-team/YATAS/commit/2f8d56686888234884f20f27703cc9a9a6bff68f))

### [0.11.10](https://github.com/padok-team/YATAS/compare/v0.11.9...v0.11.10) (2022-09-26)


### Features

* **plugins:** added checks results in interface passed to plugins ([894ab98](https://github.com/padok-team/YATAS/commit/894ab98bb9bb52b3a34f60db800be72035fb3407))

### [0.11.9](https://github.com/padok-team/YATAS/compare/v0.11.8...v0.11.9) (2022-09-23)


### Bug Fixes

* **init:** changed to latest instead of fix version for init config ([f8c39ff](https://github.com/padok-team/YATAS/commit/f8c39ff559fb0e14ea6aba824e50371532d8ae83))

### [0.11.8](https://github.com/padok-team/YATAS/compare/v0.11.7...v0.11.8) (2022-09-22)

### [0.11.7](https://github.com/padok-team/YATAS/compare/v0.11.6...v0.11.7) (2022-09-17)


### Bug Fixes

* **report:** changed from reporting to report ([5683163](https://github.com/padok-team/YATAS/commit/56831633bba8f6f27c59360618f8421e29994ed3))

### [0.11.6](https://github.com/padok-team/YATAS/compare/v0.11.5...v0.11.6) (2022-09-16)


### Features

* **plugins:** added possibility to run mods, reports plugins ([28d0caa](https://github.com/padok-team/YATAS/commit/28d0caa055dd8e0f950a37cd254045cd026237b8))

### [0.11.5](https://github.com/padok-team/YATAS/compare/v0.11.4...v0.11.5) (2022-09-16)

### [0.11.4](https://github.com/padok-team/YATAS/compare/v0.11.3...v0.11.4) (2022-09-11)


### Features

* **plugins:** can now dynamically call plugins and not only aws named plugins ([36ac556](https://github.com/padok-team/YATAS/commit/36ac5565d7bf1bb90e7cb74f810d72fbbfe6be04))

### [0.11.3](https://github.com/padok-team/YATAS/compare/v0.11.2...v0.11.3) (2022-09-09)


### Bug Fixes

* **print:** removed leftover print ([0269ce7](https://github.com/padok-team/YATAS/commit/0269ce7bf2cc4630587e3ed4ffb99040ac5d842a))

### [0.11.2](https://github.com/padok-team/YATAS/compare/v0.11.1...v0.11.2) (2022-09-09)


### Bug Fixes

* **plugin:** fixed registration of interface when no parameters ([e8dd022](https://github.com/padok-team/YATAS/commit/e8dd022fbb9f17e73b58f07b52e3340ea8a9d832))

### [0.11.1](https://github.com/padok-team/YATAS/compare/v0.11.0...v0.11.1) (2022-09-09)


### Features

* **config:** interface ([098090d](https://github.com/padok-team/YATAS/commit/098090d42ec09e027845b259a330a9d5aa74c4da))

## [0.11.0](https://github.com/padok-team/YATAS/compare/v0.10.8...v0.11.0) (2022-09-09)


### Features

* **plugins:** configuration for plugins passed as interface ([5a00c38](https://github.com/padok-team/YATAS/commit/5a00c381bf8aea72dcefd8f569e340a0f3298820))

### [0.10.8](https://github.com/padok-team/YATAS/compare/v0.10.7...v0.10.8) (2022-09-09)
