# DESCRIPTION

This command is used to show a metric definition from SOMA.

Metric names may not contain / characters.

# SYNOPSIS

```
soma metric show ${metric}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
metric | string | Name of the metric | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | metric | show | yes | no

# EXAMPLES

```
soma metric show icmp.echo.rtt
```
