# DESCRIPTION

This command is used to rename a repository. Since the repository name
is a prefix to the bucket names in that repository, it also renames all
buckets by updating the prefix.

# SYNOPSIS

```
soma repository rename ${repository} to ${newName} [from ${team}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Current name of the repository | | no
newName | string | New name of the repository | | no
team | string | Name of the team owning the repository | | yes

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Team scoped permissions must be granted on a team and allow to rename
all repositories of that team.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | team | | no | yes
team | repository | rename | yes | no

# EXAMPLES

```
soma repository rename eaxmple to example
```
