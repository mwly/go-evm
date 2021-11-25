# go-evm

## Rewriting agenda:

### Data Pipeline:

#### The Idea of the following table was to visualize the mutations the videostream will go through
#### The boolean WG ist set true if the Mutation can't be performed on one element of the previous process.
#### In other words it has to wait for the that the complete data of previous step is processed.

| data format             | amount to process   | Extract                 | amount to process       |
| ------------            | ------              | -----------             | ------                  |
| Videostream (x,y,ch)    | 1                   | Frames                  | N                       |
| Frames                  | N                   | Gaussian pyramid        | Levels * N              |
| Gaussian pyramid        | Levels * N          | Laplacian Pyramid       | Levels * N              |
| Laplacian pyramid       | Levels * N          | Timeline of length N    | Levels * x * y * ch     |
| Timeline of length N    | Levels * x * y * ch | FFT                     | Levels * x * y * ch     |
| FFT                     | Levels * x * y * ch | FFT + Amp * Filter(FFT) | Levels * x * y * ch     |
| FFT + Amp * Filter(FFT) | Levels * x * y * ch | iFFT                    | Levels * x * y * ch     |
| iFFT                    | Levels * x * y * ch | Real                    | Levels * x * y * ch     |
| Real                    | Levels * x * y * ch | Timeline of length N    | Levels * x * y * ch     | 
| Timeline of length N    | Levels * x * y * ch | Laplacian Pyramid       | Levels * N              |
| Laplacian Pyramid       | Levels * x * y * ch | Frames                  | N                       |
| Frames                  | N                   | magnified Video         | 1                       |






