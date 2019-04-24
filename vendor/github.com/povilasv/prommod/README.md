![prommod](https://api.travis-ci.com/povilasv/prommod.svg?branch=master)
![prommod](https://goreportcard.com/badge/github.com/povilasv/prommod)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpovilasv%2Fprommod.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpovilasv%2Fprommod?ref=badge_shield)
[![GoDoc](https://godoc.org/github.com/povilasv/prommod?status.svg)](https://godoc.org/github.com/povilasv/prommod)

# prommod

Export Go Module information to Prometheus.

Should work with any recent version of Go. Tested with Go versions starting 1.10.

# Download

```
go get github.com/povilasv/prommod
```

With modules:

```
GO111MODULE=on; go get github.com/povilasv/prommod
```

# Usage

```
import (
	"fmt"
	"log"
	"net/http"

	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.Register(prommod.NewCollector("app_name"))

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

# Example Metric Output

```
# HELP go_mod_info A metric with a constant '1' value labeled by dependency name, version, from which app_name was built.
# TYPE go_mod_info gauge
go_mod_info{name="github.com/beorn7/perks",program="app_name",version="v0.0.0-20180321164747-3a771d992973"} 1
go_mod_info{name="github.com/golang/protobuf",program="app_name",version="v1.2.0"} 1
go_mod_info{name="github.com/matttproud/golang_protobuf_extensions",program="app_name",version="v1.0.1"} 1
go_mod_info{name="github.com/povilasv/prommod",program="app_name",version="v0.0.11-0.20190309143328-e661980fc053"} 1
go_mod_info{name="github.com/prometheus/client_golang",program="app_name",version="v0.9.2"} 1
go_mod_info{name="github.com/prometheus/client_model",program="app_name",version="v0.0.0-20180712105110-5c3871d89910"} 1
go_mod_info{name="github.com/prometheus/common",program="app_name",version="v0.0.0-20181126121408-4724e9255275"} 1
go_mod_info{name="github.com/prometheus/procfs",program="app_name",version="v0.0.0-20181204211112-1dc9a6cbc91a"} 1
```

# Example Print

```
fmt.Println(prommod.Print("app_name"))
```

Output:

```
app_name
 github.com/beorn7/perks: v0.0.0-20180321164747-3a771d992973
 github.com/golang/protobuf: v1.2.0
 github.com/matttproud/golang_protobuf_extensions: v1.0.1
 github.com/povilasv/prommod: v0.0.3
 github.com/prometheus/client_golang: v0.9.2
 github.com/prometheus/client_model: v0.0.0-20180712105110-5c3871d89910
 github.com/prometheus/common: v0.0.0-20181126121408-4724e9255275
 github.com/prometheus/procfs: v0.0.0-20181204211112-1dc9a6cbc91a
```

# Example Info

```
fmt.Println(prommod.Info())
```

Output:

```
(github.com/beorn7/perks: v0.0.0-20180321164747-3a771d992973, github.com/golang/protobuf: v1.2.0, github.com/matttproud/golang_protobuf_extensions: v1.0.1, github.com/povilasv/prommod: v0.0.5, github.com/prometheus/client_golang: v0.9.2, github.com/prometheus/client_model: v0.0.0-20180712105110-5c3871d89910, github.com/prometheus/common: v0.0.0-20181126121408-4724e9255275, github.com/prometheus/procfs: v0.0.0-20181204211112-1dc9a6cbc91a)
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpovilasv%2Fprommod.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpovilasv%2Fprommod?ref=badge_large)
