# DESCRIPTION

This comamnd is used to create a new repository for a team, which then
becomes the owner team of the repository.

The repository name must be between 4 and 128 characters long.

# SYNOPSIS

```
soma repository-mgmt create ${repository} team ${team}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Name of the repository | | no
team | string | Name of the owner team | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | repository-mgmt | add | yes | no

# EXAMPLES

```
soma repository-mgmt create example team ExampleTeam
```
