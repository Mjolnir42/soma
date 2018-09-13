# DESCRIPTION

This command is used to remove state definitions from SOMA.

# SYNOPSIS

```
soma state remove ${state}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
state | string | Name of the state | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | state | remove | yes | no

# EXAMPLES

```
soma state remove clustered
soma state remove grouped
```
