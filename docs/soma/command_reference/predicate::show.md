# DESCRIPTION

This command is used to show details about a threshold predicate in SOMA.

Predicates must not contain `/` characters.

# SYNOPSIS

```
soma predicate show ${pred}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
pred | string | Symbol of the predicate | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | predicate | show | yes | no

# EXAMPLES

```
soma predicate show '>='
```
