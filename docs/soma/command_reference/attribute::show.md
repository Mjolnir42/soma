# DESCRIPTION

This command is used to show details about a specific attribute.

# SYNOPSIS

```
soma attribute show ${attribute}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
attribute | string | Name of the attribute | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | attribute | show | yes | no

# EXAMPLES

```
soma attribute show transport_proto
```
