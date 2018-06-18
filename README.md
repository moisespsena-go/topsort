# go-topsort

Topological Sorting for Golang

It's **fork** of [stevenle/topsort](https://github.com/stevenle/topsort).

Topological sorting algorithms are especially useful for dependency calculation, and so this particular implementation is mainly intended for this purpose. As a result, the direction of edges and the order of the results may seem reversed compared to other implementations of topological sorting.

For example, if:

* A depends on B
* B depends on C
* B depends on D
* E depends on D

The graph is represented as:

![Graph image](https://www.planttext.com/plantuml/img/SoWkIImgAStDuKh9J2zABCXGS5Uevb800aS5NJi59p3J2SHqHZ1Tm4nN2BDMWSiXDIy5Q0G0)

## INSTALLATION

```go get -u github.com/moisespsena/go-topsort/topsort```

## CLI

```
$GOPATH/bin/topsort -h

# or

export PATH=$GOPATH/bin:$PATH
topsort -h
```

Output:

```
Topological sorting algorithms are especially useful for dependency calculation, 
and so this particular implementation is mainly intended for this purpose. 

As a result, the direction of edges and the order of the results may seem reversed 
compared to other implementations of topological sorting.

Home Page: https://github.com/moisespsena/go-topsort

EXAMPLES
--------

$ echo "A-B,B-C,B-D,E-D,F" | topsort
$ echo "A-B:B-C:B-D:E-D,F" | topsort -p :
$ topsort pairs.txt

Ordered input files including STDIN (file name is '-')
$ echo "A-B,B-C,B-D,E-D,F" | topsort pairs1.txt pairs2.txt - pairs3.txt

Usage:
  topsort [flags] [file...]

Flags:
  -e, --edge-sep string   Set the edge separator (default "-")
  -h, --help              help for topsort
  -p, --pair-sep string   Set the pairs separator (default ",")
  -t, --toggle            Help message for toggle
  -T, --top-sort          Use Topological node classifier, otherwise, Depth-first classifier.

```

## How To

The code for previous example would look something like:

```go
import fmt
// Initialize the graph
graph := topsort.NewGraph()

// Add nodes
graph.AddNode("A", "B", "C", "D", "E")

// Add edges
graph.AddEdge("A", "B")
graph.AddEdge("B", "C")
graph.AddEdge("B", "D")
graph.AddEdge("E", "D")

// Topologically sort only node A in dependency order and your edges, but not sort D and E.
results, err := graph.TopSort("A")
if err != nil {
    panic(err)
}
fmt.Println(results) // => [C D B A]

// in depth-first order
results, err = graph.DepthFirst("A")
if err != nil {
    panic(err)
}
fmt.Println(results) // => [A B D C]
```

Sort all nodes:

```go
// Topologically sort all nodes in the graph
results, err := graph.TopSort()
if err != nil {
    panic(err)
}
fmt.Println(results) // => [B C D B A E]

// all nodes in depth-first order
results, err = graph.DepthFirst()
if err != nil {
    panic(err)
}
fmt.Println(results) // => [E A B D C]
```
See [Examples Test](examples_test.go) for more examples.
