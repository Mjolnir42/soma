# oncall duty management

Oncall combines the function for maintaining the information base about
oncall duty teams within SOMA.

Oncall duty teams can be assigned per-view as a property in SOMA
configuration trees.

# SYNOPSIS OVERVIEW

```
soma oncall add ${name} phone ${extension}
soma oncall update ${current-name} [phone ${new-extension}] [name ${new-name}]
soma oncall remove ${name}
soma oncall show ${name}
soma oncall list
soma oncall member assign ${username} to ${oncallduty}
soma oncall member unassign ${username} from ${oncallduty}
soma oncall member list ${oncallduty}
```

See `soma oncall help ${command}` or `soma oncall member help
${command}` for detailed help.
