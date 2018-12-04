# DESCRIPTION

This command is used to list all available repositories.

# SYNOPSIS

```
soma repository list
```

# ARGUMENT TYPES

This command takes no argument.

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Repository scoped permissions must be granted on a repository and adds
that repository to the result.
Team scoped permissions must be granted on a team and adds all repositories
owned by that team to the result.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | team | no | yes
system | repository | | no | yes
team | repository | list | no | yes
repository | repository-config | list | no | yes

# EXAMPLES

```
soma repository list
```
