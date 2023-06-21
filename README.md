# gomod-version-lint

gomod-version-lint analysis dependency version in `go.mod` and report dependency that not match with provided branch regex

# example

```
gomod-version branches --module github.com/katanomi/.* --branches-not (^main$|^release-.*$) -o json

### report
{
    "name": "github.com/katanomi/builds",
    "git": {
        "repository": "https://github.com/katanomi/pkg",
    },
    "commit": "xx",

    "depends": {
        {
            "line": 112,
            "name": "github.com/katanomi/pkg",
            "git": {
                "repository": "https://github.com/katanomi/pkg",
                "branch": "fix/bug"
            }
            "version": "v0.7.1-0.20230617092407-120ce1cbefce"
        }
    }
}
```

```
# comment pr
gomod-version comment --report ./report.json --pull-request 1122 --commit 1qaz2wsx
gomod-version comment --report ./report.json --commit xx
```

# internal

``` sh
git clone --filter=blob:none --no-checkout https://github.com/katanomi/builds
cd builds
git branch -r --contains e93ab8d
```
