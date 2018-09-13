# DESCRIPTION

This command shows details for an view defined in SOMA.

# SYNOPSIS

```
soma view show ${view}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
view | string | Name of the view | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | view | show | yes | no

# EXAMPLES

```
soma view show local
soma view show any
```
