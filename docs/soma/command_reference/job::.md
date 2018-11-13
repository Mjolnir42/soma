# job management

Job management provides all the functions provided for checking
asynchronous server jobs.

# SYNOPSIS OVERVIEW

```
soma job update
soma job show ${jobID}
soma job wait ${jobID}
soma job list outstanding
soma job list local
soma job list remote
soma job prune

soma job result-mgmt add ${result}
soma job result-mgmt remove ${result}
soma job result-mgmt show ${result}
soma job result-mgmt list
soma job result-mgmt search [id ${uuid}] [name ${result}]
```

See `soma job help ${command}`, `soma job result-mgmt help ${command}` for detailed help.
