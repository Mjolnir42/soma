# DESCRIPTION

This command shows details about a server from SOMA.

# SYNOPSIS

```
soma server show ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the server | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | server | show | yes | no

# EXAMPLES

```
soma server show example-server-a
soma server show example-server-b
```
