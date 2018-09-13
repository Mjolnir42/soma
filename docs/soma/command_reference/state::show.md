# DESCRIPTION

This command shows details for a state defined in SOMA.

# SYNOPSIS

```
soma state show ${state}
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
global | state | show | yes | no

# EXAMPLES

```
soma state show grouped
```
