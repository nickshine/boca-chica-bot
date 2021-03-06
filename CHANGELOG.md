## [1.2.13](https://github.com/nickshine/boca-chica-bot/compare/v1.2.12...v1.2.13) (2021-06-25)


### Bug Fixes

* parse any date ([82e0ce2](https://github.com/nickshine/boca-chica-bot/commit/82e0ce2b39996a76bcc0344438912283f0f62c91))

## [1.2.12](https://github.com/nickshine/boca-chica-bot/compare/v1.2.11...v1.2.12) (2021-05-21)


### Bug Fixes

* ignore rempty rows ([7f76938](https://github.com/nickshine/boca-chica-bot/commit/7f769386b9ff4b67363ac877579a64ca3f8f1af8))

## [1.2.11](https://github.com/nickshine/boca-chica-bot/compare/v1.2.10...v1.2.11) (2021-02-03)


### Bug Fixes

* more lenient time range parsing ([863b55b](https://github.com/nickshine/boca-chica-bot/commit/863b55b5b57a223b0b8af47e234d88a925a1c6b1))

## [1.2.10](https://github.com/nickshine/boca-chica-bot/compare/v1.2.9...v1.2.10) (2021-01-28)


### Bug Fixes

* allow for misspelled closure statuses ([ffee5e1](https://github.com/nickshine/boca-chica-bot/commit/ffee5e1d3cac398f37cc008925f73a8a038ddb9e))

## [1.2.9](https://github.com/nickshine/boca-chica-bot/compare/v1.2.8...v1.2.9) (2021-01-26)


### Bug Fixes

* allow for misspelled weekdays and gracefully skip unparsable closures ([3532c59](https://github.com/nickshine/boca-chica-bot/commit/3532c59268c42c67adfb6a548bea592e7c0ad71e))

## [1.2.8](https://github.com/nickshine/boca-chica-bot/compare/v1.2.7...v1.2.8) (2021-01-14)


### Bug Fixes

* account for closures that were already canceled during a time status change ([e7083de](https://github.com/nickshine/boca-chica-bot/commit/e7083de3d70018018cc0804a3f62b8650b7bd0dc))

## [1.2.7](https://github.com/nickshine/boca-chica-bot/compare/v1.2.6...v1.2.7) (2021-01-14)


### Bug Fixes

* account for mismatched date layouts ([7a56ec7](https://github.com/nickshine/boca-chica-bot/commit/7a56ec72037d0ba6c83c1e62d11463c7ba608f56))

## [1.2.6](https://github.com/nickshine/boca-chica-bot/compare/v1.2.5...v1.2.6) (2021-01-12)


### Bug Fixes

* use american spelling for canceled ([a5a95d1](https://github.com/nickshine/boca-chica-bot/commit/a5a95d1c0ae4e9aca5a9ca3bf59299bc758a90de))

## [1.2.5](https://github.com/nickshine/boca-chica-bot/compare/v1.2.4...v1.2.5) (2021-01-08)


### Bug Fixes

* increase publisher function timeout ([493e206](https://github.com/nickshine/boca-chica-bot/commit/493e20688a20c42394b9b460125697c3e24cdb6d))

## [1.2.4](https://github.com/nickshine/boca-chica-bot/compare/v1.2.3...v1.2.4) (2021-01-05)


### Bug Fixes

* don't fail on error from twitter ([d0a8b70](https://github.com/nickshine/boca-chica-bot/commit/d0a8b70f30bffa2f8fe449968a73332340666192))

## [1.2.3](https://github.com/nickshine/boca-chica-bot/compare/v1.2.2...v1.2.3) (2021-01-05)


### Bug Fixes

* don't fail on errors from twitter api ([8147aa5](https://github.com/nickshine/boca-chica-bot/commit/8147aa584033551d60718e913562fcc46907cd8b))

## [1.2.2](https://github.com/nickshine/boca-chica-bot/compare/v1.2.1...v1.2.2) (2021-01-03)


### Bug Fixes

* explicitly remove closures from db if timestamp is < current time ([475ac57](https://github.com/nickshine/boca-chica-bot/commit/475ac577930b786c2558c114aa40b14df0c09174))

## [1.2.1](https://github.com/nickshine/boca-chica-bot/compare/v1.2.0...v1.2.1) (2021-01-02)


### Bug Fixes

* change dynamodb sort key to avoid duplicates when time changes ([70caa3b](https://github.com/nickshine/boca-chica-bot/commit/70caa3b6c8a8464b0310761c4d67c0338bc204ec))

# [1.2.0](https://github.com/nickshine/boca-chica-bot/compare/v1.1.2...v1.2.0) (2021-01-01)


### Bug Fixes

* use len check instead of nil check on dynamodb event record images ([eab6318](https://github.com/nickshine/boca-chica-bot/commit/eab631826c5dd7d62dcb566eb5b106a495933a73))


### Features

* implement initial discord integration ([a726076](https://github.com/nickshine/boca-chica-bot/commit/a72607664e1934f05c6cb4ac17927f58be8dec58))

## [1.1.2](https://github.com/nickshine/boca-chica-bot/compare/v1.1.1...v1.1.2) (2020-12-28)


### Bug Fixes

* skip publishing expired closures if status is cancelled ([e24ac51](https://github.com/nickshine/boca-chica-bot/commit/e24ac51b77c6e1dc93523dce0084b0a51dde8e4e))

## [1.1.1](https://github.com/nickshine/boca-chica-bot/compare/v1.1.0...v1.1.1) (2020-12-27)


### Bug Fixes

* publish separate zips for each lambda ([ec8d709](https://github.com/nickshine/boca-chica-bot/commit/ec8d709bd1c26f4718cbe576c8f562c39a048041))

# [1.1.0](https://github.com/nickshine/boca-chica-bot/compare/v1.0.0...v1.1.0) (2020-12-27)


### Features

* tweet on start and end of closure ([3eb0a38](https://github.com/nickshine/boca-chica-bot/commit/3eb0a385905842d6b2416655e816df9cf4765f8e))
