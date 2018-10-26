# DESCRIPTION

This command shows details about a monitoring system defined in SOMA.

# SYNOPSIS

```
soma monitoringsystem-mgmt show ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the monitoring system | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions. Permissions in category
`monitoring` must be granted on the specific monitoring system being
shown.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | monitoring | no | yes
monitoring | monitoringsystem | show | yes | no

# EXAMPLES

```
soma monitoringsystem-mgmt show ExampleMonitoring
```
