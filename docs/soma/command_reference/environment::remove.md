# DESCRIPTION

This command is used to remove environment definitions from SOMA.

Environments must be unused for them to be able to be removed.

Environment names must not contain `/` characters.

# SYNOPSIS

```
soma environment remove ${environment}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
environment | string | Name of the environment | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | environment | remove | yes | no

# EXAMPLES

```
soma environment remove QA
```
