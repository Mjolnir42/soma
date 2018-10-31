# DESCRIPTION

This command is used to unassign a user from an oncall duty team.
This command is used to list the users of an oncall duty team.

If the oncall duty name is specified as a valid UUID, that ID is
used as the oncallID of the oncall duty to list the user of.

# SYNOPSIS

```
soma oncall member list ${oncallduty}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
oncallduty | string | Name of the oncall duty | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | oncall | member-list | yes | no

# EXAMPLES

```
soma oncall member list "Emergency Phone"
```
