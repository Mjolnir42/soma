# DESCRIPTION

This command is used to list all users from a specific team.

# SYNOPSIS

```
soma team-mgmt member list ${team}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
team | string | Name of the team | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | team-mgmt | member-list | yes | no

# EXAMPLES

```
soma team-mgmt member list wheel
```
