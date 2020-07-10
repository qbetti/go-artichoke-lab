# Laboratory for Go-Artichoke
Experiments for collecting data from [Go-Artichoke](https://github.com/qbetti/go-artichoke) implementation of peer-action sequences.

## Requirements

Just make sure Go 1.13 or higher installed is installed on your system, and that its binary is added to your `PATH` environment variable.

## Disclaimer

Some of the experiments can use a huge amout of memory, so be careful (up to 10 GB for some of them).
All experiments have been successfully run with 16 GB of RAM, but in case you have a less memory, I recommend you change a few parameters for the appropriate functions (like the `maxLogFactor` or the `actionSizes` values).

## Running the experiments

Download the package with Go:

```
go get -u github.com/qbetti/go-artichoke-lab
```

Just run the package with the following command:
```shell
~/go/src/github.com/qbetti/go-artichoke-lab>$ run lab.go
```

... And that's it! The experiments should be running!

If you have a dependency error, I suggest that you install manually the two dependencies of this project: [Go-Ethereum](https://github.com/ethereum/go-ethereum/) (this one may be tricky) and [Go-Artichoke](https://github.com/qbetti/go-artichoke).

```
go get -u github.com/ethereum/go-ethereum
go get -u github.com/qbetti/go-artichoke
```


## Collecting data

Upon completion, experiments' results can be found in the `data/` directory as CSV files.

