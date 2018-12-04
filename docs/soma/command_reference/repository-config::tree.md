# DESCRIPTION

This command is used to export a tree representation of the repository
and its children.

# SYNOPSIS

```
soma repository dumptree ${repository}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Name of the repository | | no

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Repository scoped permissions must be granted on a repository and allow to
dump that specific repository.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | repository | | no | yes
repository | repository-config | tree | yes | no

# EXAMPLES

```
soma repository dumptree example
```
