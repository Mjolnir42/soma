# DESCRIPTION

This command is used to remove teams from SOMA. Only unused teams can be
deleted since this action represents a hard delete, not a deletion flag.

# SYNOPSIS

```
soma team-mgmt remove ${team}
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
identity | team-mgmt | remove | yes | no

# EXAMPLES

```
soma team-mgmt remove wheel
```
