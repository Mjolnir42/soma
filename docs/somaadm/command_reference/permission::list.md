# somaadm permission list

This command is used to list all permissions in a specific category.

# SYNOPSIS

```
somaadm permission list in ${category}
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
permission | permission | list | yes | no

# EXAMPLES

```
./somaadm permission list in global
./somaadm permission list in self
./somaadm permission list in repository
./somaadm permission list in monitoring
```
