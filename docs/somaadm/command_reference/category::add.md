# somaadm category add

This command is used to add a new permission category
to the system. Categories group permission sections
with the same scope, ie. if the actions inside the
section are for example global or per monitoringsystem.

The list of categories is defined by what is used by the
server's code. All categories have to be created.

# SYNOPSIS

```
soma category add ${category}
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
permission | category | add | yes | no

# EXAMPLES

```
soma category add global
soma category add repository
soma category add monitoring
```
