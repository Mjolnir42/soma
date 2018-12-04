# DESCRIPTION

This command is used to display full information about a repository.

# SYNOPSIS

```
soma repository show ${repository} [from ${team}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Name of the repository to search | | no
team | string | Name of the team owning the repositories | | yes

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Repository scoped permissions must be granted on a repository and allows
to request that repositories' details.
Team scoped permissions must be granted on a team and authorizes to
request the details for all repositories owned by that team.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | team | no | yes
system | repository | | no | yes
team | repository | show | no | yes
repository | repository-config | show | no | yes

# EXAMPLES

```
soma repository show example
soma repository show example from 'Example Team'
```
