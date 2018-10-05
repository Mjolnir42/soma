# DESCRIPTION

This command shows details for a monitoring system mode defined in SOMA.

# SYNOPSIS

```
soma mode show ${mode}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
mode | string | Name of the mode | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | mode | show | yes | no

# EXAMPLES

```
soma mode show private
```
