# DESCRIPTION

This command is used to delete a permission. Deleting a permission
revokes all grants of that permission.

The permission to remove can be specified via shorthand or regular
syntax.

# SYNOPSIS

```
soma permission remove [${category}::]${permission} [from ${category}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
permission | string | Name of the permission | | no
category | string | Name of the category | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | permission | remove | yes | no

# EXAMPLES

```
soma permission remove auditor from permission
soma permission remove permission::designer
soma permission remove self::information
soma permission remove browse from global
```
