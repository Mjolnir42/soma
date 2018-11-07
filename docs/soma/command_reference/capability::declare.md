# DESCRIPTION

This command is used to declare a monitoring system capability within SOMA.

# SYNOPSIS

```
soma capability declare ${monitoring} view ${view} metric ${path} thresholds ${num}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
monitoring | string | Name of the monitoring system | | no
view | string | Name of the view | | no
path | string | Metric path (name) of the metric | | no
num | integer | Number of supported thresholds for this metric | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Permissions in category monitoring must be granted on the specific
monitoring systems.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | monitoring | | no | yes
monitoring | capability | declare | yes | no

# EXAMPLES

```
soma capability declare ExampleMonitoring \
     view internal \
     metric icmp.echo.rtt \
     thresholds 3
```
