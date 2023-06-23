# work log

## log your tasks
```bash
$ wlog stuff
$ wlog forgotten stuff at 11:00
$ wlog other stuff at 10:00
$ wlog forgotten stuff at 9:00
$ wlog more forgotten stuff at 15:00 yesterday
$ wlog even more forgotten stuff at 13:45 tuesday
$ wlog eod # mark the end of your day
```

## print your log
```bash
$ wlog
$ wlog -l 100 # print latest 100 entries
```

## print how long you spent on each task
```bash
$ wlog -ld
$ wlog -ld 100 # print durations from the latest 100 entries
```

## list/find your tasks
```bash
$ wlog -lt
$ wlog -lt 100 # print 100 tasks
$ wlog -lt forgotten # print tasks containing "forgotten"
```

## todo:
- support marking tasks done
- support removing entries
- support adding/editing externalIds


