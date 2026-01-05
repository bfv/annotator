# `.githooks` directory

Because the `.git\hooks` is invisible and therefor harder to add to the repo the following construct is made. All hooks are part of the repo and are located in the `.githooks` directory, which can be committed.

To notify git to look in the `.githooks` directory, execute the following:

```
git config core.hooksPath .githooks
```

If you want this behavior for all repos:
```
git config --global core.hooksPath .githooks
```
