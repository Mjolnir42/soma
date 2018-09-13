# somaadm permission map

This command is used to map specific sections or actions to the
permission, including them in what the permission authorizes to perform.

If a section is mapped, then all actions within that section are mapped,
including any actions that might be added in the future.

Only actions and sections can be mapped that are from the same category
as the permission, ie. it is not possible to grant access to global
sections or actions via a repository scoped permission.

# SYNOPSIS

```
somaadm permission map ${section}[::${action}] to ${category}::${permission}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
section | string | Name of the section | | no
action | string | Name of the action | | yes
category | string | Name of the category | | no
permission | string | Name of the permission | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | permission | map | yes | no

# EXAMPLES

```
./somaadm permission map right::grant to permission::admin
./somaadm permission map right::revoke to permission::admin
./somaadm permission map datacenter to global::dc-maintainer
```
