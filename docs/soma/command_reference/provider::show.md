# DESCRIPTION

This command shows details for a provider defined in SOMA.

Provider names may not contain / characters.

# SYNOPSIS

```
soma provider show ${provider}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
provider | string | Name of the provider | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | provider | show | yes | no

# EXAMPLES

```
soma provider show prometheus
```
