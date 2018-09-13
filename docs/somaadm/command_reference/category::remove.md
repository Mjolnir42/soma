# somaadm category remove

This command removes a permission category from the
server. Using this should only be required to correct
a typo during `category add`.

# SYNOPSIS

```
soma category remove ${category}
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
permission | category | remove | yes | no


# EXAMPLES

```
soma category remove global
```
