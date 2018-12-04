# DESCRIPTION

This command is used to change the owning team of a repository.

# SYNOPSIS

```
soma repository repossess ${repository} to ${newTeam} [from ${team}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Name of the repository | | no
newTeam | string | Name of the team that will be owning the repository | | no
team | string | Name of the team owning the repository | | yes

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Team scoped permissions must be granted on a team and allow to give away
all repositories of that team.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | team | | no | yes
team | repository | repossess | yes | no

# EXAMPLES

```
soma repository repossess example from 'Example Team' to TheUsurpers
```
