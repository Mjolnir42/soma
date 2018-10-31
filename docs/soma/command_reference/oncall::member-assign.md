# DESCRIPTION

This command is used to assign a user to an oncall duty team.

If the oncall duty name is specified as a valid UUID, that ID is
used as the oncallID of the oncall duty to assign the user to.

If the username is specified as a valid UUID, that ID is
used as the userID of the user to assign.

# SYNOPSIS

```
soma oncall member assign ${username} to ${oncallduty}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
username | string | Name of the user to assign | | no
oncallduty | string | Name of the oncall duty | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | oncall | member-assign | yes | no

# EXAMPLES

```
soma oncall member assign jd to "Emergency Phone"
```
