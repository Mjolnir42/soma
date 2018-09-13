# DESCRIPTION

This command is used to list all permissions in a specific category.

# SYNOPSIS

```
soma permission list in ${category}
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
soma permission list in global
soma permission list in self
soma permission list in repository
soma permission list in monitoring
```
