# Listener Docker Logs on go

1. Create yml file with configs. What docker container we should check
2. docker logs container_name --tail 1000 for example
3. Read logs and find by regexp special error and critical logs
4. Set up UPD client and send alerts on UDP server



```yml
containers:
  - name: test
    regexp: "[CRITICAL]"

  - name: mysql
    regexp: "[ERROR]"
```
