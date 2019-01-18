# DESCRIPTION

This command is used to search for buckets by specific attributes.
All specified search conditions must be true for a bucket to be part
of the result set. While all conditions are optional, a search with no
condition is not accepted.

# SYNOPSIS

```
soma bucket search [id ${uuid}] [name ${bucket}] [repository ${repository}] [environment ${environment}] [deleted ${isDeleted}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
id | uuid | UUID of the bucket to search | | yes
name | string | Name of the bucket to search | | yes
repository | string | Name of the repository of the bucket to search | | yes
environment | string | Name of the environment of the bucket to search | | yes
isDeleted | boolean | Whether the bucket must be flagged as deleted | | yes

# PERMISSIONS

The request is authorized if the user has at least one of the sufficient
permissions or all required permissions.
Repository scoped permissions must be granted on a repository and adds
it to the pool of possible results.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | repository | no | yes
repository | bucket | search | no | yes

# EXAMPLES

```
soma bucket search deleted true
soma bucket search environment live
```
