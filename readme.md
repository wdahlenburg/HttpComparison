# HttpComparison

This small library attempts to determine the HTTP responses that are significantly different from a baseline.

This is useful for fuzzing or other testing where you are looking to correlate certain inputs having a significant change in response.

## Determining Significance

The significance threshold depends on the application using the library. In certain cases the threshold should be very large to require major changes to be considered significant. In other scenarios a small  threshold will allow minor changes to be considered significant. Some tuning is required to figure out a threshold that works for your application.

## Algorithm

The base algorithm does a string comparison with the Jaccard index. It applies a 25% weight reduction if the HTTP status code changes.

## Enhancements

If you think you have a better algorithm, I'd encourage you to submit a PR and it'll get well tested.

The algorithm needs to balance determining significant differences vs speed. The https://github.com/adrg/strutil repo was used to initially evaluate different similarity indexes and speed.