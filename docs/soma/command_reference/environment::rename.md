# DESCRIPTION

This command renames an environment defined in SOMA. The environment
must be unused for the rename to succeed.

Environment names must not contain `/` characters.

# SYNOPSIS

```
soma environment rename ${environment} to ${new-environment}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
environment | string | Current name of the environment | | no
new-environment | string | New name of the environment | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | environment | rename | yes | no

# EXAMPLES

```
soma environment rename staeg to stage
```
