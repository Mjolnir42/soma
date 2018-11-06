# DESCRIPTION

This command removes a server from SOMA by flagging it as deleted
within the database.

# SYNOPSIS

```
soma server remove ${name}
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
global | server | remove | yes | no

# EXAMPLES

```
soma server remove example-server-a
soma server remove example-server-b
```
