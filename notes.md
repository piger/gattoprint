# notes

on linux I see:

```
discovering services
discovering characteristics on: 0000ae3a-0000-1000-8000-00805f9b34fb
found characteristic: 0000ae3c-0000-1000-8000-00805f9b34fb
found characteristic: 0000ae3b-0000-1000-8000-00805f9b34fb

discovering characteristics on: 0000ae30-0000-1000-8000-00805f9b34fb <- this is the one
found characteristic: 0000ae03-0000-1000-8000-00805f9b34fb
found characteristic: 0000ae10-0000-1000-8000-00805f9b34fb
found characteristic: 0000ae01-0000-1000-8000-00805f9b34fb <- this is the one 
found characteristic: 0000ae05-0000-1000-8000-00805f9b34fb
found characteristic: 0000ae04-0000-1000-8000-00805f9b34fb
found characteristic: 0000ae02-0000-1000-8000-00805f9b34fb

service: 0000ae3a-0000-1000-8000-00805f9b34fb
```

```
notification (getDevState):
0x51 0x78 0xA3 0x1 0x3 0x0 0x0 0xE 0x24 0x2A 0xFF
[--------] header                       crc? [tail]
          0xA3 = cmdGetDevState
```
