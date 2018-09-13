# somaadm permission remove

This command is used to delete a permission. Deleting a permission
revokes all grants of that permission.

The permission to remove can be specified via shorthand or regular
syntax.

# SYNOPSIS

```
somaadm permission remove ${permission} from ${category}
somaadm permission remove ${category}::${permission}
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
./somaadm permission remove auditor from permission
./somaadm permission remove permission::designer
./somaadm permission remove self::information
./somaadm permission remove browse from global
```
