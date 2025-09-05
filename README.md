# bloom-filter

$$
m = -\frac{n \ln p}{(\ln 2)^2}
$$

$$
k = \frac{m}{n} \ln 2
$$


## Definition
m = total number of bits in the Bloom filter’s bit array
n = expected number of distinct items
p = target false positive rate
k = number of hash functions

## Feature
- Given n and p, calculate m and k
- Draw the statistic result to find the sweet point of m and k for given n and p 