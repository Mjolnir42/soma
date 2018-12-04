# DESCRIPTION

This command is used to destroy a repository.

# SYNOPSIS

```
soma repository destroy ${repository} [from ${team}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Name of the repository | | no
team | string | Name of the team owning the repository | | yes

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Team scoped permissions must be granted on a team and allow to destroy
all repositories of that team.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | team | | no | yes
team | repository | destroy | yes | no

# EXAMPLES

```
soma repository destroy example
soma repository destroy example from 'Example Team'
```
