# DESCRIPTION

This command is used to remove an oncall duty team from SOMA.

If the oncall duty name is specified as a valid UUID, that ID is
used as the oncallID of the oncall duty to remove.

# SYNOPSIS

```
soma oncall remove ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the oncall duty | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | oncall | remove | yes | no

# EXAMPLES

```
soma oncall remove "Emergency Phone"
```
