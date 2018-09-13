# somaadm permission unmap

This command is used to unmap specific sections or actions from a
permission, removing them from what the permission authorizes to perform.

Sections can only be unmapped if sections were originally mapped, ie. it
is not possible to unmap a section to unmap individually mapped
actions.

# SYNOPSIS

```
somaadm permission unmap ${section}[::${action}] from ${category}::${permission}
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
permission | permission | unmap | yes | no

# EXAMPLES

```
./somaadm permission unmap right::grant from permission::admin
./somaadm permission unmap right::revoke from permission::admin
./somaadm permission unmap datacenter from global::dc-maintainer
```
