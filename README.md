[![Build go app](https://github.com/Charlie-Root/npv/actions/workflows/build.yml/badge.svg)](https://github.com/Charlie-Root/npv/actions/workflows/build.yml)

# Network Path Collector

With npv, routes to multiple final destinations can be mapped and visualized. npv uses an mtr implementation and can be configured in various ways, see config.yml for this.


## Usage

### Using go source

Start a test:
```
# Clone the repo
git clone https:/github.com/Charlie-Root/npv.git

# Open the folder and test 1.1.1.1. Must run as root!
sudo go run ./main.go run single 1.1.1.1
# Result saved.
```

Check the result:
```
go run ./main.go serve
started server on http://localhost:3000
```

You can also run a batch of hosts, i have tested it with batches of 500-2000 ip's at a time. 

Create a hosts file (hosts.json) with the following structure:

```
{ "hosts": ["1.1.1.1","8.8.4.4"] }
```

Run a batch job:
```
sudo go run ./main.go run batch
```

You will see a processbar of the run :-)
```


Cleanup:
```
go run ./main.go cleanup
```

## Screenshots

