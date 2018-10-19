# DESCRIPTION

This command lists all per-team service properties defined in SOMA
of a specific team.

# SYNOPSIS

```
soma property-mgmt service list team ${team}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
team | string | Name of the team | | no


# PERMISSIONS

The request is authorized if the user either has at least one
sufficient, all system or all required permissions.

Category | Section | Action | Required | System | Sufficient
 ------- | ------- | ------ | -------- | ------ | ----------
omnipotence | | | no | no | yes
system | global | | no | yes | no
system | team | | no | yes | no
global | property-mgmt | list | yes | no | no
team | property-service | list | yes | no | no

# EXAMPLES

```
soma property-mgmt service list team ExampleTeam
```
