# Golang Kinesis Consumer

Kinesis consumer applications written in Go. This library is intended to be a lightweight wrapper around the Kinesis API to read records, save checkpoints (with swappable backends), and gracefully recover from service timeouts/errors.

Some alterative options:

* [Kinesis to Firehose](http://docs.aws.amazon.com/firehose/latest/dev/writing-with-kinesis-streams.html) can be used to archive data directly to S3, Redshift, or Elasticsearch without running a consumer application. 

* [Process Kinensis Streams with Golang and AWS Lambda](https://medium.com/@harlow/processing-kinesis-streams-w-aws-lambda-and-golang-264efc8f979a) for serverless processing and checkpoint management.

## Installation

Get the package source:

    $ go get github.com/harlow/kinesis-consumer

## Overview

The consumer leverages a handler func that accepts a Kinesis record. The `Scan` method will consume all shards concurrently and call the callback func as it receives records from the stream.

```go
import(
	// ...
	consumer "github.com/harlow/kinesis-consumer"
	checkpoint "github.com/harlow/kinesis-consumer/checkpoint/redis"
)

func main() {
	var (
		app    = flag.String("app", "", "App name")
		stream = flag.String("stream", "", "Stream name")
	)
	flag.Parse()

	// new checkpoint
	ck, err := checkpoint.New(*app, *stream)
	if err != nil {
		log.Fatalf("checkpoint error: %v", err)
	}

	// new consumer
	c, err := consumer.New(ck, *app, *stream)
	if err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	// scan stream
	err = c.Scan(context.TODO(), func(r *consumer.Record) bool {
		fmt.Println(string(r.Data))
		return true // continue scanning
	})
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}

	// Note: If you need to aggregate based on a specific shard the `ScanShard` 
	// method should be leverged instead.
}
```

## Checkpoint

To record the progress of the consumer in the stream we use a checkpoint to store the last sequence number the consumer has read from a particular shard. 

This will allow consumers to re-launch and pick up at the position in the stream where they left off.

The uniq identifier for a consumer is `[appName, streamName, shardID]`

<img width="722" alt="kinesis-checkpoints" src="https://user-images.githubusercontent.com/739782/33085867-d8336122-ce9a-11e7-8c8a-a8afeb09dff1.png">

There are currently two storage types for checkpoints:

### Redis Checkpoint

The Redis checkpoint requries App Name, and Stream Name:

```go
import checkpoint "github.com/harlow/kinesis-consumer/checkpoint/redis"

// redis checkpoint
ck, err := checkpoint.New(appName, streamName)
if err != nil {
	log.Fatalf("new checkpoint error: %v", err)
}
```

### DynamoDB Checkpoint

The DynamoDB checkpoint requires Table Name, App Name, and Stream Name:

```go
import checkpoint "github.com/harlow/kinesis-consumer/checkpoint/ddb"

// ddb checkpoint
ck, err := checkpoint.New(tableName, appName, streamName)
if err != nil {
	log.Fatalf("new checkpoint error: %v", err)
}
```

To leverage the DDB checkpoint we'll also need to create a table:

<img width="659" alt="screen shot 2017-11-20 at 9 16 14 am" src="https://user-images.githubusercontent.com/739782/33033316-db85f848-cdd8-11e7-941a-0a87d8ace479.png">

## Options

The consumer allows the following optional overrides:

* Kinesis Client
* Logger

```go
// new kinesis client
svc := kinesis.New(session.New(aws.NewConfig()))

// new consumer with custom client
c, err := consumer.New(
	consumer,
	streamName,
	consumer.WithClient(svc),
)
```

## Logging

The package defaults to `ioutil.Discard` which will silence log output. This can be overridden with the preferred logging strategy:

```go
func main() {
	// ...

	// logger
	logger := log.New(os.Stdout, "consumer-example: ", log.LstdFlags)

	// consumer
	c, err := consumer.New(checkpoint, appName, streamName, consumer.WithLogger(logger))
}
```

## Contributing

Please see [CONTRIBUTING.md] for more information. Thank you, [contributors]!

[LICENSE]: /MIT-LICENSE
[CONTRIBUTING.md]: /CONTRIBUTING.md

## License

Copyright (c) 2015 Harlow Ward. It is free software, and may
be redistributed under the terms specified in the [LICENSE] file.

[contributors]: https://github.com/harlow/kinesis-connectors/graphs/contributors

> [www.hward.com](http://www.hward.com) &nbsp;&middot;&nbsp;
> GitHub [@harlow](https://github.com/harlow) &nbsp;&middot;&nbsp;
> Twitter [@harlow_ward](https://twitter.com/harlow_ward)
