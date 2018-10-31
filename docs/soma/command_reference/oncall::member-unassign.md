# DESCRIPTION

This command is used to unassign a user from an oncall duty team.

If the oncall duty name is specified as a valid UUID, that ID is
used as the oncallID of the oncall duty to unassign the user from.

If the username is specified as a valid UUID, that ID is
used as the userID of the user to unassign.

# SYNOPSIS

```
soma oncall member unassign ${username} from ${oncallduty}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
username | string | Name of the user to unassign | | no
oncallduty | string | Name of the oncall duty | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | oncall | member-unassign | yes | no

# EXAMPLES

```
soma oncall member unassign jd from "Emergency Phone"
```
