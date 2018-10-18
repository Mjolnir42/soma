# DESCRIPTION

This command is used to show details about a team from SOMA.

# SYNOPSIS

```
soma team-mgmt show ${team}
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
identity | team-mgmt | show | yes | no

# EXAMPLES

```
soma team-mgmt show wheel
```
