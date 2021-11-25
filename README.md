# go-evm

## Rewriting agenda:

### Data Pipeline:

#### The Idea of the following table was to visualize the mutations the videostream will go through
#### The boolean WG ist set true if the Mutation can't be performed on one element of the previous process.
#### In other words it has to wait for the that the complete data of previous step is processed.

| data format             | amount to process   | Extract                 | amount to process       | WG   |
| ------------            | ------              | -----------             | ------                  | ---- |
| Videostream (x,y,ch)    | 1                   | Frames                  | N                       | f    |
| Frames                  | N                   | Gaussian pyramid        | Levels * N              | f    |
| Gaussian pyramid        | Levels * N          | Laplacian Pyramid       | Levels * N              | f    |
| Laplacian pyramid       | Levels * N          | Timeline of length N    | Levels * x * y * ch     | true |
| Timeline of length N    | Levels * x * y * ch | FFT                     | Levels * x * y * ch     | f    |
| FFT                     | Levels * x * y * ch | FFT + Amp * Filter(FFT) | Levels * x * y * ch     | f    |
| FFT + Amp * Filter(FFT) | Levels * x * y * ch | iFFT                    | Levels * x * y * ch     | f    |
| iFFT                    | Levels * x * y * ch | Real                    | Levels * x * y * ch     | f    |
| Real                    | Levels * x * y * ch | Timeline of length N    | Levels * x * y * ch     | f    |
| Timeline of length N    | Levels * x * y * ch | Laplacian Pyramid       | Levels * N              | true |
| Laplacian Pyramid       | Levels * x * y * ch | Frames                  | N                       | f    |
| Frames                  | N                   | magnified Video         | 1                       | f    |






