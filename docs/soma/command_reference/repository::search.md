# DESCRIPTION

This command is used to search for repositories by specific attributes.
All search conditions must be true for a repository to be part of the
result set. While all conditions are optional, a search with no
condition is not accepted.

# SYNOPSIS

```
soma repository search [id ${uuid}] [name ${repository}] [team ${team}] [deleted ${isDeleted}] [active ${isActive}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
id | uuid | UUID of the repository to search | | yes
repository | string | Name of the repository to search | | yes
team | string | Name of the team owning the repositories | | yes
isDeleted | boolean | Whether the repository must be flagged as deleted | | yes
isActive | boolean | Whether the repository must be flagged as active | | yes

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Repository scoped permissions must be granted on a repository and adds
it to the pool of possible results.
Team scoped permissions must be granted on a team and add all
repositories owned by the team to the pool of possible results.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | team | | no | yes
system | repository | no | yes
team | repository | search | no | yes
repository | repository-config | search | no | yes

# EXAMPLES

```
soma repository search deleted true
soma reposutory search team 'Example Team'
```
