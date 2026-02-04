# PingIV

This is an essential tool for network administrators who still use legacy IPv4. The stack feels old, like the ancients Romans, maybe it is, who knows.

This is not a wrapper for the ping command, it's implemented from scratch (well almost the pro-bing library has been used)


### Usage

```
usage: pingIV [OPTIONS] <Roman IP address>

This ping utility accepts Roman numerals IPv4 addresses.

Options:
  -c int
    	Number of ping packets to send (default 4)
  -i duration
    	Interval between pings (default 1s)
  -t duration
    	Timeout for each ping (default 5s)
  -v	Verbose output (show conversion)

Examples:
  pingIV CXXVII.N.N.I                  # Ping 127.0.0.1
  pingIV CXXVII...I                    # Ping 127.0.0.1 (725 BC format)
  pingIV CXXVII.nulla.nulla.I          # Ping 127.0.0.1 (725 BC latin format)
  pingIV -c 10 VIII.VIII.VIII.VIII     # Send 10 pings to 8.8.8.8
  pingIV -v CXCII.CLXVIII.I.I          # Verbose mode (show the coversion on top)

For the compatiblity notation after 725 BC there is an automatically pad with zero/s
```

### Build

Requiments golang compiler and run the commmand:

If you want to build for your platform:

```make```

If you want to cross compile:

```make all```


### FAQ

*Q.* Will there be an IPv6 version?

*A.* I don't think so there are too many overlaps with the Latin letters, not even in Roman times was IPv6 fully supported imagine in the years MM

*Q.* Is $random feature from standard ping supported?

*A.* You can open an issue or write to me, and maybe I will implement it or find a very compelling historical argument to not do it

*Q.* Why additive notation do weird things with more than 3 symbols?

*A.* Well go to some historical building, you will found out that not even the ancient Roman used a single standard


*Q.* What the hell are Roman Numerals?

*A.* It's never to late to learn something new, check Wikipedia [Roman Numerals](https://en.wikipedia.org/wiki/Roman_numerals)


### License

CC-BY-SA-4.0 license
