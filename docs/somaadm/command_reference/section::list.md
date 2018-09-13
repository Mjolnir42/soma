# somaadm section list

This command lists all permission sections registered in the system.

# SYNOPSIS

```
somaadm section list in ${category}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
category | string | Name of the category | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | section | list | yes | no

# EXAMPLES

```
./somaadm section list in global
```
