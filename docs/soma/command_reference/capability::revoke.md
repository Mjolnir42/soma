# DESCRIPTION

This command is used to revoke a monitoring system capability within SOMA.

# SYNOPSIS

```
soma capability revoke ${capability}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
capability | string | Name of the monitoring capability | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Permissions in category monitoring must be granted on the specific
monitoring systems.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | monitoring | | no | yes
monitoring | capability | revoke | yes | no

# EXAMPLES

```
soma capability revoke ExampleMonitoring.internal.icmp.echo.rtt
```
