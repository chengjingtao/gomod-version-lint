# gomod-version-lint

gomod-version-lint analysis dependency version in `go.mod` and report dependency that not match with provided branch regex

# dependency git revision of dependency checking

print all dependency branches

``` bash
gomod-version-lint branches --module "github.com/demo/.*"

✅️  github.com/demo/demo1@v0.0.0-20230314042448-bf45d9fa206a                                 create-catalog-info-yaml,main,release-0.7
⚠️   github.com/demo/demo2@v0.7.1-0.20230620020346-5e946b016f71s
✅️  github.com/demo/demo3@v0.7.0                                                        0.7,main,release-0.7 

```

# git file comment

comment on git file in pull request
```bash
gomod-version-lint comment --commit-id a13889b --pr-id 1208 --repo https://gitlab.example/demo/demo --file ./.git-comments
```
# internal

``` sh
git clone --filter=blob:none --no-checkout https://github.com/demo/demo-1
cd builds
git branch -r --contains e93ab8d
```
