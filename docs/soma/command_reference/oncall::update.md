# DESCRIPTION

This command is used to update an oncall duty team in SOMA. While both the
phone and the name update are optional, the client requires at least one of
them to be specified to avoid requesting a nop action from the server.

The new oncall duty name must not be formatted as a UUID.
The new phone extension must be a number of 4 or 5 digits length.

If the current oncall duty name is specified as a valid UUID, that ID is
used as the oncallID of the oncall duty to update.

# SYNOPSIS

```
soma oncall update ${current-name} [phone ${new-extension}] [name ${new-name}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
current-name | string | Current name of the oncall duty | | no
new-extension | integer | Numeric phone extension of this oncall duty | | yes
new-name | string | New name of the oncall duty | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | oncall | update | yes | no

# EXAMPLES

```
soma oncall update "Emergency Phone" phone 4321
soma oncall update "Emergncey Phone" name "Emergency Phone"
soma oncall update "Emergency Phone" name "Emergency Oncall Duty" phone 6667
```
