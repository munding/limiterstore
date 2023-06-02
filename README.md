# LimiterStore

LimiterStore is a storage for the `rate.Limiter` structs from the "golang.org/x/time/rate" package. It allows you to store and retrieve previously created `rate.Limiter` using string keys, such as IP addresses.

## Features

1. **Thread-Safe**: LimiterStore uses the `sync.Map` structure from the standard library to safely handle concurrent read/write requests.
2. **Automatic Cleanup**: It automatically cleans up unused limiters based on a specified `cleanupInterval`.
3. **Dynamic Update**: The rate limit and burst size can be updated at runtime based on a specified `updateInterval`. If a limiter hasn't been updated in the given interval, it will be automatically updated with the provided rate limit and burst size.

## Usage

### Importing the Package

```go
import (
    "github.com/munding/limiterstore"
)
```

### Creating a new LimiterStore

```go
// creating a new LimiterStore with cleanupInterval and updateInterval
limiterStore := limiterstore.NewLimiterStore(cleanupInterval, updateInterval)
```

### Loading and Updating a Rate Limiter

```go
// loading and updating a rate limiter
limiter := limiterStore.LoadAndUpdate(key, rateLimit, burst)
```

This will return a `rate.Limiter` from the "golang.org/x/time/rate" package. If a limiter with the provided key doesn't exist, a new one will be created with the provided rate limit and burst size.

The `LoadAndUpdate` function also checks if the last update time exceeds the specified `updateInterval`. If it does, the limiter's rate limit and burst size will be updated.

### Cleanup of Unused Limiters
The `cleanupUnusedLimiters` function runs in the background and cleans up limiters that haven't been used in the specified cleanupInterval. It uses a `time.Ticker` to trigger the cleanup operation at regular intervals.

This function is automatically started when a new `LimiterStore` is created and doesn't need to be manually managed.

## Contributions
Feel free to contribute to this project by creating issues for bugs or enhancements, or by submitting pull requests. Please follow the standard Go coding conventions and make sure your code is properly tested.

## License
This project is licensed under the MIT License.